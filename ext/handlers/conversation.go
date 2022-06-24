package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// TODO: Add a "block" option to force linear processing. Also a "waiting" state to handle blocked handlers.
// TODO: Allow for timeouts (and a "timeout" state to handle that)

type KeyStrategy int64

// Note: If you add a new keystrategy here, make sure to add it to the getStateKey method!
const (
	// KeyStrategySenderAndChat ensures that each sender get a unique conversation in each chats.
	KeyStrategySenderAndChat KeyStrategy = iota
	// KeyStrategySender gives a unique conversation to each sender, but that conversation is available in all chats.
	KeyStrategySender
	// KeyStrategyChat gives a unique conversation to each chat, which all senders can interact in together
	KeyStrategyChat
)

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
	// KeyStrategy defines how to calculate keys for each conversation.
	KeyStrategy KeyStrategy
	// StateStorage is responsible for storing all running conversations.
	StateStorage ConversationStorage
}

func NewConversation(entryPoints []ext.Handler, states map[string][]ext.Handler, fallbacks []ext.Handler) Conversation {
	return Conversation{
		EntryPoints: entryPoints,
		States:      states,
		Fallbacks:   fallbacks,
		// By default, conversations are per-user and per-chat; so each user gets a unique conversation for each chat.
		KeyStrategy: KeyStrategySenderAndChat,
		// Instantiate default map-based storage
		StateStorage: &ConversationStorageMap{
			lock:               sync.RWMutex{},
			conversationStates: map[string]string{},
		},
	}
}

var ConversationKeyNotFound = errors.New("conversation key not found")

// ConversationStorage allows you to define custom backends for retaining conversation states.
// If you are looking to persist conversation data, you should implement this interface with you backend of choice.
type ConversationStorage interface {
	// Get returns the state for the specified conversation key.
	//
	// If the key is not found (and as such, this conversation has not yet started), this method should return the
	// ConversationKeyNotFound error.
	Get(key string) (string, error)
	// Set updates the conversation state.
	Set(key string, state string) error
	// Delete ends the conversation, removing the key from the storage.
	Delete(key string) error
}

// ConversationStorageMap is a thread-safe in-memory implementation of the ConversationStorage interface.
type ConversationStorageMap struct {
	conversationStates map[string]string
	lock               sync.RWMutex
}

func (c *ConversationStorageMap) Get(key string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.conversationStates == nil {
		return "", ConversationKeyNotFound
	}

	s, ok := c.conversationStates[key]
	if !ok {
		return "", ConversationKeyNotFound
	}
	return s, nil
}

func (c *ConversationStorageMap) Set(key string, state string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conversationStates == nil {
		c.conversationStates = map[string]string{}
	}

	c.conversationStates[key] = state
	return nil
}

func (c *ConversationStorageMap) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conversationStates == nil {
		return nil
	}

	delete(c.conversationStates, key)
	return nil
}

// TODO: should this be exported?
func (c Conversation) getStateKey(ctx *ext.Context) string {
	switch c.KeyStrategy {
	case KeyStrategySender:
		return strconv.FormatInt(ctx.EffectiveSender.Id(), 10)
	case KeyStrategyChat:
		return strconv.FormatInt(ctx.EffectiveChat.Id, 10)
	case KeyStrategySenderAndChat:
		fallthrough
	default:
		// Default to KeyStrategySenderAndChat if unknown strategy
		return fmt.Sprintf("%d/%d", ctx.EffectiveSender.Id(), ctx.EffectiveChat.Id)
	}
}

// CurrentState is exposed for testing purposes.
// TODO: Should we un-export this?
func (c Conversation) CurrentState(ctx *ext.Context) (string, error) {
	return c.StateStorage.Get(c.getStateKey(ctx))
}

func (c Conversation) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	// Note: Kinda sad that this error gets lost.
	h, _ := c.getNextHandler(c.getStateKey(ctx), b, ctx)
	return h != nil
}

func (c Conversation) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	key := c.getStateKey(ctx)

	next, err := c.getNextHandler(key, b, ctx)
	if err != nil {
		return err
	}
	if next == nil {
		// Note: this should be impossible
		return nil
	}

	var stateChange *conversationStateChange
	err = next.HandleUpdate(b, ctx)
	if !errors.As(err, &stateChange) {
		return err
	}

	if stateChange.End {
		// Mark the conversation as ended by deleting the conversation reference.
		c.StateStorage.Delete(key)
	}

	if stateChange.NextState != nil {
		// If the next state is defined, then move to it.
		if _, ok := c.States[*stateChange.NextState]; !ok {
			// Check if the "next" state is a supported state.
			return fmt.Errorf("unknown state: %w", stateChange)
		}
		err := c.StateStorage.Set(key, *stateChange.NextState)
		if err != nil {
			return fmt.Errorf("failed to set new conversation state: %w", err)
		}
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
func (c Conversation) getNextHandler(conversationKey string, b *gotgbot.Bot, ctx *ext.Context) (ext.Handler, error) {
	// Check if ongoing conversation
	currState, err := c.StateStorage.Get(conversationKey)
	if err != nil {
		if errors.Is(err, ConversationKeyNotFound) {
			// If this is an unknown conversation key, then we know this is a new conversation, so we check all
			// entrypoints
			return checkHandlerList(c.EntryPoints, b, ctx), nil
		}
		// Else, we need to handle the error.
		return nil, err
	}

	// If reentry is allowed, check the entrypoints again
	if c.AllowReEntry {
		if next := checkHandlerList(c.EntryPoints, b, ctx); next != nil {
			return next, nil
		}
	}

	// else, check state mappings
	if next := checkHandlerList(c.States[currState], b, ctx); next != nil {
		return next, nil
	}

	// TODO: do we check fallbacks BEFORE or AFTER the main states?
	// else, fallbacks -> handle any cancellations
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
