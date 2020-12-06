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
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	updater := ext.NewUpdater(b)
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewMessage(filters.All, test))

	updater.StartCleanPolling(b)
	fmt.Printf("%s has been started...\n", b.User.Username)
	updater.Idle()
}

func test(ctx *ext.Context) error {
	ctx.Update.Message.Reply(ctx.Bot, ctx.Update.Message.Text, gotgbot.SendMessageOpts{})
	return nil
}
