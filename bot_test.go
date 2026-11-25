package gotgbot

import (
	"testing"
)

func TestNewBot_DisableTokenCheck(t *testing.T) {
	b, err := NewBot("123:abc", &BotOpts{
		DisableTokenCheck: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.Id != 123 {
		t.Errorf("expected bot ID 123, got %d", b.Id)
	}
	if !b.IsBot {
		t.Errorf("expected IsBot to be true")
	}
	if b.FirstName != "" {
		t.Errorf("expected empty FirstName, got %q", b.FirstName)
	}
	if b.Username != "" {
		t.Errorf("expected empty Username, got %q", b.Username)
	}
	if b.MatchesUsername("<missing>") {
		t.Errorf("MatchesUsername('<missing>') should return false")
	}
}
