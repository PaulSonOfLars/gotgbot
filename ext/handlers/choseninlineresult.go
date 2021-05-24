package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type ChosenInlineResult struct {
	Filter   filters.ChosenInlineResult
	Response Response
}

func NewChosenInlineResult(filter filters.ChosenInlineResult, r Response) ChosenInlineResult {
	return ChosenInlineResult{
		Filter:   filter,
		Response: r,
	}
}

func (i ChosenInlineResult) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return i.Response(b, ctx)
}

func (i ChosenInlineResult) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.ChosenInlineResult == nil {
		return false
	}

	return i.Filter == nil || i.Filter(u.ChosenInlineResult)
}

func (i ChosenInlineResult) Name() string {
	return fmt.Sprintf("choseninlineresult_%p", i.Response)
}
