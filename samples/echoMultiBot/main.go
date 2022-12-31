package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func main() {
	// This example defines how to define an echo bot that works for multiple bot instances at once.
	// It also shows how to stop the bot gracefully using the Updater.Stop() mechanism.

	// Get comma separated bot tokens from environment variable
	tokens := os.Getenv("TOKENS")
	if tokens == "" {
		panic("TOKENS environment variable is empty")
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fmt.Println("an error occurred while handling update:", err.Error())
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
		return stop(b, ctx, &updater)
	}))
	// Add echo handler to reply to all other text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

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

		// Start polling for the current bot.
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
			panic(fmt.Errorf("bot %d: failed to start polling: %w", idx, err))
		}

		fmt.Printf("bot %d: %s has been started...\n", idx, b.User.Username)
	}

	// Idle, to keep updates coming in, and thus avoid our bots from stopping.
	fmt.Println("Idling to keep main thread active, while incoming updates are handled.")
	updater.Idle()

	// If we get here, the updater.Idle() has ended.
	// This means the bots have been gracefully stopped via updater.Stop()
	fmt.Println("All bots have been gracefully stopped.")
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
