package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// This bot demonstrates how to pass around variables to all handlers without changing any function signatures.
// That is - by using a struct with methods instead of simple functions.
// This pattern is great to avoid passing data around through global variables. The client can store database clients,
// cache clients, in memory clients, and many more.
func main() {
	// Get token from the environment variable
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
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

	// We create our stateful bot client here.
	c := &client{}

	// Add the /start command.
	// Note the use of `c.start` rather than the usual `start`.
	dispatcher.AddHandler(handlers.NewCommand("start", c.start))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
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

// start counts how many times the start command has been used.
// This is set up as a method on the client struct, such that we can access shared fields across method calls.
// NOTE: Make sure to use a pointer receiver to avoid copying the client data!
func (c *client) start(b *gotgbot.Bot, ctx *ext.Context) error {
	// Get the existing data. "ok" will be false if the data doesn't exist yet.
	countVal, ok := c.getUserData(ctx, "count")
	if !ok {
		c.setUserData(ctx, "count", 1)
		ctx.EffectiveMessage.Reply(b, "This is the first time you press start.", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Cast the data to an int so it can be used.
	count, ok := countVal.(int)
	if !ok {
		ctx.EffectiveMessage.Reply(b, "'count' was not an integer, as was expected - this is a programmer error!", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Increment our count (one more press!)
	count += 1
	c.setUserData(ctx, "count", count)
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("You have pressed start %d times.", count), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
