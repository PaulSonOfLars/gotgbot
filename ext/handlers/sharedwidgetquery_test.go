package handlers_test

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
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

	if !sharedwidgetquery.All(swq) {
		t.Errorf("All filter failed")
	}
	if !sharedwidgetquery.Equal("view_details_done")(swq) {
		t.Errorf("Equal filter failed")
	}
	if !sharedwidgetquery.Prefix("view_")(swq) {
		t.Errorf("Prefix filter failed")
	}
	if !sharedwidgetquery.Suffix("_done")(swq) {
		t.Errorf("Suffix filter failed")
	}
	if !sharedwidgetquery.Contains("details")(swq) {
		t.Errorf("Contains filter failed")
	}
	if !sharedwidgetquery.FromUserID(999)(swq) {
		t.Errorf("FromUserID filter failed")
	}
	if !sharedwidgetquery.WidgetID("wid_555")(swq) {
		t.Errorf("WidgetID filter failed")
	}

	// Nil safety checks
	if sharedwidgetquery.Equal("test")(nil) {
		t.Errorf("Equal filter should return false for nil")
	}
	if sharedwidgetquery.Prefix("test")(nil) {
		t.Errorf("Prefix filter should return false for nil")
	}
	if sharedwidgetquery.Suffix("test")(nil) {
		t.Errorf("Suffix filter should return false for nil")
	}
	if sharedwidgetquery.Contains("test")(nil) {
		t.Errorf("Contains filter should return false for nil")
	}
	if sharedwidgetquery.FromUserID(123)(nil) {
		t.Errorf("FromUserID filter should return false for nil")
	}
	if sharedwidgetquery.WidgetID("wid")(nil) {
		t.Errorf("WidgetID filter should return false for nil")
	}
}
