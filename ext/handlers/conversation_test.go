package handlers_test

import (
	"math/rand"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func TestBasicConversation(t *testing.T) {
	b := NewTestBot()
	updateIdCounter := rand.Int63()

	const nextStep = "nextStep"
	var started bool
	var ended bool

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			started = true
			return handlers.NextConversationState(nextStep)
		})},
		map[string][]ext.Handler{
			nextStep: {handlers.NewMessage(message.Contains("message"), func(b *gotgbot.Bot, ctx *ext.Context) error {
				ended = true
				return handlers.EndConversation()
			})},
		},
		[]ext.Handler{},
	)

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startCommand := NewCommandMessage(updateIdCounter, userId, chatId, "start", []string{})
	checkEntrypointHasRun(t, b, conv, startCommand, &started)

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(updateIdCounter, userId, chatId, "message")
	checkInternalHandlerHasRun(t, b, conv, nextStep, textMessage, &ended)

	// Ensure conversation has ended.
	if _, ok := conv.CurrentState(textMessage); ok {
		t.Errorf("expected the conversation to be finished")
	}
}

func checkEntrypointHasRun(t *testing.T, b *gotgbot.Bot, conv handlers.Conversation, startCommand *ext.Context, started *bool) {
	if _, ok := conv.CurrentState(startCommand); ok {
		t.Errorf("expected the conversation to not be started")
	}
	if !conv.CheckUpdate(b, startCommand) {
		t.Errorf("expected the entrypoint handler to match")
	}
	if err := conv.HandleUpdate(b, startCommand); err != nil {
		t.Errorf("unexpected error from update handling entrypoint")
	}
	if !*started {
		t.Errorf("expected the entrypoiny handler to have run")
	}
}

func checkInternalHandlerHasRun(t *testing.T, b *gotgbot.Bot, conv handlers.Conversation, expectedState string, textMessage *ext.Context, ended *bool) {
	if state, ok := conv.CurrentState(textMessage); !ok || state != expectedState {
		t.Errorf("expected the conversation to be at '%s', was '%s'", expectedState, state)
	}
	if !conv.CheckUpdate(b, textMessage) {
		t.Errorf("expected the internal handler to match")
	}
	if err := conv.HandleUpdate(b, textMessage); err != nil {
		t.Errorf("unexpected error from update handling internal handler")
	}
	if !*ended {
		t.Errorf("expected the interna; handler to have run")
	}
}
