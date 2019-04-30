package main

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

func main() {
	log.Println("Starting gotgbot...")
	updater, err := gotgbot.NewUpdater("YOUR_TOKEN_HERE")
	if err != nil {
		log.Fatal(err)
	}
	// reply to /start messages
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", start))
	// reply to messages satisfying this regex
	updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)hello", hi))
	// reply to all messages satisfying the filter
	updater.Dispatcher.AddHandler(handlers.NewMessage(Filters.Sticker, stickerDeleter))

	// start getting updates
	updater.StartPolling()

	// wait
	updater.Idle()
}

func start(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, "Congrats! You just issued a /start on your go bot!")
	return nil
}

func hi(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, "Hello to you too!")
	return gotgbot.ContinueGroups{} // will keep executing handlers, even after having been caught by this one.
}

func stickerDeleter(b ext.Bot, u *gotgbot.Update) error {
	if _, err := u.EffectiveMessage.Delete(); err != nil {
		u.EffectiveMessage.ReplyText("Can't delete, you're in PM!")
	} else {
		msg := b.NewSendableMessage(u.Message.Chat.Id, "Don't you *dare* send _stickers_ here!")
		msg.ParseMode = parsemode.Markdown
		msg.Send()
	}
	return nil
}
