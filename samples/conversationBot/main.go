package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// This bot demonstrates a basic conversation handler. Conversation handlers help you to track states across multiple
// messages.
// In this case, the bot has commands to start a conversation, which then causes it to ask for your name and your age.
func main() {
	// Get token from the environment variable.
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("start", start)},
		map[string][]ext.Handler{
			NAME: {handlers.NewMessage(noCommands, name)},
			AGE:  {handlers.NewMessage(noCommands, age)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// These are the conversation handler states.
const (
	NAME = "name"
	AGE  = "age"
)

// Create a matcher which only matches text which is not a command
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

// start introduces the bot and starts the conversation
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nWhat is your name?.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return handlers.NextConversationState(NAME)
}

// cancel cancels the conversation.
func cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Oh, goodbye!", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}

// name gets the user's name
func name(b *gotgbot.Bot, ctx *ext.Context) error {
	inputName := ctx.EffectiveMessage.Text
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Nice to meet you, %s!\n\nAnd how old are you?", html.EscapeString(inputName)), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}
	return handlers.NextConversationState(AGE)
}

// age gets the user's age
func age(b *gotgbot.Bot, ctx *ext.Context) error {
	inputAge := ctx.EffectiveMessage.Text
	ageNumber, err := strconv.ParseInt(inputAge, 10, 64)
	if err != nil {
		// If the number is not valid, try again!
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("This doesn't seem to be a number. Could you repeat?"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		// We try the age handler again
		return handlers.NextConversationState(AGE)
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Ah, you're %d years old!", ageNumber), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send age message: %w", err)
	}
	return handlers.EndConversation()
}
