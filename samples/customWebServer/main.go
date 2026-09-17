package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// This example demonstrates how to create a bot with a custom webhook server.
// This is useful when you want to have more control over the server, or
// when you just want to use a different server library. In this example,
// we will use the gin-gonic/gin library to create a server.
func main() {
	// Get token from the environment variable.
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Get the webhook domain from the environment variable.
	webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
	if webhookDomain == "" {
		panic("WEBHOOK_DOMAIN environment variable is empty")
	}
	// Get the webhook secret from the environment variable.
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		panic("WEBHOOK_SECRET environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// Add webhook for this bot. Supposed our webhookDomain is
	// "http://localhost:8080", we will set webhook to
	// "http://localhost:8080/bots/<botToken>", so that we can
	// easily differentiate between multiple bots.
	err = updater.AddWebhook(b, b.Token, &ext.AddWebhookOpts{
		SecretToken: webhookSecret,
	})
	if err != nil {
		panic("failed to add webhook: " + err.Error())
	}

	// Start the webhook server.
	StartServer(updater, webhookDomain)

	// Add echo handler to reply to all text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

	err = updater.SetAllBotWebhooks(webhookDomain+"/bots", &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookSecret,
	})
	if err != nil {
		panic("failed to set webhook: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
