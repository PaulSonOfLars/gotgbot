package ext

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func TestNewContext_PurchasedPaidMedia(t *testing.T) {
	b := &gotgbot.Bot{
		User: gotgbot.User{Id: 1, IsBot: true, FirstName: "TestBot"},
	}
	upd := &gotgbot.Update{
		PurchasedPaidMedia: &gotgbot.PaidMediaPurchased{
			From:             gotgbot.User{Id: 12345, FirstName: "Buyer"},
			PaidMediaPayload: "payload",
		},
	}
	ctx := NewContext(b, upd, nil)
	if ctx.EffectiveUser == nil || ctx.EffectiveUser.Id != 12345 {
		t.Fatalf("expected EffectiveUser.Id = 12345, got %v", ctx.EffectiveUser)
	}
	if ctx.EffectiveSender == nil || ctx.EffectiveSender.Id() != 12345 {
		t.Fatalf("expected EffectiveSender.Id = 12345, got %v", ctx.EffectiveSender)
	}
}

func TestNewContext_WidgetQuery(t *testing.T) {
	b := &gotgbot.Bot{
		User: gotgbot.User{Id: 1, IsBot: true, FirstName: "TestBot"},
	}
	upd := &gotgbot.Update{
		WidgetQuery: &gotgbot.WidgetQuery{
			Id:   "wq_123",
			From: gotgbot.User{Id: 12345, FirstName: "Alice"},
			Chat: &gotgbot.Chat{Id: -100123, Type: "supergroup"},
			Data: "action_test",
		},
	}
	ctx := NewContext(b, upd, nil)
	if ctx.EffectiveUser == nil || ctx.EffectiveUser.Id != 12345 {
		t.Fatalf("expected EffectiveUser.Id = 12345, got %v", ctx.EffectiveUser)
	}
	if ctx.EffectiveChat == nil || ctx.EffectiveChat.Id != -100123 {
		t.Fatalf("expected EffectiveChat.Id = -100123, got %v", ctx.EffectiveChat)
	}
	if ctx.EffectiveSender == nil || ctx.EffectiveSender.Id() != 12345 {
		t.Fatalf("expected EffectiveSender.Id = 12345, got %v", ctx.EffectiveSender)
	}
}

func TestNewContext_SharedWidgetQuery(t *testing.T) {
	b := &gotgbot.Bot{
		User: gotgbot.User{Id: 1, IsBot: true, FirstName: "TestBot"},
	}
	upd := &gotgbot.Update{
		SharedWidgetQuery: &gotgbot.SharedWidgetQuery{
			Id:       "swq_456",
			From:     gotgbot.User{Id: 67890, FirstName: "Bob"},
			Data:     "shared_action",
			WidgetId: "w_1",
		},
	}
	ctx := NewContext(b, upd, nil)
	if ctx.EffectiveUser == nil || ctx.EffectiveUser.Id != 67890 {
		t.Fatalf("expected EffectiveUser.Id = 67890, got %v", ctx.EffectiveUser)
	}
}
