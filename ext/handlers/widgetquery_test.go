package handlers_test

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
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

	if !widgetquery.All(wq) {
		t.Errorf("All filter should return true")
	}
	if !widgetquery.Equal("test_action_run")(wq) {
		t.Errorf("Equal filter failed")
	}
	if !widgetquery.Prefix("test_")(wq) {
		t.Errorf("Prefix filter failed")
	}
	if !widgetquery.Suffix("_run")(wq) {
		t.Errorf("Suffix filter failed")
	}
	if !widgetquery.Contains("action")(wq) {
		t.Errorf("Contains filter failed")
	}
	if !widgetquery.FromUserID(456)(wq) {
		t.Errorf("FromUserID filter failed")
	}
	if !widgetquery.ChatInstance("inst_1")(wq) {
		t.Errorf("ChatInstance filter failed")
	}
	if !widgetquery.ChatType("private")(wq) {
		t.Errorf("ChatType filter failed")
	}
	if !widgetquery.WidgetID("wid_1")(wq) {
		t.Errorf("WidgetID filter failed")
	}

	// Nil safety checks
	if widgetquery.Equal("test")(nil) {
		t.Errorf("Equal filter should return false for nil")
	}
	if widgetquery.Prefix("test")(nil) {
		t.Errorf("Prefix filter should return false for nil")
	}
	if widgetquery.Suffix("test")(nil) {
		t.Errorf("Suffix filter should return false for nil")
	}
	if widgetquery.Contains("test")(nil) {
		t.Errorf("Contains filter should return false for nil")
	}
	if widgetquery.FromUserID(123)(nil) {
		t.Errorf("FromUserID filter should return false for nil")
	}
	if widgetquery.ChatInstance("inst")(nil) {
		t.Errorf("ChatInstance filter should return false for nil")
	}
	if widgetquery.ChatType("private")(nil) {
		t.Errorf("ChatType filter should return false for nil")
	}
	if widgetquery.WidgetID("wid")(nil) {
		t.Errorf("WidgetID filter should return false for nil")
	}
}
