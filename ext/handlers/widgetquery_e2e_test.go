package handlers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/sharedwidgetquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/widgetquery"
)

func TestBotEndToEnd_WidgetUpdates(t *testing.T) {
	var widgetQueryHandled bool
	var sharedWidgetQueryHandled bool
	var widgetMessageHandled bool
	var answeredQueryId string

	// Create test bot with mocked client
	b := &gotgbot.Bot{
		Token: "12345:mock_token",
		User: gotgbot.User{
			Id:        12345,
			IsBot:     true,
			FirstName: "WidgetTestBot",
			Username:  "WidgetTestBot",
		},
		BotClient: &mockBotClient{
			requestFunc: func(ctx context.Context, token string, method string, params map[string]any, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
				if method == "answerWidgetQuery" {
					if qId, ok := params["widget_query_id"].(string); ok {
						answeredQueryId = qId
					}
					return json.Marshal(true)
				}
				return json.Marshal(true)
			},
		},
	}

	// Create real dispatcher
	dispatcher := ext.NewDispatcher(nil)

	// 1. Add WidgetQuery handler
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("click_"), func(bot *gotgbot.Bot, ctx *ext.Context) error {
		widgetQueryHandled = true
		_, err := bot.AnswerWidgetQuery(ctx.WidgetQuery.Id, nil)
		return err
	}))

	// 2. Add SharedWidgetQuery handler
	dispatcher.AddHandler(handlers.NewSharedWidgetQuery(sharedwidgetquery.Equal("share_action"), func(bot *gotgbot.Bot, ctx *ext.Context) error {
		sharedWidgetQueryHandled = true
		return nil
	}))

	// 3. Add Message handler with message.Widget filter
	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Widget != nil
	}, func(bot *gotgbot.Bot, ctx *ext.Context) error {
		widgetMessageHandled = true
		return nil
	}))

	// --- Test Dispatching WidgetQuery ---
	wqUpdate := &gotgbot.Update{
		UpdateId: 1,
		WidgetQuery: &gotgbot.WidgetQuery{
			Id:   "query_999",
			From: gotgbot.User{Id: 42, FirstName: "Tester"},
			Data: "click_buy",
		},
	}
	err := dispatcher.ProcessUpdate(b, wqUpdate, nil)
	if err != nil {
		t.Fatalf("failed to process WidgetQuery: %v", err)
	}
	if !widgetQueryHandled {
		t.Errorf("WidgetQuery handler did not run")
	}
	if answeredQueryId != "query_999" {
		t.Errorf("expected answeredQueryId to be 'query_999', got %q", answeredQueryId)
	}

	// --- Test Dispatching SharedWidgetQuery ---
	swqUpdate := &gotgbot.Update{
		UpdateId: 2,
		SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
			Id:       "shared_888",
			From:     gotgbot.User{Id: 43, FirstName: "Sharer"},
			Data:     "share_action",
			WidgetId: "widget_abc",
		},
	}
	err = dispatcher.ProcessUpdate(b, swqUpdate, nil)
	if err != nil {
		t.Fatalf("failed to process SharedWidgetQuery: %v", err)
	}
	if !sharedWidgetQueryHandled {
		t.Errorf("SharedWidgetQuery handler did not run")
	}

	// --- Test Dispatching Message with Widget ---
	widgetMsgUpdate := &gotgbot.Update{
		UpdateId: 3,
		Message: &gotgbot.Message{
			MessageId: 101,
			From:      &gotgbot.User{Id: 42},
			Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
			Widget:    &gotgbot.Widget{},
		},
	}
	err = dispatcher.ProcessUpdate(b, widgetMsgUpdate, nil)
	if err != nil {
		t.Fatalf("failed to process Widget Message: %v", err)
	}
	if !widgetMessageHandled {
		t.Errorf("Widget message handler did not run")
	}

	// --- Test Dispatching Normal Message (should NOT trigger widget filter) ---
	normalMsgUpdate := &gotgbot.Update{
		UpdateId: 4,
		Message: &gotgbot.Message{
			MessageId: 102,
			From:      &gotgbot.User{Id: 42},
			Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
			Text:      "just text",
		},
	}
	widgetMessageHandled = false
	err = dispatcher.ProcessUpdate(b, normalMsgUpdate, nil)
	if err != nil {
		t.Fatalf("failed to process normal Message: %v", err)
	}
	if widgetMessageHandled {
		t.Errorf("Widget message handler ran on normal message")
	}
}

type mockBotClient struct {
	requestFunc func(ctx context.Context, token string, method string, params map[string]any, opts *gotgbot.RequestOpts) (json.RawMessage, error)
}

func (m *mockBotClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]any, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	if m.requestFunc != nil {
		return m.requestFunc(ctx, token, method, params, opts)
	}
	return json.Marshal(true)
}

func (m *mockBotClient) GetAPIURL(opts *gotgbot.RequestOpts) string {
	return "https://api.telegram.org"
}

func (m *mockBotClient) FileURL(token string, tgFilePath string, opts *gotgbot.RequestOpts) string {
	return "https://api.telegram.org/file"
}
