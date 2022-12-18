package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// This bot demonstrates some example interactions with commands on telegram.
// It has a basic start command with a bot intro.
// It also has a source command, which sends the bot sourcecode, as a file.
func main() {
	// Get token from the environment variable
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Create bot from environment value.
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

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// /source command to send the bot source code
	dispatcher.AddHandler(handlers.NewCommand("source", source))

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

func source(b *gotgbot.Bot, ctx *ext.Context) error {
	f, err := os.Open("samples/commandBot/main.go")
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}

	_, err = b.SendDocument(ctx.EffectiveChat.Id, f, &gotgbot.SendDocumentOpts{
		Caption:          "Here is my source code.",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}

	// Alternative file sending solutions:

	// --- By file_id:
	// _, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, "file_id", &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	// })
	// if err != nil {
	//	return fmt.Errorf("failed to send source: %w", err)
	// }

	// --- By []byte:
	// bs, err := ioutil.ReadFile("samples/commandBot/main.go")
	// if err != nil {
	//	return fmt.Errorf("failed to open source: %w", err)
	// }
	//
	// _, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, bs, &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	// })
	// if err != nil {
	//	return fmt.Errorf("failed to send source: %w", err)
	// }

	// --- By custom name:
	// f2, err := os.Open("samples/commandBot/main.go")
	// if err != nil {
	//	return fmt.Errorf("failed to open source: %w", err)
	//	return err
	// }
	//
	// _, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, gotgbot.NamedFile{
	//	File:     f2,
	//	FileName: "NewFileName",
	// }, &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	// })
	// if err != nil {
	//	return fmt.Errorf("failed to send source: %w", err)
	//	return err
	// }

	return nil
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s. I <b>repeat</b> all your messages.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
