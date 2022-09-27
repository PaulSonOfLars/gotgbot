package handlers_test

import (
	"errors"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func TestBasicConversation(t *testing.T) {
	b := NewTestBot()

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
		nil,
	)

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startCommand := NewCommandMessage(userId, chatId, "start", []string{})
	runHandler(t, b, &conv, startCommand, "", nextStep)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(userId, chatId, "message")
	runHandler(t, b, &conv, textMessage, nextStep, "")
	if !ended {
		t.Fatalf("expected the internal handler to have run")
	}

	// Ensure conversation has ended.
	checkExpectedState(t, &conv, textMessage, "")
}

func TestBasicKeyedConversation(t *testing.T) {
	b := NewTestBot()

	const nextStep = "nextStep"

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			return handlers.NextConversationState(nextStep)
		})},
		map[string][]ext.Handler{
			nextStep: {handlers.NewMessage(message.Contains("message"), func(b *gotgbot.Bot, ctx *ext.Context) error {
				return handlers.EndConversation()
			})},
		},
		&handlers.ConversationOpts{
			// Make sure that we key by sender in one chat
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySender),
		},
	)

	var userIdOne int64 = 123
	var userIdTwo int64 = 456
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startFromUserOne := NewCommandMessage(userIdOne, chatId, "start", []string{})
	messageFromTwo := NewMessage(userIdTwo, chatId, "message")

	runHandler(t, b, &conv, startFromUserOne, "", nextStep)

	// We have now started a conversation with user one
	checkExpectedState(t, &conv, startFromUserOne, nextStep)

	// But user two doesnt exist
	checkExpectedState(t, &conv, messageFromTwo, "")
}

func TestFallbackConversation(t *testing.T) {
	b := NewTestBot()

	const nextStep = "nextStep"
	var started bool
	var fallback bool

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			started = true
			return handlers.NextConversationState(nextStep)
		})},
		map[string][]ext.Handler{
			nextStep: {handlers.NewMessage(message.Contains("message"), func(b *gotgbot.Bot, ctx *ext.Context) error {
				t.Fatalf("internal handler should not have run")
				return handlers.EndConversation()
			})},
		},
		&handlers.ConversationOpts{
			Exits: []ext.Handler{handlers.NewCommand("cancel", func(b *gotgbot.Bot, ctx *ext.Context) error {
				fallback = true
				return handlers.EndConversation()
			})},
		},
	)

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startCommand := NewCommandMessage(userId, chatId, "start", []string{})
	runHandler(t, b, &conv, startCommand, "", nextStep)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "cancel" command, triggering the fallback handler (and causing it to "end").
	cancelCommand := NewCommandMessage(userId, chatId, "cancel", []string{})
	runHandler(t, b, &conv, cancelCommand, nextStep, "")
	if !fallback {
		t.Fatalf("expected the fallback handler to have run")
	}

	// Ensure conversation has ended.
	checkExpectedState(t, &conv, cancelCommand, "")
}
func TestReEntryConversation(t *testing.T) {
	b := NewTestBot()

	const nextStep = "nextStep"
	startCount := 0

	conv := handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			startCount++
			return handlers.NextConversationState(nextStep)
		})},
		map[string][]ext.Handler{
			nextStep: {handlers.NewMessage(message.Contains("message"), func(b *gotgbot.Bot, ctx *ext.Context) error {
				// We don't want this to happen!
				t.Fatalf("internal handler should not have run")
				return handlers.EndConversation()
			})},
		},
		nil,
	)
	// Enable reentry
	conv.AllowReEntry = true

	var userId int64 = 123
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startCommand := NewCommandMessage(userId, chatId, "start", []string{})
	runHandler(t, b, &conv, startCommand, "", nextStep)
	if startCount != 1 {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Send a message which matches both the entrypoint, and the "nextStep" state.
	cancelCommand := NewCommandMessage(userId, chatId, "start", []string{"message"})
	runHandler(t, b, &conv, cancelCommand, nextStep, nextStep) // Should hit
	if startCount != 2 {
		t.Fatalf("expected the entrypoint handler to have run a second time")
	}

	checkExpectedState(t, &conv, cancelCommand, nextStep)
}

func TestNestedConversation(t *testing.T) {
	b := NewTestBot()

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
		},
		nil,
	)

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
	start := NewCommandMessage(userId, chatId, startCmd, []string{})
	runHandler(t, b, &conv, start, "", firstStep)

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(userId, chatId, messageText)
	runHandler(t, b, &conv, textMessage, firstStep, secondStep)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedStart := NewCommandMessage(userId, chatId, nestedStartCmd, []string{})
	willRunHandler(t, b, &nestedConv, nestedStart, "")
	runHandler(t, b, &conv, nestedStart, secondStep, secondStep)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedFinish := NewMessage(userId, chatId, finishNestedText)
	willRunHandler(t, b, &nestedConv, nestedFinish, nestedStep)
	runHandler(t, b, &conv, nestedFinish, secondStep, thirdStep)

	// Ensure nested conversation has ended.
	checkExpectedState(t, &nestedConv, nestedFinish, "")
	t.Log("Nested conversation finished")

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	finish := NewMessage(userId, chatId, finishText)
	runHandler(t, b, &conv, finish, thirdStep, "")

	checkExpectedState(t, &conv, textMessage, "")
}

func runHandler(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, message *ext.Context, currentState string, nextState string) {
	willRunHandler(t, b, conv, message, currentState)
	if err := conv.HandleUpdate(b, message); err != nil {
		t.Fatalf("unexpected error from handler: %s", err.Error())
	}

	checkExpectedState(t, conv, message, nextState)
}

func willRunHandler(t *testing.T, b *gotgbot.Bot, conv *handlers.Conversation, message *ext.Context, expectedState string) {
	t.Logf("conv %p: checking message for %d in %d with text: %s", conv, message.EffectiveSender.Id(), message.EffectiveChat.Id, message.Message.Text)

	checkExpectedState(t, conv, message, expectedState)

	if !conv.CheckUpdate(b, message) {
		t.Fatalf("conv %p: expected the handler to match text: %s", conv, message.Message.Text)
	}
}

func checkExpectedState(t *testing.T, conv *handlers.Conversation, message *ext.Context, nextState string) {
	currentState, err := conv.StateStorage.Get(message)
	if nextState == "" {
		if !errors.Is(err, conversation.KeyNotFound) {
			t.Fatalf("should not have a conversation, but got currentState: %s", currentState)
		}
	} else if err != nil {
		t.Fatalf("unexpected error while checking the current currentState of the conversation")
	} else if currentState == nil || currentState.Key != nextState {
		t.Fatalf("expected the conversation to be at '%s', was '%s'", nextState, currentState)
	}
}
