package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Filter func(msg *gotgbot.Message) bool

type Message struct {
	AllowEdited  bool
	AllowChannel bool
	Filter       Filter
	Response     func(ctx *ext.Context) error
}

func NewMessage(f Filter, r func(ctx *ext.Context) error) Message {
	return Message{
		AllowEdited:  false,
		AllowChannel: false,
		Filter:       f,
		Response:     r,
	}
}

func (m Message) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.Message != nil {
		if u.Message.Text == "" && u.Message.Caption == "" {
			return false
		}
		return m.Filter(u.Message)
	}

	// if no edits and message is edited
	if m.AllowEdited && u.EditedMessage != nil {
		if u.EditedMessage.Text == "" && u.EditedMessage.Caption == "" {
			return false
		}
		return m.Filter(u.EditedMessage)
	}
	// if no channel and message is channel message
	if m.AllowChannel && u.ChannelPost != nil {
		if u.ChannelPost.Text == "" && u.ChannelPost.Caption == "" {
			return false
		}
		return m.Filter(u.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if m.AllowChannel && m.AllowEdited && u.EditedChannelPost != nil {
		if u.EditedChannelPost.Text == "" && u.EditedChannelPost.Caption == "" {
			return false
		}
		return m.Filter(u.EditedChannelPost)
	}

	return false
}

func (m Message) HandleUpdate(ctx *ext.Context) error {
	return m.Response(ctx)
}

func (m Message) Name() string {
	return fmt.Sprintf("message_%T", m.Response)
}
