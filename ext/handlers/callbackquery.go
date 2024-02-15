package handlers

import (
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type CallbackQuery struct {
	AllowChannel bool
	Filter       filters.CallbackQuery
	Response     func(b *gotgbot.Bot, ctx *ext.Context) error
}

func NewCallback(filter filters.CallbackQuery, response func(b *gotgbot.Bot, ctx *ext.Context) error) *CallbackQuery {
	return &CallbackQuery{
		Filter:   filter,
		Response: response,
	}
}

// SetAllowChannel enables channel messages for this handler.
func (cb *CallbackQuery) SetAllowChannel(allow bool) *CallbackQuery {
	cb.AllowChannel = allow
	return cb
}

func (cb *CallbackQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	err := cb.Response(b, ctx)
	if err != nil {
		log.Printf("Error handling callback query: %v", err)
	}
	return err
}

func (cb *CallbackQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.CallbackQuery == nil {
		return false
	}

	if !cb.AllowChannel && ctx.CallbackQuery.Message != nil && ctx.CallbackQuery.Message.GetChat().Type == "channel" {
		return false
	}

	return cb.Filter == nil || cb.Filter(ctx.CallbackQuery)
}

func (cb *CallbackQuery) Name() string {
	return fmt.Sprintf("callback_query_handler_%p", cb.Response)
}
