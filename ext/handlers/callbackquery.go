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

func NewCallback(filter filters.CallbackQuery, response Response) CallbackQuery {
	return CallbackQuery{
		Filter:   filter,
		Response: response,
	}
}

func (cb CallbackQuery) HandleUpdate(ctx *ext.Context) error {
	return cb.Response(ctx)
}

func (cb CallbackQuery) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.CallbackQuery == nil {
		return false
	}

	if !cb.AllowChannel && u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat.Type == "channel" {
		return false
	}

	return cb.Filter == nil || cb.Filter(u.CallbackQuery)
}

func (cb CallbackQuery) Name() string {
	return fmt.Sprintf("inlinequery_%p", cb.Response)
}
