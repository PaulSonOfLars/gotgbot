package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type ManagedBot struct {
	Filter   filters.ManagedBot
	Response Response
}

func NewManagedBot(filter filters.ManagedBot, r Response) ManagedBot {
	return ManagedBot{
		Filter:   filter,
		Response: r,
	}
}

func (cb ManagedBot) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return cb.Response(b, ctx)
}

func (cb ManagedBot) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.ManagedBot == nil {
		return false
	}

	return cb.Filter == nil || cb.Filter(ctx.ManagedBot)
}

func (cb ManagedBot) Name() string {
	return fmt.Sprintf("managed_bot_%p", cb.Response)
}
