package gotgbot

import (
	"encoding/json"
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

// TestUnmarshalMessage_UnknownRichBlockType verifies that a Message containing
// a RichBlock of a type this version of gotgbot doesn't know about unmarshals
// successfully with the unknown block preserved/skipped.
func TestUnmarshalMessage_UnknownRichBlockType(t *testing.T) {
	// Minimal Message JSON with a rich_message.blocks array containing
	// one legitimate block and one block of an unknown "type".
	raw := []byte(`{
		"message_id": 1,
		"date": 1234567890,
		"chat": {"id": 123, "type": "private"},
		"rich_message": {
			"blocks": [
				{"type": "paragraph", "text": {"text": "hello"}},
				{"type": "THIS IS NEW", "some field": [[{"some value": "test"}]]}
			]
		}
	}`)

	var msg Message
	err := json.Unmarshal(raw, &msg)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling message with unknown RichBlock type: %v", err)
	}
	if len(msg.RichMessage.Blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(msg.RichMessage.Blocks))
	}
	unknown, ok := msg.RichMessage.Blocks[1].(RichBlockUnknown)
	if !ok {
		t.Fatalf("expected second block to be RichBlockUnknown, got %T", msg.RichMessage.Blocks[1])
	}
	if unknown.GetType() != "THIS IS NEW" {
		t.Fatalf("expected unknown raw type 'THIS IS NEW', got %q", unknown.GetType())
	}
}
