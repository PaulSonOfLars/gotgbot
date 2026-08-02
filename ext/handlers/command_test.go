package handlers

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func TestCheckMessage(t *testing.T) {
	bot := &gotgbot.Bot{User: gotgbot.User{Username: "MyBot"}}
	c := Command{Triggers: []rune{'!'}, Command: "start"}

	tests := []struct {
		name        string
		text        string
		richMessage *gotgbot.RichMessage
		want        bool
	}{
		{name: "matching command", text: "!start", want: true},
		{name: "matching command with target", text: "!start@MyBot", want: true},
		{name: "wrong bot target", text: "!start@OtherBot", want: false},
		{name: "wrong command", text: "!stop", want: false},
		{name: "no trigger", text: "start", want: false},
		{name: "empty text", text: "", want: false},
		{name: "whitespace only", text: "   ", want: false},
		{name: "command case insensitive", text: "!START", want: true}, // Command should be case-insensitive
		{name: "bot username case insensitive", text: "!start@mybot", want: true},
		{name: "rich message / command", richMessage: &gotgbot.RichMessage{
			Blocks: []gotgbot.RichBlock{
				gotgbot.RichBlockParagraph{Text: gotgbot.RichTextBotCommand{
					Text:       gotgbot.RichTextString("/start"),
					BotCommand: "start",
				}},
			},
		}, want: false}, // Only have ! enabled as a trigger
		{name: "rich message ! command", richMessage: &gotgbot.RichMessage{
			Blocks: []gotgbot.RichBlock{
				gotgbot.RichBlockParagraph{Text: gotgbot.RichTextString("!start")},
			},
		}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &gotgbot.Message{Text: tt.text, RichMessage: tt.richMessage}
			got := c.checkMessage(bot, msg)
			if got != tt.want {
				t.Errorf("checkMessage(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestExtractCommand(t *testing.T) {
	c := Command{Triggers: []rune{'!', '/', '?'}}

	tests := []struct {
		name    string
		command string
		want    string
	}{
		// Basic cases
		{"ascii trigger match", "!hello", "hello"},
		{"second trigger match", "/hello", "hello"},
		{"third trigger match", "?hello", "hello"},
		{"no match", "hello", ""},
		{"unrecognised trigger", "$hello", ""},

		// Edge cases
		{"empty string", "", ""},
		{"trigger only", "!", ""},
		{"trigger with spaces", "! hello", " hello"},

		// Multi-byte rune triggers
		{"multibyte no match", "你好", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.extractCommand(tt.command)
			if got != tt.want {
				t.Errorf("extractCommand(%q) = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}

func TestExtractCommandMultibyteTrigger(t *testing.T) {
	c := Command{Triggers: []rune{'你'}}
	got := c.extractCommand("你好")
	if got != "好" {
		t.Errorf("expected 好, got %q", got)
	}
}
