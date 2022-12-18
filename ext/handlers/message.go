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

	// If edits are allowed, and message is edited.
	if m.AllowEdited && ctx.EditedMessage != nil {
		return m.Filter == nil || m.Filter(ctx.EditedMessage)
	}
	// If channel posts are allowed, and message is channel post.
	if m.AllowChannel && ctx.ChannelPost != nil {
		return m.Filter == nil || m.Filter(ctx.ChannelPost)
	}
	// If channel posts and edits are allowed, and post is edited.
	if m.AllowChannel && m.AllowEdited && ctx.EditedChannelPost != nil {
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
