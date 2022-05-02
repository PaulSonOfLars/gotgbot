package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"), &gotgbot.BotOpts{
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
		DispatcherOpts: ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fmt.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		},
	})
	dispatcher := updater.Dispatcher

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// Answer callback query sent in the /start command.
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("start_callback"), startCB))
	// /source command to send the bot source code
	dispatcher.AddHandler(handlers.NewCommand("source", source))
	// Add echo handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewMessage(message.All, echo))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func source(b *gotgbot.Bot, ctx *ext.Context) error {
	f, err := os.Open("samples/echoBot/main.go")
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
	// bs, err := ioutil.ReadFile("samples/echoBot/main.go")
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
	// f2, err := os.Open("samples/echoBot/main.go")
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
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Press me", CallbackData: "start_callback"},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

// startCB edits the start message.
func startCB(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "You pressed a button!",
	})
	if err != nil {
		return fmt.Errorf("failed to answer start callback query: %w", err)
	}

	_, _, err = cb.Message.EditText(b, "You edited the start message.", nil)
	if err != nil {
		return fmt.Errorf("failed to edit start message text: %w", err)
	}
	return nil
}

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
