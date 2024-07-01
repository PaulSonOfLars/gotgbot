package handlers

import (
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
)

// TODO: Add a "block" option to force linear processing. Also a "waiting" state to handle blocked handlers.
// TODO: Allow for timeouts (and a "timeout" state to handle that)

// ConversationFilter is much wider than regular filters, because it allows for any kind of update; we may want
// messages, commands, callbacks, etc.
type ConversationFilter func(ctx *ext.Context) bool

// The Conversation handler is an advanced handler which allows for running a sequence of commands in a stateful manner.
// An example of this flow can be found at t.me/Botfather; upon receiving the "/newbot" command, the user is asked for
// the name of their bot, which is sent as a separate message.
//
// The bot's internal state allows it to check at which point of the conversation the user is, and decide how to handle
// the next update.
type Conversation struct {
	// EntryPoints is the list of handlers to start the conversation.
	EntryPoints []ext.Handler
	// States is the map of possible states, with a list of possible handlers for each one.
	States map[string][]ext.Handler
	// StateStorage is responsible for storing all running conversations.
	StateStorage conversation.Storage

	// The following are all optional fields:
	// Exits is the list of handlers to exit the current conversation partway (eg /cancel commands)
	Exits []ext.Handler
	// Fallbacks is the list of handlers to handle updates which haven't been matched by any states.
	Fallbacks []ext.Handler
	// If True, a user can restart the conversation by hitting one of the entry points.
	AllowReEntry bool
	// Filter allows users to set a conversation-wide filter to any incoming updates. This can be useful to only target
	// one specific chat, or to avoid unwanted updates which may interfere with the conversation key strategy
	// (eg polls).
	Filter ConversationFilter
}

type ConversationOpts struct {
	// Exits is the list of handlers to exit the current conversation partway (eg /cancel commands). This returns
	// EndConversation() by default, unless otherwise specified.
	Exits []ext.Handler
	// Fallbacks is the list of handlers to handle updates which haven't been matched by any other handlers.
	Fallbacks []ext.Handler
	// If True, a user can restart the conversation at any time by hitting one of the entry points again.
	AllowReEntry bool
	// StateStorage is responsible for storing all running conversations.
	StateStorage conversation.Storage
	// Filter allows users to set a conversation-wide filter to any incoming updates. This can be useful to only target
	// one specific chat, or to avoid unwanted updates which may interfere with the conversation key strategy
	// (eg polls).
	Filter ConversationFilter
}

func NewConversation(entryPoints []ext.Handler, states map[string][]ext.Handler, opts *ConversationOpts) Conversation {
	c := Conversation{
		EntryPoints: entryPoints,
		States:      states,
		// Setup a default storage medium
		StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
	}

	if opts != nil {
		c.Exits = opts.Exits
		c.Fallbacks = opts.Fallbacks
		c.AllowReEntry = opts.AllowReEntry
		c.Filter = opts.Filter

		// If no StateStorage is specified, we should keep the default.
		if opts.StateStorage != nil {
			c.StateStorage = opts.StateStorage
		}
	}

	return c
}

func (c Conversation) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	// Note: Kinda sad that this error gets lost.
	h, _ := c.getNextHandler(b, ctx)
	return h != nil
}

func (c Conversation) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	next, err := c.getNextHandler(b, ctx)
	if err != nil {
		return fmt.Errorf("failed to get next handler in conversation: %w", err)
	}
	if next == nil {
		// Note: this should be impossible
		return nil
	}

	var stateChange *ConversationStateChange
	err = next.HandleUpdate(b, ctx)
	if !errors.As(err, &stateChange) {
		// We don't wrap this error, as users might want to handle it explicitly
		return err
	}

	if stateChange.End {
		// Mark the conversation as ended by deleting the conversation reference.
		err := c.StateStorage.Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to end conversation: %w", err)
		}
	}

	if stateChange.NextState != nil {
		// If the next state is defined, then move to it.
		if _, ok := c.States[*stateChange.NextState]; !ok {
			// Check if the "next" state is a supported state.
			return fmt.Errorf("unknown state: %w", stateChange)
		}
		err := c.StateStorage.Set(ctx, conversation.State{Key: *stateChange.NextState})
		if err != nil {
			return fmt.Errorf("failed to update conversation state: %w", err)
		}
	}

	if stateChange.ParentState != nil {
		// If a parent state is set, return that state for it to be handled.
		return stateChange.ParentState
	}

	return nil
}

