package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type ShippingQuery struct {
	Filter   filters.ShippingQuery
	Response Response
}

func NewShippingQuery(f filters.ShippingQuery, r Response) ShippingQuery {
	return ShippingQuery{
		Filter:   f,
		Response: r,
	}
}

func (r ShippingQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.ShippingQuery == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.ShippingQuery)
}

func (r ShippingQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r ShippingQuery) Name() string {
	return fmt.Sprintf("shippingquery_%p", r.Response)
}
