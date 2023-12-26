package shippingquery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.ShippingQuery) bool {
	return true
}

func FromUserID(id int64) filters.ShippingQuery {
	return func(p *gotgbot.ShippingQuery) bool {
		return p.From.Id == id
	}
}

func HasPayloadPrefix(pre string) filters.ShippingQuery {
	return func(p *gotgbot.ShippingQuery) bool {
		return strings.HasPrefix(p.InvoicePayload, pre)
	}
}
