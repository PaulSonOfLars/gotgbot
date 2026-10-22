package message

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func TestFromUserID_NilFrom(t *testing.T) {
	// Should not panic when From is nil (e.g. channel post)
	msg := &gotgbot.Message{
		From: nil,
	}
	f := FromUserID(1234567890)
	if f(msg) {
		t.Errorf("expected false for nil From")
	}
}

func TestFromUsername_NilFrom(t *testing.T) {
	// Should not panic when From is nil (e.g. channel post)
	msg := &gotgbot.Message{
		From: nil,
	}
	f := FromUsername("example")
	if f(msg) {
		t.Errorf("expected false for nil From")
	}
}
