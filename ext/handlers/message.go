package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type Message struct {
	AllowEdited   bool
	AllowChannel  bool
	AllowBusiness bool
	Filter        filters.Message
	Response      Response
}

func NewMessage(f filters.Message, r Response) Message {
	return Message{
		AllowEdited:  false,
		AllowChannel: false,
		Filter:       f,
		Response:     r,
	}
}

// SetAllowEdited Enables edited messages for this handler.
func (m Message) SetAllowEdited(allow bool) Message {
	m.AllowEdited = allow
	return m
}

// SetAllowChannel Enables channel messages for this handler.
func (m Message) SetAllowChannel(allow bool) Message {
	m.AllowChannel = allow
	return m
}

// SetAllowBusiness Enables business messages for this handler.
func (m Message) SetAllowBusiness(allow bool) Message {
	m.AllowBusiness = allow
	return m
}

func (m Message) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	// Handle NEW normal messages ONLY
	if ctx.Message != nil && ctx.EditedMessage == nil {
		return m.Filter == nil || m.Filter(ctx.Message)
	}
	// Handle EDITED normal messages ONLY if enabled
	if m.AllowEdited && ctx.EditedMessage != nil {
		return m.Filter == nil || m.Filter(ctx.EditedMessage)
	}
	// Handle NEW channel posts ONLY
	if m.AllowChannel && ctx.ChannelPost != nil && ctx.EditedChannelPost == nil {
		return m.Filter == nil || m.Filter(ctx.ChannelPost)
	}
	// Handle EDITED channel posts ONLY if both allowed
	if m.AllowChannel && m.AllowEdited && ctx.EditedChannelPost != nil {
		return m.Filter == nil || m.Filter(ctx.EditedChannelPost)
	}
	// Handle NEW business messages ONLY
	if m.AllowBusiness && ctx.BusinessMessage != nil && ctx.EditedBusinessMessage == nil {
		return m.Filter == nil || m.Filter(ctx.BusinessMessage)
	}
	// Handle EDITED business messages ONLY if both allowed
	if m.AllowBusiness && m.AllowEdited && ctx.EditedBusinessMessage != nil {
		return m.Filter == nil || m.Filter(ctx.EditedBusinessMessage)
	}
	return false
}

func (m Message) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return m.Response(b, ctx)
}

func (m Message) Name() string {
	return fmt.Sprintf("message_%p", m.Response)
}
