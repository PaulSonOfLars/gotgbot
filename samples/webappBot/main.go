package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// This bot shows how to use this library to server a webapp.
// Webapps are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
// For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
// Simply install ngrok, make an account on the website, and run:
// `ngrok http 8080`
// Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
// from the samples/webappBot directory:
// `URL="<your_url_here>" TOKEN="<your_token_here>" go run .`
// Then, simply send /start to your bot, and enjoy your webapp demo.
//
// This example also demonstrates how to use the updater's handler in a user-provided server.
func main() {
	// Get token from the environment variable
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// This MUST be an HTTPS URL for telegram to accept it.
	webappURL := os.Getenv("URL")
	if webappURL == "" {
		panic("URL environment variable is empty")
	}

	// Get the webhook secret from the environment variable.
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		panic("WEBHOOK_SECRET environment variable is empty")
	}

	// Create our bot.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher to handle updates in a simple manner.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot and send the URL
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		// We can wrap commands with anonymous functions to pass in extra variables, like the webapp URL, or other
		// configuration.
		return start(b, ctx, webappURL)
	}))

	// We add the bot webhook to our updater, such that we can populate the updater's http.Handler.
	err = updater.AddWebhook(b, b.Token, &ext.AddWebhookOpts{SecretToken: webhookSecret})
	if err != nil {
		panic("Failed to add bot webhooks to updater: " + err.Error())
	}

	// We select a subpath to specify where the updater handler is found on the http.Server.
	updaterSubpath := "/bots/"
	err = updater.SetAllBotWebhooks(webappURL+updaterSubpath, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookSecret,
	})
	if err != nil {
		panic("Failed to set bot webhooks: " + err.Error())
	}

	// Setup new HTTP server mux to handle different paths.
	mux := http.NewServeMux()
	// This serves the home page.
	mux.HandleFunc("/", index(webappURL))
	// This serves our "validation" API, which checks if the input data is valid.
	mux.HandleFunc("/validate", validate(token))
	// This serves the updater's webhook handler.
	mux.HandleFunc(updaterSubpath, updater.GetHandlerFunc(updaterSubpath))

	server := http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8080",
	}

	log.Printf("%s has been started...\n", b.User.Username)
	// Start the webserver displaying the page.
	// Note: ListenAndServe is a blocking operation, so we don't need to call updater.Idle() here.
	if err := server.ListenAndServe(); err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context, webappURL string) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nYou can use me to run a (very) simple telegram webapp demo!", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Press me", WebApp: &gotgbot.WebAppInfo{Url: webappURL}},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
