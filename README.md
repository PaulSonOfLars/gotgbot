# Go Telegram Bot

**This library is WIP; it does not currently support all of the telegram api methods.**

This library attempts to create a user-friendly wrapper around the telegram bot api.

Heavily inspired by the [python-telegram-bot library](github.com/python-telegram-bot/python-telegram-bot),
this aims to create a simple way to manage a concurrent and scalable bot.

## Getting started
Install it as you would install your usual go library: `go get github.com/PaulSonOfLars/gotgbot`

A sample bot would look something like this:

```
func main() {
	log.Println("Starting gotgbot...")
	updater := gotgbot.NewUpdater(YOUR_TOKEN_HERE)
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

func start(b ext.Bot, u gotgbot.Update) {
	b.SendMessage(u.Message.Chat.Id, "Congrats! You just issued a /start on your go bot!")
}

func hi(b ext.Bot, u gotgbot.Update) {
	b.SendMessage(u.Message.Chat.Id, "Hello to you too!")
}

func stickerDeleter(b ext.Bot, u gotgbot.Update) {
	if _, err := u.EffectiveMessage.Delete(); err != nil {
		u.EffectiveMessage.ReplyMessage("Can't delete, you're in PM!")
	} else {
		msg := b.NewSendableMessage(u.Message.Chat.Id, "Don't you *dare* send _stickers_ here!")
		msg.ParseMode = parsemode.Markdown
		msg.Send()
	}
}
```


An interesting feature to take note of is that due to go's
handling of exceptions, if you choose not to handle an exception, your bot
will simply keep on going happily and ignore any issues.

All handlers are async; theyre all executed in their own go routine,
so can communicate accross channels if needed.

## Message sending

As seen in the example, message sending can be done in two ways; via each received message's
ReplyMessage() function, or by building your own; and calling msg.Send(). This allows for
ease of use by having the most commonly used shortcuts readily available, while
retaining the flexibility of building each message yourself, which wouldnt be
available otherwise.

