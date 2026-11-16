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

	dispatcher := ext.NewDispatcher(nil)

	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("click_"), func(bot *gotgbot.Bot, ctx *ext.Context) error {
		widgetQueryHandled = true
		_, err := bot.AnswerWidgetQuery(ctx.WidgetQuery.Id, nil)
		return err
	}))

	dispatcher.AddHandler(handlers.NewSharedWidgetQuery(sharedwidgetquery.Equal("share_action"), func(bot *gotgbot.Bot, ctx *ext.Context) error {
		sharedWidgetQueryHandled = true
		return nil
	}))

	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Widget != nil
	}, func(bot *gotgbot.Bot, ctx *ext.Context) error {
		widgetMessageHandled = true
		return nil
	}))

	tests := []struct {
		name       string
		update     *gotgbot.Update
		checkState func(t *testing.T)
	}{
		{
			name: "WidgetQuery",
			update: &gotgbot.Update{
				UpdateId: 1,
				WidgetQuery: &gotgbot.WidgetQuery{
					Id:   "query_999",
					From: gotgbot.User{Id: 42, FirstName: "Tester"},
					Data: "click_buy",
				},
			},
			checkState: func(t *testing.T) {
				if !widgetQueryHandled {
					t.Errorf("expected WidgetQuery handler to run")
				}
				if answeredQueryId != "query_999" {
					t.Errorf("expected answeredQueryId 'query_999', got %q", answeredQueryId)
				}
			},
		},
		{
			name: "SharedWidgetQuery",
			update: &gotgbot.Update{
				UpdateId: 2,
				SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
					Id:       "shared_888",
					From:     gotgbot.User{Id: 43, FirstName: "Sharer"},
					Data:     "share_action",
					WidgetId: "widget_abc",
				},
			},
			checkState: func(t *testing.T) {
				if !sharedWidgetQueryHandled {
					t.Errorf("expected SharedWidgetQuery handler to run")
				}
			},
		},
		{
			name: "MessageWithWidget",
			update: &gotgbot.Update{
				UpdateId: 3,
				Message: &gotgbot.Message{
					MessageId: 101,
					From:      &gotgbot.User{Id: 42},
					Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
					Widget:    &gotgbot.Widget{},
				},
			},
			checkState: func(t *testing.T) {
				if !widgetMessageHandled {
					t.Errorf("expected widget message handler to run")
				}
			},
		},
		{
			name: "NormalMessageWithoutWidget",
			update: &gotgbot.Update{
				UpdateId: 4,
				Message: &gotgbot.Message{
					MessageId: 102,
					From:      &gotgbot.User{Id: 42},
					Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
					Text:      "plain text without widget",
				},
			},
			checkState: func(t *testing.T) {
				if widgetMessageHandled {
					t.Errorf("widget message handler should not run for normal message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widgetMessageHandled = false
			err := dispatcher.ProcessUpdate(b, tt.update, nil)
			if err != nil {
				t.Fatalf("unexpected error processing update: %v", err)
			}
			tt.checkState(t)
		})
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
