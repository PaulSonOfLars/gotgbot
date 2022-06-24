package handlers

import (
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// TODO: Add a "block" option to force linear processing. Also a "waiting" state to handle blocked handlers.
// TODO: Allow for timeouts (and a "timeout" state to handle that)

// The Conversation handler is an advanced handler which allows for running a sequence of commands in a stateful manner.
// An example of this flow can be found at t.me/Botfather; upon receiving the "/newbot" command, the user is asked for
// the name of their bot, which is sent as a separate message.
//
// The bot's internal state allows it to check at which point of the conversation the user is, and decide how to handle
// the next update.
type Conversation struct {
	// EntryPoints is the list of handlers to start the conversation.
	EntryPoints []ext.Handler
	// States maintain the list of states and mappings for each action.
	States map[string][]ext.Handler
	// Fallbacks is the list of handlers to end the conversation halfway (eg, with /cancel commands)
	Fallbacks []ext.Handler
	// If True, a user can restart the conversation by hitting one of the entry points.
	AllowReEntry bool

	// TODO: use sync/map to ensure concurrent safety? Or protect with rwmutex?
	// TODO: dump/restore conversation states for persistence
	// TODO: Allow for custom conversation state management without a map (interface with get/check/set)
	conversationStates map[string]string
}

func NewConversation(entryPoints []ext.Handler, states map[string][]ext.Handler, fallbacks []ext.Handler) Conversation {
	return Conversation{
		EntryPoints: entryPoints,
		States:      states,
		Fallbacks:   fallbacks,

		conversationStates: map[string]string{},
	}
}

func (c Conversation) getStateKey(ctx *ext.Context) string {
	// TODO: Need to allow for customising the state key by userid/chatid (messageid?)
	return fmt.Sprintf("%d-%d", ctx.EffectiveSender.Id(), ctx.EffectiveChat.Id)
}

// CurrentState is exposed for testing purposes.
func (c Conversation) CurrentState(ctx *ext.Context) (string, bool) {
	s, ok := c.conversationStates[c.getStateKey(ctx)]
	return s, ok
}

func (c Conversation) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	// TODO: should checkUpdate return a method pointer to execute instead of a bool?
	return c.getNextHandler(c.getStateKey(ctx), b, ctx) != nil
}

func (c Conversation) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	key := c.getStateKey(ctx)

	next := c.getNextHandler(key, b, ctx)
	if next == nil {
		// Note: this should be impossible
		return nil
	}

	var stateChange *conversationStateChange
	err := next.HandleUpdate(b, ctx)
	if !errors.As(err, &stateChange) {
		return err
	}

	if stateChange.End {
		// Mark the conversation as ended by deleting the conversation reference.
		delete(c.conversationStates, key)
	}

	if stateChange.NextState != nil {
		// If the next state is defined, then move to it.
		if _, ok := c.States[*stateChange.NextState]; !ok {
			// Check if the "next" state is a supported state.
			return fmt.Errorf("unknown state: %w", stateChange)
		}
		c.conversationStates[key] = *stateChange.NextState
	}

	if stateChange.ParentState != nil {
		// If a parent state is set, return that state for it to be handled.
		return stateChange.ParentState
	}

	return nil
}

// conversationStateChange handles all the possible states that can be returned from a conversation.
type conversationStateChange struct {
	// The next state to handle in the current conversation.
	NextState *string
	// End the current conversation
	End bool
	// Move the parent conversation (if any) to the desired state.
	ParentState *conversationStateChange
}

func (s *conversationStateChange) Error() string {
	// Avoid infinite print recursion by changing type
	type tmp *conversationStateChange
	return fmt.Sprintf("conversation state change: %+v", tmp(s))
}

// NextConversationState moves to the defined state in the current conversation.
func NextConversationState(nextState string) *conversationStateChange {
	return &conversationStateChange{NextState: &nextState}
}

// NextParentConversationState moves to the defined state in the parent conversation, without changing the state of the current one.
func NextParentConversationState(parentState *conversationStateChange) error {
	return &conversationStateChange{ParentState: parentState}
}

// NextConversationStateAndParentState moves both the current conversation state, as well as the parent conversation
// state.
// Can be helpful in the case of certain circular conversations.
func NextConversationStateAndParentState(nextState string, parentState *conversationStateChange) error {
	return &conversationStateChange{NextState: &nextState, ParentState: parentState}
}

// EndConversation ends the current conversation.
func EndConversation() error {
	return &conversationStateChange{End: true}
}

// EndConversationToParentState ends the current conversation and moves the parent conversation to the new state.
func EndConversationToParentState(parentState *conversationStateChange) error {
	return &conversationStateChange{End: true, ParentState: parentState}
}

func (c Conversation) Name() string {
	return fmt.Sprintf("conversation_%p", c.States)
}

// getNextHandler goes through all the handlers in the conversation, until finds a handler that matches.
// If no matching handler is found, returns nil.
func (c Conversation) getNextHandler(conversationKey string, b *gotgbot.Bot, ctx *ext.Context) ext.Handler {
	// Check if ongoing conversation
	currState, ok := c.conversationStates[conversationKey]
	if !ok {
		// If no ongoing conversations, we can check entrypoints
		return checkHandlerList(c.EntryPoints, b, ctx)
	}

	// If reentry is allowed, check the entrypoints again
	if c.AllowReEntry {
		if next := checkHandlerList(c.EntryPoints, b, ctx); next != nil {
			return next
		}
	}

	// else, check state mappings
	if next := checkHandlerList(c.States[currState], b, ctx); next != nil {
		return next
	}

	// TODO: do we check fallbacks BEFORE or AFTER the main states?
	// else, fallbacks -> handle any cancellations
	if next := checkHandlerList(c.Fallbacks, b, ctx); next != nil {
		return next
	}

	return nil
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
