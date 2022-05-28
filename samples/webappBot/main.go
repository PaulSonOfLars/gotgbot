package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// This bot is slightly more complex to run, since it requires a running webserver, as well as an HTTPS domain.
// For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
// Simply install ngrok, make an account on the website, and run:
// ngrok http 8080
// Then, copy paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
// from the samples/webappBot directory:
// URL="<your_url_here>" TOKEN="<your_token_here>" go run .
// Then, simply send /start to your bot, and enjoy your webapp demo.
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

	// Create our bot.
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

	// Create updater and dispatcher to handle updates in a simple manner.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fmt.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		},
	})
	dispatcher := updater.Dispatcher

	// /start command to introduce the bot and send the URL
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		return start(b, ctx, webappURL)
	}))

	// Start receiving (and handling) updates.
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)

	// Setup new HTTP server mux to handle different paths.
	mux := http.NewServeMux()
	// This serves the home page.
	mux.HandleFunc("/", index(webappURL))
	// This serves our "validation" API, which checks if the input data is valid.
	mux.HandleFunc("/validate", validate(token))
	server := http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8080",
	}

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
