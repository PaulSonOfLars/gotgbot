package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type PollAnswer struct {
	Filter   filters.PollAnswer
	Response Response
}

func NewPollAnswer(f filters.PollAnswer, r Response) PollAnswer {
	return PollAnswer{
		Filter:   f,
		Response: r,
	}
}

func (r PollAnswer) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	if u.PollAnswer == nil {
		return false
	}
	return r.Filter == nil || r.Filter(u.PollAnswer)
}

func (r PollAnswer) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r PollAnswer) Name() string {
	return fmt.Sprintf("pollanswer_%p", r.Response)
}
