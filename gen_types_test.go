package gotgbot

import (
	"testing"
)

func TestInlineQueriesHavePointers(t *testing.T) {
	// This test is somewhat ridiculous, in that it actually only tests compilation.
	// But it broke once, and I won't have it break again.

	// Future readers, this is required because inline keyboard buttons can be passed many different things.
	// One of those things, is an empty string "switch inline query" field. Go's type system sees this as an empty
	// value, so doesn't include it in the JSON marshalling, since it has an omitempty tag. We therefore use a pointer
	// here to differentiate between empty field (nil), and empty value ("").
	// Reported as a bug here: https://t.me/GotgbotChat/4537
	// Fixed here: https://github.com/PaulSonOfLars/gotgbot/pull/31, and again here https://github.com/PaulSonOfLars/gotgbot/pull/63

	stringValue := "Foo"
	_ = InlineKeyboardButton{
		Text:                         "Barr",
		SwitchInlineQuery:            nil,          // ilq can be nil
		SwitchInlineQueryCurrentChat: &stringValue, // or can be a pointer
	}
}
