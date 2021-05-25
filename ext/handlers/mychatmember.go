package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type MyChatMember struct {
	Response Response
	Filter   filters.ChatMember
}

func NewMyChatMember(f filters.ChatMember, r Response) MyChatMember {
	return MyChatMember{
		Response: r,
		Filter:   f,
	}
}

func (m MyChatMember) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.MyChatMember == nil {
		return false
	}
	return m.Filter == nil || m.Filter(u.MyChatMember)
}

func (m MyChatMember) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return m.Response(b, ctx)
}

func (m MyChatMember) Name() string {
	return fmt.Sprintf("mychatmember_%p", m.Response)
}
