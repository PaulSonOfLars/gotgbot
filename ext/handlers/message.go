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

func (m Message) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.Message != nil {
		return m.Filter == nil || m.Filter(ctx.Message)
	}

	// if no edits and message is edited
	if m.AllowEdited && ctx.EditedMessage != nil {
		if ctx.EditedMessage.Text == "" && ctx.EditedMessage.Caption == "" {
			return false
		}
		return m.Filter == nil || m.Filter(ctx.EditedMessage)
	}
	// if no channel and message is channel message
	if m.AllowChannel && ctx.ChannelPost != nil {
		if ctx.ChannelPost.Text == "" && ctx.ChannelPost.Caption == "" {
			return false
		}
		return m.Filter == nil || m.Filter(ctx.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if m.AllowChannel && m.AllowEdited && ctx.EditedChannelPost != nil {
		if ctx.EditedChannelPost.Text == "" && ctx.EditedChannelPost.Caption == "" {
			return false
		}
		return m.Filter == nil || m.Filter(ctx.EditedChannelPost)
	}

	return false
}

func (m Message) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return m.Response(b, ctx)
}

func (m Message) Name() string {
	return fmt.Sprintf("message_%p", m.Response)
}
