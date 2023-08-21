package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
)

// This bot is as basic as inline query bots can get. It simply links the bot library every time.
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

	if !b.SupportsInlineQueries {
		panic("bot does not support inline queries - enable them in botfather first!")
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
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

	// Create an inline query handler to reply to all inline queries
	dispatcher.AddHandler(handlers.NewInlineQuery(inlinequery.All, source))

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

// source is a very boring inline query handler that always just responds with a link to the library.
func source(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.InlineQuery.Answer(b, []gotgbot.InlineQueryResult{gotgbot.InlineQueryResultArticle{
		Id:      strconv.Itoa(rand.Int()),
		Title:   "Bot Library",
		Url:     "github.com/PaulSonOfLars/gotgbot",
		HideUrl: true,
		InputMessageContent: gotgbot.InputTextMessageContent{
			MessageText: "Bot library source code:\ngithub.com/PaulSonOfLars/gotgbot",
		},
		Description: "Link to the bot source code",
	}}, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
	})
	if err != nil {
		return fmt.Errorf("failed to send source ILQ: %w", err)
	}
	return nil
}
