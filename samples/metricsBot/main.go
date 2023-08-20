package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// This bot shows how to gather various prometheus metrics from the dispatcher and bot client, such that they can be
// scraped and placed on a dashboard.
// Note: this example is NOT a bot to gather useful chat metrics; it simply demonstrates how various bot metrics
// could be collected.
func main() {
	// Get token from the environment variable
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		// Setup botclient to collect metrics on client requests
		BotClient: newMetricsClient(),
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create the dispatcher with our custom processor.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// Collect metrics on the dispatcher's update processing
		Processor: metricsProcessor{processor: ext.BaseProcessor{}},
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	// Collect metrics on the state of the dispatcher's buffer.
	go monitorDispatcherBuffer(dispatcher)

	// Create the updater with our customised dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: dispatcher})

	// Add echo handler to reply to all text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

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

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
