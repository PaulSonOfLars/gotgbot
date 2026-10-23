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
