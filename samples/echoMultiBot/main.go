package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// This bot demonstrates how to create echo bot that works with multiple bot instances at once.
// It also shows how to stop the bot gracefully using the Updater.Stop() mechanism.
// It has options to use either polling or webhooks.
func main() {

	// Get comma separated bot tokens from environment variable
	tokens := os.Getenv("TOKENS")
	if tokens == "" {
		panic("TOKENS environment variable is empty")
	}

	// Get the webhook domain from the environment variable.
	webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
	if webhookDomain == "" {
		log.Println("no webhook domain specified; using long polling.")
	}
	// Get the webhook secret from the environment variable.
	webhookSecret := os.Getenv("WEBHOOK_SECRET")

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
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

	// Add stop handler to stop all bots gracefully.
	dispatcher.AddHandler(handlers.NewCommand("stop", func(b *gotgbot.Bot, ctx *ext.Context) error {
		// Using an anonymous function here can be a nice way of passing additional parameters into bot methods.
		// In this case we pass the updater - this will allow us to stop it and exit the program.
		return stop(b, ctx, updater)
	}))
	// Add echo handler to reply to all other text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

	var bots []*gotgbot.Bot
	// We iterate over all the tokens provided, to create and start polling on each one.
	for idx, token := range strings.Split(tokens, ",") {
		b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		})
		if err != nil {
			panic(fmt.Errorf("bot %d: failed to create bot: %w", idx, err))
		}
		bots = append(bots, b)
	}

	if webhookDomain != "" {
		err := startWebhookBots(updater, bots, webhookDomain, webhookSecret)
		if err != nil {
			panic("Failed to start bots via webhook: " + err.Error())
		}

	} else {
		err := startLongPollingBots(updater, bots)
		if err != nil {
			panic("Failed to start bots via polling: " + err.Error())
		}
	}

	// Idle, to keep updates coming in, and thus avoid our bots from stopping.
	log.Println("Idling to keep main thread active while incoming updates are handled.")
	updater.Idle()

	// If we get here, the updater.Idle() has ended.
	// This means the bots have been gracefully stopped via updater.Stop()
	log.Println("All bots have been gracefully stopped.")
}

// startLongPollingBots demonstrates how to start multiple bots with long-polling.
func startLongPollingBots(updater *ext.Updater, bots []*gotgbot.Bot) error {
	for idx, b := range bots {
		err := updater.StartPolling(b, &ext.PollingOpts{
			DropPendingUpdates: true,
			GetUpdatesOpts: gotgbot.GetUpdatesOpts{
				Timeout: 9,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 10,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to start polling for bot %d: %w", idx, err)
		}

		log.Printf("bot %d: %s has started polling...\n", idx, b.User.Username)
	}

	return nil
}

// startWebhookBots demonstrates how to start multiple bots with webhooks.
func startWebhookBots(updater *ext.Updater, bots []*gotgbot.Bot, domain string, webhookSecret string) error {
	opts := ext.WebhookOpts{
		Listen:      "0.0.0.0", // This example assumes you're in a dev environment running ngrok on 8080.
		Port:        8080,
		SecretToken: webhookSecret,
	}

	// We start the server before we set the webhooks, so that incoming requests can be processed immediately.
	err := updater.StartServer(opts)
	if err != nil {
		return fmt.Errorf("failed to start the webhook server: %w", err)
	}

	// We add all the bots to the updater.
	for idx, b := range bots {
		updater.AddWebhook(b, b.GetToken(), opts)
		log.Printf("bot %d: %s has been added to the updater\n", idx, b.User.Username)
	}

	// We set the webhooks for all the added bots, so telegram starts sending updates.
	return updater.SetAllBotWebhooks(domain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		AllowedUpdates:     nil,
		DropPendingUpdates: false,
		SecretToken:        webhookSecret,
	})
}

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

// echo replies to a messages with its own contents.
func stop(b *gotgbot.Bot, ctx *ext.Context, updater *ext.Updater) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Stopping all bots...", nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}

	// We stop the updater in a separate goroutine, otherwise it would be stuck waiting for itself.
	go func() {
		err = updater.Stop()
		if err != nil {
			ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Failed to stop updater: %s", err.Error()), nil)
			return
		}
	}()

	return nil
}
