package message

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func TestFromFilters(t *testing.T) {
	tests := []struct {
		name     string
		filter   filters.Message
		msg      *gotgbot.Message
		expected bool
	}{
		{
			name:     "FromUserID with nil From",
			filter:   FromUserID(123),
			msg:      &gotgbot.Message{From: nil},
			expected: false,
		},
		{
			name:     "FromUserID matching ID",
			filter:   FromUserID(123),
			msg:      &gotgbot.Message{From: &gotgbot.User{Id: 123}},
			expected: true,
		},
		{
			name:     "FromUserID mismatching ID",
			filter:   FromUserID(123),
			msg:      &gotgbot.Message{From: &gotgbot.User{Id: 456}},
			expected: false,
		},
		{
			name:     "FromUsername with nil From",
			filter:   FromUsername("alice"),
			msg:      &gotgbot.Message{From: nil},
			expected: false,
		},
		{
			name:     "FromUsername matching username",
			filter:   FromUsername("alice"),
			msg:      &gotgbot.Message{From: &gotgbot.User{Username: "alice"}},
			expected: true,
		},
		{
			name:     "FromUsername mismatching username",
			filter:   FromUsername("alice"),
			msg:      &gotgbot.Message{From: &gotgbot.User{Username: "bob"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter(tt.msg)
			if got != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
