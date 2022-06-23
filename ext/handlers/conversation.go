package handlers

import (
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// TODO: Add a "block" option to force linear processing. Also a "waiting" state to handle blocked handlers.
// TODO: Allow for timeouts (and a "timeout" state to handle that)

type Conversation struct {
	// EntryPoints is the list of handlers to start the conversation.
	EntryPoints []ext.Handler
	// States maintain the list of states and mappings for each action.
	States map[string][]ext.Handler
	// Fallbacks is the list of handlers to end the conversation halfway (eg, with /cancel commands)
	Fallbacks []ext.Handler
	// If True, a user can restart the conversation by hitting one of the entry points.
	AllowReEntry bool

	// TODO: dump/restore conversation states
	// TODO: use sync/map to ensure concurrent safety? Or protect with rwmutex?
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
	return fmt.Sprintf("%s-%s", ctx.EffectiveSender.Id(), ctx.EffectiveChat.Id)
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
		// TODO:  double-check end conversation flow when parents are used
		delete(c.conversationStates, key)
		// TODO: can we end parent conversation from child?
		// TODO: how do we know which parent state to map to?
		return nil
	}

	if _, ok := c.States[stateChange.NextState]; !ok {
		// TODO: checked that wrapped statechange remains unwrappable and is used by parent as expected
		return fmt.Errorf("unknown state: %w", stateChange)
	}
	c.conversationStates[key] = stateChange.NextState
	return nil
}

type conversationStateChange struct {
	NextState string
	End       bool
}

func (s *conversationStateChange) Error() string {
	if s.End {
		return "not an error; conversation end"
	}
	return "conversation state change to: " + s.NextState
}

func NextConversationState(nextState string) error {
	return &conversationStateChange{NextState: nextState}
}

func EndConversation() error {
	return &conversationStateChange{End: true}
}

func (c Conversation) Name() string {
	return fmt.Sprintf("conversation_%p", c.States)
}

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

func checkHandlerList(handlers []ext.Handler, b *gotgbot.Bot, ctx *ext.Context) ext.Handler {
	for _, h := range handlers {
		if h.CheckUpdate(b, ctx) {
			return h
		}
	}
	return nil
}
