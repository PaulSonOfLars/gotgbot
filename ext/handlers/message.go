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
	if u.Message != nil {
		return m.Filter == nil || m.Filter(u.Message)
	}

	// if no edits and message is edited
	if m.AllowEdited && u.EditedMessage != nil {
		if u.EditedMessage.Text == "" && u.EditedMessage.Caption == "" {
			return false
		}
		return m.Filter == nil || m.Filter(u.EditedMessage)
	}
	// if no channel and message is channel message
	if m.AllowChannel && u.ChannelPost != nil {
		if u.ChannelPost.Text == "" && u.ChannelPost.Caption == "" {
			return false
		}
		return m.Filter == nil || m.Filter(u.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if m.AllowChannel && m.AllowEdited && u.EditedChannelPost != nil {
		if u.EditedChannelPost.Text == "" && u.EditedChannelPost.Caption == "" {
			return false
		}
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
