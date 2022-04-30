package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type Poll struct {
	Filter   filters.Poll
	Response Response
}

func NewPoll(f filters.Poll, r Response) Poll {
	return Poll{
		Filter:   f,
		Response: r,
	}
}

func (r Poll) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.Poll == nil {
		return false
	}
	return r.Filter == nil || r.Filter(u.Poll)
}

func (r Poll) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r Poll) Name() string {
	return fmt.Sprintf("poll_%p", r.Response)
}
