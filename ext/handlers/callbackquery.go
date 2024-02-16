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

// SetAllowChannel Enables channel messages for this handler.
func (cb CallbackQuery) SetAllowChannel(allow bool) CallbackQuery {
	cb.AllowChannel = allow
	return cb
}

func (cb CallbackQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return cb.Response(b, ctx)
}

func (cb CallbackQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.CallbackQuery == nil {
		return false
	}

	if !cb.AllowChannel && ctx.CallbackQuery.Message != nil && ctx.CallbackQuery.Message.GetChat().Type == "channel" {
		return false
	}

	return cb.Filter == nil || cb.Filter(ctx.CallbackQuery)
}

func (cb CallbackQuery) Name() string {
	return fmt.Sprintf("callback_query_handler_%p", cb.Response)
}
