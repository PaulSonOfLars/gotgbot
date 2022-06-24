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
	runEntrypoint(t, b, &conv, startCommand)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(updateIdCounter, userId, chatId, "message")
	runInternalHandler(t, b, &conv, nextStep, textMessage)
	if !ended {
		t.Fatalf("expected the internal handler to have run")
	}

	// Ensure conversation has ended.
	if _, ok := conv.CurrentState(textMessage); ok {
		t.Fatalf("expected the conversation to be finished")
	}
}
func TestFallbackConversation(t *testing.T) {
	b := NewTestBot()
	updateIdCounter := rand.Int63()

	const nextStep = "nextStep"
	var started bool
	var internal bool
	var fallback bool

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			started = true
			return handlers.NextConversationState(nextStep)
		})},
		map[string][]ext.Handler{
			nextStep: {handlers.NewMessage(message.Contains("message"), func(b *gotgbot.Bot, ctx *ext.Context) error {
				internal = true
				return handlers.EndConversation()
			})},
		},
		[]ext.Handler{handlers.NewCommand("cancel", func(b *gotgbot.Bot, ctx *ext.Context) error {
			fallback = true
			return handlers.EndConversation()
		})},
	)

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startCommand := NewCommandMessage(updateIdCounter, userId, chatId, "start", []string{})
	runEntrypoint(t, b, &conv, startCommand)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "cancel" command, triggering the fallback handler (and causing it to "end").
	cancelCommand := NewCommandMessage(updateIdCounter, userId, chatId, "cancel", []string{})
	runInternalHandler(t, b, &conv, nextStep, cancelCommand)
	if !fallback {
		t.Fatalf("expected the fallback handler to have run")
	}
	if internal {
		t.Fatalf("did not expect the internal handler to have run")
	}

	// Ensure conversation has ended.
	if _, ok := conv.CurrentState(cancelCommand); ok {
		t.Fatalf("expected the conversation to be finished")
	}
}

func TestNestedConversation(t *testing.T) {
	b := NewTestBot()
	updateIdCounter := rand.Int63()

	const firstStep = "firstStep"
	const secondStep = "secondStep"
	const nestedStep = "nestedStep"
	const thirdStep = "thirdStep"

	const startCmd = "start"
	const nestedStartCmd = "nested_start"
	const messageText = "message"
	const finishNestedText = "finish nested"
	const finishText = "finish"

	nestedConv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand(nestedStartCmd, func(b *gotgbot.Bot, ctx *ext.Context) error {
			return handlers.NextConversationState(nestedStep)
		})},
		map[string][]ext.Handler{
			nestedStep: {handlers.NewMessage(message.Contains(finishNestedText), func(b *gotgbot.Bot, ctx *ext.Context) error {
				return handlers.EndConversationToParentState(handlers.NextConversationState(thirdStep))
			})},
		}, nil)

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand(startCmd, func(b *gotgbot.Bot, ctx *ext.Context) error {
			return handlers.NextConversationState(firstStep)
		})},
		map[string][]ext.Handler{
			firstStep: {handlers.NewMessage(message.Contains(messageText), func(b *gotgbot.Bot, ctx *ext.Context) error {
				return handlers.NextConversationState(secondStep)
			})},
			secondStep: {nestedConv},
			thirdStep: {handlers.NewMessage(message.Contains(finishText), func(b *gotgbot.Bot, ctx *ext.Context) error {
				return handlers.EndConversation()
			})},
		},
		nil,
	)

	t.Logf("main   conv: %p", &conv)
	t.Logf("nested conv: %p", &nestedConv)

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	start := NewCommandMessage(updateIdCounter, userId, chatId, startCmd, []string{})
	runEntrypoint(t, b, &conv, start)

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(updateIdCounter, userId, chatId, messageText)
	runInternalHandler(t, b, &conv, firstStep, textMessage)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedStart := NewCommandMessage(updateIdCounter, userId, chatId, nestedStartCmd, []string{})
	willRunEntrypoint(t, b, &nestedConv, nestedStart)
	runInternalHandler(t, b, &conv, secondStep, nestedStart)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedFinish := NewMessage(updateIdCounter, userId, chatId, finishNestedText)
	willRunInternalHandler(t, b, &nestedConv, nestedStep, nestedFinish)
	runInternalHandler(t, b, &conv, secondStep, nestedFinish)

	// Ensure nested conversation has ended.
	if _, ok := nestedConv.CurrentState(nestedFinish); ok {
		t.Fatalf("expected the nested conversation to be finished")
	}
	t.Log("Nested conversation finished")

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	finish := NewMessage(updateIdCounter, userId, chatId, finishText)
	runInternalHandler(t, b, &conv, thirdStep, finish)

	// Ensure conversation has ended.
	if _, ok := conv.CurrentState(textMessage); ok {
		t.Fatalf("expected the conversation to be finished")
	}
}

func runEntrypoint(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, message *ext.Context) {
	willRunEntrypoint(t, b, conv, message)
	if err := conv.HandleUpdate(b, message); err != nil {
		t.Fatalf("unexpected error from update handling entrypoint: %s", err.Error())
	}
}

func willRunEntrypoint(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, message *ext.Context) {
	t.Logf("conv %p: checking entrypoint message for %d in %d with text: %s", conv, message.EffectiveSender.Id(), message.EffectiveChat.Id, message.Message.Text)

	if _, ok := conv.CurrentState(message); ok {
		t.Fatalf("expected the conversation to not be started")
	}
	if !conv.CheckUpdate(b, message) {
		t.Fatalf("conv %p: expected the entrypoint handler to match text: %s", conv, message.Message.Text)
	}
}

func runInternalHandler(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, expectedState string, message *ext.Context) {
	willRunInternalHandler(t, b, conv, expectedState, message)
	if err := conv.HandleUpdate(b, message); err != nil {
		t.Fatalf("unexpected error from update handling internal handler: %s", err.Error())
	}
}

func willRunInternalHandler(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, expectedState string, message *ext.Context) {
	t.Logf("conv %p: checking internal message for %d in %d with text: %s", conv, message.EffectiveSender.Id(), message.EffectiveChat.Id, message.Message.Text)

	if state, ok := conv.CurrentState(message); !ok || state != expectedState {
		t.Fatalf("expected the conversation to be at '%s', was '%s'", expectedState, state)
	}
	if !conv.CheckUpdate(b, message) {
		t.Fatalf("conv %p: expected the internal handler to match text: %s", conv, message.Message.Text)
	}
}
