package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type Reaction struct {
	Filter   filters.Reaction
	Response Response
}

func NewReaction(f filters.Reaction, r Response) Reaction {
	return Reaction{
		Filter:   f,
		Response: r,
	}
}

func (r Reaction) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.MessageReaction == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.MessageReaction)
}

func (r Reaction) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r Reaction) Name() string {
	return fmt.Sprintf("reaction_%p", r.Response)
}
