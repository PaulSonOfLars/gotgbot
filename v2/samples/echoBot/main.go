package main

import (
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(b)
	dispatcher := updater.Dispatcher

	// Add echo handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewMessage(filters.All, echo))

	// Start receiving updates.
	updater.StartCleanPolling(b)
	fmt.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func echo(ctx *ext.Context) error {
	// Reply to message with its own contents
	ctx.EffectiveMessage.Reply(ctx.Bot, ctx.EffectiveMessage.Text, gotgbot.SendMessageOpts{})
	return nil
}
