package handlers_test

import (
	"errors"
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
	runHandler(t, b, &conv, startCommand, "", nextStep)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(updateIdCounter, userId, chatId, "message")
	runHandler(t, b, &conv, textMessage, nextStep, "")
	if !ended {
		t.Fatalf("expected the internal handler to have run")
	}

	// Ensure conversation has ended.
	if _, err := conv.CurrentState(textMessage); err == nil || !errors.Is(err, handlers.ConversationKeyNotFound) {
		t.Fatalf("expected the conversation to be finished")
	}
}

func TestBasicKeyedConversation(t *testing.T) {
	b := NewTestBot()
	updateIdCounter := rand.Int63()

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
		[]ext.Handler{},
	)
	// Make sure that we key by sender in one chat
	conv.KeyStrategy = handlers.KeyStrategySender

	var userIdOne int64 = 123
	var userIdTwo int64 = 456
	var chatId int64 = 1234

	// Emulate sending the "start" command, triggering the entrypoint.
	startFromUserOne := NewCommandMessage(updateIdCounter, userIdOne, chatId, "start", []string{})
	messageFromTwo := NewMessage(updateIdCounter, userIdTwo, chatId, "message")

	runHandler(t, b, &conv, startFromUserOne, "", nextStep)

	// We have now started a conversation with user one
	s, err := conv.CurrentState(startFromUserOne)
	if err != nil {
		t.Fatalf("%d should now have started the conversation and be at %s", userIdOne, nextStep)
	}
	if s != nextStep {
		t.Fatalf("%d should now have moved state to %s, but was %s", userIdOne, nextStep, s)
	}

	// But user two doesnt exist
	s2, err := conv.CurrentState(messageFromTwo)
	if err == nil {
		t.Fatalf("%d should not have a conversation at this point, but got %s", userIdTwo, s2)
	}
	if !errors.Is(err, handlers.ConversationKeyNotFound) {
		t.Fatalf("%d should not have a conversation at this point", userIdTwo)
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
	runHandler(t, b, &conv, startCommand, "", nextStep)
	if !started {
		t.Fatalf("expected the entrypoint handler to have run")
	}

	// Emulate sending the "cancel" command, triggering the fallback handler (and causing it to "end").
	cancelCommand := NewCommandMessage(updateIdCounter, userId, chatId, "cancel", []string{})
	runHandler(t, b, &conv, cancelCommand, nextStep, "")
	if !fallback {
		t.Fatalf("expected the fallback handler to have run")
	}
	if internal {
		t.Fatalf("did not expect the internal handler to have run")
	}

	// Ensure conversation has ended.
	if _, err := conv.CurrentState(cancelCommand); err == nil || !errors.Is(err, handlers.ConversationKeyNotFound) {
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
	runHandler(t, b, &conv, start, "", firstStep)

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	textMessage := NewMessage(updateIdCounter, userId, chatId, messageText)
	runHandler(t, b, &conv, textMessage, firstStep, secondStep)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedStart := NewCommandMessage(updateIdCounter, userId, chatId, nestedStartCmd, []string{})
	willRunHandler(t, b, &nestedConv, nestedStart, "")
	runHandler(t, b, &conv, nestedStart, secondStep, secondStep)

	// Emulate sending the "nested_start" command, triggering the entrypoint of the nested conversation.
	nestedFinish := NewMessage(updateIdCounter, userId, chatId, finishNestedText)
	willRunHandler(t, b, &nestedConv, nestedFinish, nestedStep)
	runHandler(t, b, &conv, nestedFinish, secondStep, thirdStep)

	// Ensure nested conversation has ended.
	if _, err := nestedConv.CurrentState(nestedFinish); err == nil || !errors.Is(err, handlers.ConversationKeyNotFound) {
		t.Fatalf("expected the nested conversation to be finished")
	}
	t.Log("Nested conversation finished")

	// Emulate sending the "message" text, triggering the internal handler (and causing it to "end").
	finish := NewMessage(updateIdCounter, userId, chatId, finishText)
	runHandler(t, b, &conv, finish, thirdStep, "")

	// Ensure conversation has ended.
	if _, err := conv.CurrentState(textMessage); err == nil || !errors.Is(err, handlers.ConversationKeyNotFound) {
		t.Fatalf("expected the conversation to be finished")
	}
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
	state, err := conv.CurrentState(message)
	if nextState == "" {
		if !errors.Is(err, handlers.ConversationKeyNotFound) {
			t.Fatalf("should not have a conversation, but got state: %s", state)
		}
	} else if err != nil {
		t.Fatalf("unexpected error while checking the current state of the conversation")
	} else if state != nextState {
		t.Fatalf("expected the conversation to be at '%s', was '%s'", nextState, state)
	}
}
