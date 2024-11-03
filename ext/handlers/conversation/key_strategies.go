package conversation

import (
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var ErrEmptyKey = errors.New("empty conversation key")

// KeyStrategy is the function used to obtain the current key in the ongoing conversation.
//
// Use one of the existing keys, or define your own if you need external data (eg a DB or other state).
type KeyStrategy func(ctx *ext.Context) (string, error)

var (
	// Ensure key strategy methods match the function signatures.
	_ KeyStrategy = KeyStrategyChat
	_ KeyStrategy = KeyStrategySender
	_ KeyStrategy = KeyStrategySenderAndChat
)

// KeyStrategySenderAndChat ensures that each sender get a unique conversation, even in different chats.
func KeyStrategySenderAndChat(ctx *ext.Context) (string, error) {
	if ctx.EffectiveSender == nil || ctx.EffectiveChat == nil {
		return "", fmt.Errorf("missing sender or chat fields: %w", ErrEmptyKey)
	}
	return fmt.Sprintf("%d/%d/%d", ctx.Bot.Id, ctx.EffectiveSender.Id(), ctx.EffectiveChat.Id), nil
}

// KeyStrategySender gives a unique conversation to each sender, and that single conversation is available in all chats.
func KeyStrategySender(ctx *ext.Context) (string, error) {
	if ctx.EffectiveSender == nil {
		return "", fmt.Errorf("missing sender field: %w", ErrEmptyKey)
	}
	return fmt.Sprintf("%d/%d", ctx.Bot.Id, ctx.EffectiveSender.Id()), nil
}

// KeyStrategyChat gives a unique conversation to each chat, which all senders can interact in together.
func KeyStrategyChat(ctx *ext.Context) (string, error) {
	if ctx.EffectiveChat == nil {
		return "", fmt.Errorf("missing chat field: %w", ErrEmptyKey)
	}
	return fmt.Sprintf("%d/%d", ctx.Bot.Id, ctx.EffectiveChat.Id), nil
}

// StateKey provides a sane default for handling incoming updates.
func StateKey(ctx *ext.Context, strategy KeyStrategy) (string, error) {
	if strategy == nil {
		return KeyStrategySenderAndChat(ctx)
	}
	return strategy(ctx)
}
