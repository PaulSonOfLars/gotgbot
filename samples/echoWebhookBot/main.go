package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// This bot repeats everything you say - but it uses webhooks instead of long polling.
// Webhooks are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
// For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
// Simply install ngrok, make an account on the website, and run:
// `ngrok http 8080`
// Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
// from the samples/echoWebhookBot directory:
// `TOKEN="<your_token_here>" WEBHOOK_DOMAIN="<your_domain_here>"  WEBHOOK_SECRET="<random_string_here>" go run .`
// Then, simply send /start to your bot; if it replies, you've successfully set up webhooks!
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

	// Add echo handler to reply to all text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

	// Start the webhook server. We start the server before we set the webhook itself, so that when telegram starts
	// sending updates, the server is already ready.
	webhookOpts := ext.WebhookOpts{
		Listen:      "0.0.0.0", // This example assumes you're in a dev environment running ngrok on 8080.
		Port:        8080,
		SecretToken: webhookSecret, // Setting a webhook secret here allows you to ensure the webhook is set by you (must be set here AND in SetWebhook!).
	}
	// We use the token as the urlPath for the webhook, as using a secret ensures that strangers aren't crafting fake updates.
	err = updater.StartWebhook(b, token, webhookOpts)
	if err != nil {
		panic("failed to start webhook: " + err.Error())
	}

	err = updater.SetAllBotWebhooks(webhookDomain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
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
