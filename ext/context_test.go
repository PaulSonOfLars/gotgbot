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

func TestNewContext_ChatBoostNilSource(t *testing.T) {
	b := &gotgbot.Bot{
		User: gotgbot.User{Id: 1, IsBot: true, FirstName: "TestBot"},
	}
	// Should not panic on nil Source
	upd := &gotgbot.Update{
		ChatBoost: &gotgbot.ChatBoostUpdated{
			Chat:  gotgbot.Chat{Id: -100},
			Boost: gotgbot.ChatBoost{},
		},
	}
	ctx := NewContext(b, upd, nil)
	if ctx.EffectiveChat == nil || ctx.EffectiveChat.Id != -100 {
		t.Fatalf("expected EffectiveChat.Id = -100, got %v", ctx.EffectiveChat)
	}
}
