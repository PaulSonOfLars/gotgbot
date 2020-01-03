package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

func main() {
	log.Println("Starting gotgbot...")
	token := os.Getenv("TOKEN")
	log.Println("token:", token)
	updater, err := gotgbot.NewUpdater(token)
	if err != nil {
		log.Fatal(err)
	}
	// reply to /start messages
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", start))
	// get the message HTML
	updater.Dispatcher.AddHandler(handlers.NewCommand("get", get))
	// reply to messages satisfying this regex
	updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)hello", hi))
	// reply to all messages satisfying the filter
	updater.Dispatcher.AddHandler(handlers.NewMessage(Filters.Sticker, stickerDeleter))

	if os.Getenv("USE_WEBHOOKS") == "t" {
		// start getting updates
		webhook := gotgbot.Webhook{
			Serve:          "0.0.0.0",
			ServePort:      8080,
			ServePath:      updater.Bot.Token,
			URL:            os.Getenv("WEBHOOK_URL"),
			MaxConnections: 30,
		}
		updater.StartWebhook(webhook)
		ok, err := updater.SetWebhook(updater.Bot.Token, webhook)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to start bot due to: ", err)
		}
		if !ok {
			logrus.Fatal("Failed to set webhook")
		}
	} else {
		err := updater.StartPolling()
		if err != nil {
			logrus.WithError(err).Fatal("failed to start polling")
		}
	}

	// wait
	updater.Idle()
}

// return the HTML
func get(b ext.Bot, u *gotgbot.Update) error {
	if u.EffectiveMessage.ReplyToMessage == nil {
		u.EffectiveMessage.ReplyText("Please reply to a message!")
		return nil
	}

	u.EffectiveMessage.ReplyText(u.EffectiveMessage.ReplyToMessage.OriginalHTML())
	return nil
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

// can be used to run local load tests
func doLocalLoadTest(u *gotgbot.Updater) {
	lim := 30000
	text := "hi"

	ups := make([]gotgbot.RawUpdate, lim)
	for i := 0; i < lim; i++ {
		x, _ := json.Marshal(gotgbot.Update{
			UpdateId: 1 + i,
			Message: &ext.Message{
				MessageId: 1 + i,
				From: &ext.User{
					Id:           666,
					IsBot:        false,
					FirstName:    "test",
					LanguageCode: "en",
				},
				Chat: &ext.Chat{
					Id:        666,
					FirstName: "testchat",
					Type:      "private",
				},
				Text: text,
				Date: int(time.Now().Unix()),
			},
		})
		ups[i] = x
	}

	t := time.Now()

	for _, x := range ups {
		u.Updates <- &x
	}

	logrus.Println(time.Since(t), "to send", lim, "updates.")
	time.Sleep(1)
}
