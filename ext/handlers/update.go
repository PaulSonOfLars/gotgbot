package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Update struct {
	Response Response
}

func NewUpdate(r Response) Update {
	return Update{
		Response: r,
	}
}

func (u Update) CheckUpdate(b *gotgbot.Bot, upd *gotgbot.Update) bool {
	return upd != nil
}

func (u Update) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return u.Response(b, ctx)
}

func (u Update) Name() string {
	return fmt.Sprintf("update_%p", u.Response)
}