// ConversationStateChange handles all the possible states that can be returned from a conversation.
type ConversationStateChange struct {
	// The next state to handle in the current conversation.
	NextState *string
	// End the current conversation
	End bool
	// Move the parent conversation (if any) to the desired state.
	ParentState *ConversationStateChange
}

func (s *ConversationStateChange) Error() string {
	// Avoid infinite print recursion by changing type
	type tmp *ConversationStateChange
	return fmt.Sprintf("conversation state change: %+v", tmp(s))
}

// NextConversationState moves to the defined state in the current conversation.
func NextConversationState(nextState string) *ConversationStateChange {
	return &ConversationStateChange{NextState: &nextState}
}

// NextParentConversationState moves to the defined state in the parent conversation, without changing the state of the current one.
func NextParentConversationState(parentState *ConversationStateChange) error {
	return &ConversationStateChange{ParentState: parentState}
}

// NextConversationStateAndParentState moves both the current conversation state and the parent conversation state.
// Can be helpful in the case of certain circular conversations.
func NextConversationStateAndParentState(nextState string, parentState *ConversationStateChange) error {
	return &ConversationStateChange{NextState: &nextState, ParentState: parentState}
}

// EndConversation ends the current conversation.
func EndConversation() error {
	return &ConversationStateChange{End: true}
}

// EndConversationToParentState ends the current conversation and moves the parent conversation to the new state.
func EndConversationToParentState(parentState *ConversationStateChange) error {
	return &ConversationStateChange{End: true, ParentState: parentState}
}

func (c Conversation) Name() string {
	return fmt.Sprintf("conversation_%p", c.States)
}

// getNextHandler goes through all the handlers in the conversation, until it finds a handler that matches.
// If no matching handler is found, returns nil.
func (c Conversation) getNextHandler(b *gotgbot.Bot, ctx *ext.Context) (ext.Handler, error) {
	// If the user has defined a filter, and this filter does NOT return true, then we do NOT want to consider this
	// update for the conversation.
	if c.Filter != nil && !c.Filter(ctx) {
		return nil, nil
	}

	// Check if a conversation has already started for this user.
	currState, err := c.StateStorage.Get(ctx)
	if err != nil {
		if errors.Is(err, conversation.ErrKeyNotFound) {
			// If this is an unknown conversation key, then we know this is a new conversation, so we check all
			// entrypoints.
			return checkHandlerList(c.EntryPoints, b, ctx), nil
		}
		// Else, we need to handle the error.
		return nil, fmt.Errorf("failed to get state from conversation storage: %w", err)
	}

	// If reentry is allowed, check the entrypoints again.
	if c.AllowReEntry {
		if next := checkHandlerList(c.EntryPoints, b, ctx); next != nil {
			return next, nil
		}
	}

	// Else, exits -> handle any conversation exits/cancellations.
	if next := checkHandlerList(c.Exits, b, ctx); next != nil {
		return wrappedExitHandler{h: next}, nil
	}

	// Else, check state mappings (the magic happens here!).
	if next := checkHandlerList(c.States[currState.Key], b, ctx); next != nil {
		return next, nil
	}

	// Else, fallbacks -> handle any updates which haven't been caught by the state or exit handlers.
	if next := checkHandlerList(c.Fallbacks, b, ctx); next != nil {
		return next, nil
	}

	return nil, nil
}

// checkHandlerList iterates over a list of handlers until a match is found; at which point it is returned.
func checkHandlerList(handlers []ext.Handler, b *gotgbot.Bot, ctx *ext.Context) ext.Handler {
	for _, h := range handlers {
		if h.CheckUpdate(b, ctx) {
			return h
		}
	}
	return nil
}

// wrappedExitHandler ensures that exit handlers return conversation ends by default.
type wrappedExitHandler struct {
	h ext.Handler
}

func (w wrappedExitHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	return w.h.CheckUpdate(b, ctx)
}

func (w wrappedExitHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	err := w.h.HandleUpdate(b, ctx)
	if err != nil {
		return err
	}
	return EndConversation()
}

func (w wrappedExitHandler) Name() string {
	return w.h.Name()
}
