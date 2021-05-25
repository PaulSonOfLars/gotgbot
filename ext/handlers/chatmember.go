package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type ChatMember struct {
	Response Response
	Filter   filters.ChatMember
}

func NewChatMember(f filters.ChatMember, r Response) ChatMember {
	return ChatMember{
		Response: r,
		Filter:   f,
	}
}

func (c ChatMember) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.ChatMember == nil {
		return false
	}
	
	return c.Filter == nil || c.Filter(u.ChatMember)
}

func (c ChatMember) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.Response(b, ctx)
}

func (c ChatMember) Name() string {
	return fmt.Sprintf("chatmember_%p", c.Response)
}
