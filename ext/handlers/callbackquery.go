package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type CallbackQuery struct {
	AllowChannel bool
	Filter       filters.CallbackQuery
	Response     Response
}

func NewCallback(filter filters.CallbackQuery, r Response) CallbackQuery {
	return CallbackQuery{
		Filter:   filter,
		Response: r,
	}
}

func (cb CallbackQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return cb.Response(b, ctx)
}

func (cb CallbackQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.CallbackQuery == nil {
		return false
	}

	if ctx.CallbackQuery.Message != nil {
		m, ok := ctx.CallbackQuery.Message.(*gotgbot.Message)
		if !cb.AllowChannel && ok && m.Chat.Type == "channel" {
			return false
		}
	}

	return cb.Filter == nil || cb.Filter(ctx.CallbackQuery)
}

func (cb CallbackQuery) Name() string {
	return fmt.Sprintf("inlinequery_%p", cb.Response)
}
