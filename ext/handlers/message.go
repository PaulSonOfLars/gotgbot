package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type Message struct {
	AllowEdited  bool
	AllowChannel bool
	Filter       filters.Message
	Response     Response
}

func NewMessage(f filters.Message, r Response) Message {
	return Message{
		AllowEdited:  false,
		AllowChannel: false,
		Filter:       f,
		Response:     r,
	}
}

func (m Message) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	// Normal incoming message in a group/private chat.
	if u.Message != nil {
		return m.Filter == nil || m.Filter(u.Message)
	}

	// If edits are allowed, and message is edited.
	if m.AllowEdited && u.EditedMessage != nil {
		return m.Filter == nil || m.Filter(u.EditedMessage)
	}

	// If channels are allowed, and message is a channel post.
	if m.AllowChannel && u.ChannelPost != nil {
		return m.Filter == nil || m.Filter(u.ChannelPost)
	}

	// If edits AND channels are allowed, and message is a channel post.
	if m.AllowChannel && m.AllowEdited && u.EditedChannelPost != nil {
		return m.Filter == nil || m.Filter(u.EditedChannelPost)
	}

	return false
}

func (m Message) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return m.Response(b, ctx)
}

func (m Message) Name() string {
	return fmt.Sprintf("message_%p", m.Response)
}
