package handlers_test

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/widgetquery"
)

func TestWidgetQueryHandler(t *testing.T) {
	b := NewTestBot()

	var handled bool
	h := handlers.NewWidgetQuery(widgetquery.Prefix("action_"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		handled = true
		return nil
	})

	// Matching update
	matchingUpd := ext.NewContext(b, &gotgbot.Update{
		WidgetQuery: &gotgbot.WidgetQuery{
			Id:   "wq_1",
			From: gotgbot.User{Id: 123},
			Data: "action_submit",
		},
	}, nil)

	if !h.CheckUpdate(b, matchingUpd) {
		t.Errorf("expected handler to match widget query with prefix action_")
	}

	err := h.HandleUpdate(b, matchingUpd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handled {
		t.Errorf("expected response function to run")
	}

	// Nil context check
	if h.CheckUpdate(b, nil) {
		t.Errorf("expected handler not to match nil context")
	}

	// Empty handler response
	emptyHandler := handlers.WidgetQuery{}
	if err := emptyHandler.HandleUpdate(b, matchingUpd); err != nil {
		t.Errorf("expected no error from empty handler: %v", err)
	}

	// Non-matching update
	nonMatchingUpd := ext.NewContext(b, &gotgbot.Update{
		WidgetQuery: &gotgbot.WidgetQuery{
			Id:   "wq_2",
			From: gotgbot.User{Id: 123},
			Data: "other_data",
		},
	}, nil)

	if h.CheckUpdate(b, nonMatchingUpd) {
		t.Errorf("expected handler not to match widget query with different prefix")
	}

	// Nil widget query update
	msgUpd := ext.NewContext(b, &gotgbot.Update{
		Message: &gotgbot.Message{Text: "hello"},
	}, nil)
	if h.CheckUpdate(b, msgUpd) {
		t.Errorf("expected handler not to match non-widget-query update")
	}
}

func TestWidgetQueryFilters(t *testing.T) {
	wq := &gotgbot.WidgetQuery{
		Id:           "123",
		From:         gotgbot.User{Id: 456},
		Data:         "test_action_run",
		ChatInstance: "inst_1",
		ChatType:     "private",
		WidgetId:     "wid_1",
	}

	tests := []struct {
		name   string
		filter filters.WidgetQuery
		input  *gotgbot.WidgetQuery
		want   bool
	}{
		{name: "All filter", filter: widgetquery.All, input: wq, want: true},
		{name: "Equal filter match", filter: widgetquery.Equal("test_action_run"), input: wq, want: true},
		{name: "Equal filter mismatch", filter: widgetquery.Equal("other"), input: wq, want: false},
		{name: "Prefix filter match", filter: widgetquery.Prefix("test_"), input: wq, want: true},
		{name: "Prefix filter mismatch", filter: widgetquery.Prefix("other_"), input: wq, want: false},
		{name: "Suffix filter match", filter: widgetquery.Suffix("_run"), input: wq, want: true},
		{name: "Suffix filter mismatch", filter: widgetquery.Suffix("_stop"), input: wq, want: false},
		{name: "Contains filter match", filter: widgetquery.Contains("action"), input: wq, want: true},
		{name: "Contains filter mismatch", filter: widgetquery.Contains("other"), input: wq, want: false},
		{name: "FromUserID filter match", filter: widgetquery.FromUserID(456), input: wq, want: true},
		{name: "FromUserID filter mismatch", filter: widgetquery.FromUserID(999), input: wq, want: false},
		{name: "ChatInstance filter match", filter: widgetquery.ChatInstance("inst_1"), input: wq, want: true},
		{name: "ChatInstance filter mismatch", filter: widgetquery.ChatInstance("inst_2"), input: wq, want: false},
		{name: "ChatType filter match", filter: widgetquery.ChatType("private"), input: wq, want: true},
		{name: "ChatType filter mismatch", filter: widgetquery.ChatType("channel"), input: wq, want: false},
		{name: "WidgetID filter match", filter: widgetquery.WidgetID("wid_1"), input: wq, want: true},
		{name: "WidgetID filter mismatch", filter: widgetquery.WidgetID("wid_2"), input: wq, want: false},
		// Nil input checks
		{name: "Equal nil", filter: widgetquery.Equal("test"), input: nil, want: false},
		{name: "Prefix nil", filter: widgetquery.Prefix("test"), input: nil, want: false},
		{name: "Suffix nil", filter: widgetquery.Suffix("test"), input: nil, want: false},
		{name: "Contains nil", filter: widgetquery.Contains("test"), input: nil, want: false},
		{name: "FromUserID nil", filter: widgetquery.FromUserID(123), input: nil, want: false},
		{name: "ChatInstance nil", filter: widgetquery.ChatInstance("inst"), input: nil, want: false},
		{name: "ChatType nil", filter: widgetquery.ChatType("private"), input: nil, want: false},
		{name: "WidgetID nil", filter: widgetquery.WidgetID("wid"), input: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter(tt.input); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
