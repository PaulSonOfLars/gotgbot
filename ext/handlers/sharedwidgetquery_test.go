package handlers_test

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/sharedwidgetquery"
)

func TestSharedWidgetQueryHandler(t *testing.T) {
	b := NewTestBot()

	var handled bool
	h := handlers.NewSharedWidgetQuery(sharedwidgetquery.Prefix("share_"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		handled = true
		return nil
	})

	matchingUpd := ext.NewContext(b, &gotgbot.Update{
		SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
			Id:       "swq_1",
			From:     gotgbot.User{Id: 789},
			Data:     "share_confirm",
			WidgetId: "w_100",
		},
	}, nil)

	if !h.CheckUpdate(b, matchingUpd) {
		t.Errorf("expected handler to match shared widget query with prefix share_")
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
	emptyHandler := handlers.SharedWidgetQuery{}
	if err := emptyHandler.HandleUpdate(b, matchingUpd); err != nil {
		t.Errorf("expected no error from empty handler: %v", err)
	}

	nonMatchingUpd := ext.NewContext(b, &gotgbot.Update{
		SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
			Id:       "swq_2",
			From:     gotgbot.User{Id: 789},
			Data:     "other_data",
			WidgetId: "w_100",
		},
	}, nil)

	if h.CheckUpdate(b, nonMatchingUpd) {
		t.Errorf("expected handler not to match shared widget query with different prefix")
	}
}

func TestSharedWidgetQueryFilters(t *testing.T) {
	swq := &gotgbot.SharedWidgetQuery{
		Id:       "swq_1",
		From:     gotgbot.User{Id: 999},
		Data:     "view_details_done",
		WidgetId: "wid_555",
	}

	tests := []struct {
		name   string
		filter filters.SharedWidgetQuery
		input  *gotgbot.SharedWidgetQuery
		want   bool
	}{
		{name: "All filter", filter: sharedwidgetquery.All, input: swq, want: true},
		{name: "Equal filter match", filter: sharedwidgetquery.Equal("view_details_done"), input: swq, want: true},
		{name: "Equal filter mismatch", filter: sharedwidgetquery.Equal("other"), input: swq, want: false},
		{name: "Prefix filter match", filter: sharedwidgetquery.Prefix("view_"), input: swq, want: true},
		{name: "Prefix filter mismatch", filter: sharedwidgetquery.Prefix("edit_"), input: swq, want: false},
		{name: "Suffix filter match", filter: sharedwidgetquery.Suffix("_done"), input: swq, want: true},
		{name: "Suffix filter mismatch", filter: sharedwidgetquery.Suffix("_start"), input: swq, want: false},
		{name: "Contains filter match", filter: sharedwidgetquery.Contains("details"), input: swq, want: true},
		{name: "Contains filter mismatch", filter: sharedwidgetquery.Contains("settings"), input: swq, want: false},
		{name: "FromUserID filter match", filter: sharedwidgetquery.FromUserID(999), input: swq, want: true},
		{name: "FromUserID filter mismatch", filter: sharedwidgetquery.FromUserID(111), input: swq, want: false},
		{name: "WidgetID filter match", filter: sharedwidgetquery.WidgetID("wid_555"), input: swq, want: true},
		{name: "WidgetID filter mismatch", filter: sharedwidgetquery.WidgetID("other_wid"), input: swq, want: false},
		// Nil input checks
		{name: "Equal nil", filter: sharedwidgetquery.Equal("test"), input: nil, want: false},
		{name: "Prefix nil", filter: sharedwidgetquery.Prefix("test"), input: nil, want: false},
		{name: "Suffix nil", filter: sharedwidgetquery.Suffix("test"), input: nil, want: false},
		{name: "Contains nil", filter: sharedwidgetquery.Contains("test"), input: nil, want: false},
		{name: "FromUserID nil", filter: sharedwidgetquery.FromUserID(123), input: nil, want: false},
		{name: "WidgetID nil", filter: sharedwidgetquery.WidgetID("wid"), input: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter(tt.input); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
