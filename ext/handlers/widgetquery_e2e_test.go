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
	tests := []struct {
		name                     string
		update                   *gotgbot.Update
		wantWidgetQueryHandled   bool
		wantSharedQueryHandled   bool
		wantWidgetMessageHandled bool
		wantAnsweredQueryId      string
	}{
		{
			name: "WidgetQuery matching",
			update: &gotgbot.Update{
				UpdateId: 1,
				WidgetQuery: &gotgbot.WidgetQuery{
					Id:   "query_999",
					From: gotgbot.User{Id: 42, FirstName: "Tester"},
					Data: "click_buy",
				},
			},
			wantWidgetQueryHandled: true,
			wantAnsweredQueryId:    "query_999",
		},
		{
			name: "WidgetQuery non-matching prefix",
			update: &gotgbot.Update{
				UpdateId: 2,
				WidgetQuery: &gotgbot.WidgetQuery{
					Id:   "query_111",
					From: gotgbot.User{Id: 42, FirstName: "Tester"},
					Data: "other_prefix",
				},
			},
			wantWidgetQueryHandled: false,
		},
		{
			name: "SharedWidgetQuery matching",
			update: &gotgbot.Update{
				UpdateId: 3,
				SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
					Id:       "shared_888",
					From:     gotgbot.User{Id: 43, FirstName: "Sharer"},
					Data:     "share_action",
					WidgetId: "widget_abc",
				},
			},
			wantSharedQueryHandled: true,
		},
		{
			name: "SharedWidgetQuery non-matching",
			update: &gotgbot.Update{
				UpdateId: 4,
				SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
					Id:       "shared_222",
					From:     gotgbot.User{Id: 43, FirstName: "Sharer"},
					Data:     "other_action",
					WidgetId: "widget_abc",
				},
			},
			wantSharedQueryHandled: false,
		},
		{
			name: "Message with widget",
			update: &gotgbot.Update{
				UpdateId: 5,
				Message: &gotgbot.Message{
					MessageId: 101,
					From:      &gotgbot.User{Id: 42},
					Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
					Widget:    &gotgbot.Widget{},
				},
			},
			wantWidgetMessageHandled: true,
		},
		{
			name: "Normal message without widget",
			update: &gotgbot.Update{
				UpdateId: 6,
				Message: &gotgbot.Message{
					MessageId: 102,
					From:      &gotgbot.User{Id: 42},
					Chat:      gotgbot.Chat{Id: -1001, Type: "supergroup"},
					Text:      "plain text without widget",
				},
			},
			wantWidgetMessageHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			err := dispatcher.ProcessUpdate(b, tt.update, nil)
			if err != nil {
				t.Fatalf("unexpected error processing update: %v", err)
			}

			if widgetQueryHandled != tt.wantWidgetQueryHandled {
				t.Errorf("widgetQueryHandled = %v, want %v", widgetQueryHandled, tt.wantWidgetQueryHandled)
			}
			if answeredQueryId != tt.wantAnsweredQueryId {
				t.Errorf("answeredQueryId = %q, want %q", answeredQueryId, tt.wantAnsweredQueryId)
			}
			if sharedWidgetQueryHandled != tt.wantSharedQueryHandled {
				t.Errorf("sharedWidgetQueryHandled = %v, want %v", sharedWidgetQueryHandled, tt.wantSharedQueryHandled)
			}
			if widgetMessageHandled != tt.wantWidgetMessageHandled {
				t.Errorf("widgetMessageHandled = %v, want %v", widgetMessageHandled, tt.wantWidgetMessageHandled)
			}
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
