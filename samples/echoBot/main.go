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
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Add echo handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("source", source))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("start_callback"), startCB))
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
		fmt.Println("failed to open source: " + err.Error())
		return nil
	}

	_, err = b.SendDocument(ctx.EffectiveChat.Id, f, &gotgbot.SendDocumentOpts{
		Caption:          "Here is my source code.",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		fmt.Println("failed to send source: " + err.Error())
		return nil
	}

	// Alternative file sending solutions:

	// --- By file_id:
	// _, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, "file_id", &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	// })
	// if err != nil {
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	// }

	// --- By []byte:
	// bs, err := ioutil.ReadFile("samples/echoBot/main.go")
	// if err != nil {
	//	fmt.Println("failed to open source: " + err.Error())
	//	return nil
	// }
	//
	// _, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, bs, &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	// })
	// if err != nil {
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	// }

	// --- By custom name:
	// f2, err := os.Open("samples/echoBot/main.go")
	// if err != nil {
	//	fmt.Println("failed to open source: " + err.Error())
	//	return nil
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
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	// }

	return nil
}

// start introduces the bot
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
		fmt.Println("failed to send: " + err.Error())
	}
	return nil
}

// startCB edits the start message
func startCB(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	cb.Message.EditText(b, "You edited the start message.", nil)
	return nil
}

// echo replies to a messages with its own contents
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	return nil
}
