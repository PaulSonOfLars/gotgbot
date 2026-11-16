package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/widgetquery"
)

//go:embed widget.html
var widgetHTML string

type counterState struct {
	Count int `json:"count"`
}

var (
	counterMu     sync.Mutex
	counterValues = map[string]int{}
)

// This bot demonstrates how to use Telegram Bot API 10.4 HTML widgets.
// It sends an interactive counter widget via SendWidget, and handles
// button clicks using WidgetQuery handlers and AnswerWidgetQuery.
func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		Logger:      logger,
	})
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{Logger: logger})

	// /start command to send the counter widget
	dispatcher.AddHandler(handlers.NewCommand("start", sendCounterWidget))

	// Handle widget query button presses
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("counter:"), handleCounterQuery))

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

	updater.Idle()
}

func sendCounterWidget(b *gotgbot.Bot, ctx *ext.Context) error {
	initialData, _ := json.Marshal(counterState{Count: 0})
	_, err := b.SendWidget(ctx.EffectiveChat.Id, widgetHTML, &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(initialData),
	})
	return err
}

func handleCounterQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}

	key := getWidgetKey(wq)

	counterMu.Lock()
	count := counterValues[key]
	switch wq.Data {
	case "counter:inc":
		count++
	case "counter:dec":
		count--
	case "counter:reset":
		count = 0
	}
	counterValues[key] = count
	counterMu.Unlock()

	newState := counterState{Count: count}
	raw, err := json.Marshal(newState)
	if err != nil {
		return err
	}

	// Persist the updated state to the widget's chatData
	persistChatData(b, wq, string(raw))

	// Answer the widget query with the new state
	_, err = b.AnswerWidgetQuery(wq.Id, &gotgbot.AnswerWidgetQueryOpts{
		Data: string(raw),
	})
	return err
}

func getWidgetKey(wq *gotgbot.WidgetQuery) string {
	if wq.InlineMessageId != "" {
		return "inline:" + wq.InlineMessageId
	}
	if wq.Chat != nil {
		return fmt.Sprintf("msg:%d:%d", wq.Chat.Id, wq.MessageId)
	}
	return "user:" + strconv.FormatInt(wq.From.Id, 10)
}

func persistChatData(b *gotgbot.Bot, wq *gotgbot.WidgetQuery, data string) {
	opts := &gotgbot.EditMessageWidgetDataOpts{ChatData: data}
	if wq.InlineMessageId != "" {
		opts.InlineMessageId = wq.InlineMessageId
	} else if wq.Chat != nil && wq.MessageId != 0 {
		opts.ChatId = wq.Chat.Id
		opts.MessageId = wq.MessageId
	} else {
		return
	}
	_, _ = b.EditMessageWidgetData(opts)
}
