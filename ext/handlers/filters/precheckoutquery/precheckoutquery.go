package shippingquery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.PreCheckoutQuery) bool {
	return true
}

func FromUserID(id int64) filters.PreCheckoutQuery {
	return func(p *gotgbot.PreCheckoutQuery) bool {
		return p.From.Id == id
	}
}

func HasPayloadPrefix(pre string) filters.PreCheckoutQuery {
	return func(p *gotgbot.PreCheckoutQuery) bool {
		return strings.HasPrefix(p.InvoicePayload, pre)
	}
}
