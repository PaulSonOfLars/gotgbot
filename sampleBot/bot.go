package main

import (
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

func main() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel))
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()

	l.Info("Starting gotgbot...")
	token := os.Getenv("TOKEN")
	l.Info("token: ", token)
	updater, err := gotgbot.NewUpdater(token, logger)
	if err != nil {
		l.Fatalw("failed to start updater", zap.Error(err))
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
			l.Fatalw("Failed to start bot", zap.Error(err))
		}
		if !ok {
			l.Fatalw("Failed to set webhook", zap.Error(err))
		}
	} else {
		err := updater.StartPolling()
		if err != nil {
			l.Fatalw("Failed to start polling", zap.Error(err))
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

	_, err := u.EffectiveMessage.ReplyText(u.EffectiveMessage.ReplyToMessage.OriginalTextV2())
	if err != nil {
		b.Logger.Warnw("Error sending text", zap.Error(err))
	}

	_, err = u.EffectiveMessage.ReplyHTML(u.EffectiveMessage.ReplyToMessage.OriginalHTML())
	if err != nil {
		b.Logger.Warnw("Error sending HTML", zap.Error(err))
	}

	m := b.NewSendableMessage(u.EffectiveChat.Id, u.EffectiveMessage.ReplyToMessage.OriginalTextV2())
	m.ParseMode = parsemode.MarkdownV2
	_, err = m.Send()
	if err != nil {
		b.Logger.Warnw("Error sending V2", zap.Error(err))
	}
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

	u.Bot.Logger.Info(time.Since(t), "to send", lim, "updates.")
	time.Sleep(1)
}
