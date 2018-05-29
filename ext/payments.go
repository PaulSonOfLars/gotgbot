package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
)

// TODO: all the optionals here. Best option is probably to use a builder.
func (b Bot) SendInvoice(chatId int, title string, description string, payload string,
	providerToken string, startParameter string, currency string,
	prices []types.LabeledPrice) (*Message, error) {
	return b.NewSendableInvoice(chatId, title, description, payload, providerToken, startParameter, currency, prices).Send()
}

func (b Bot) AnswerShippingQuery(shippingQueryId string, ok bool) (bool, error) {
	return b.NewSendableAnswerShippingQuery(shippingQueryId, ok).Send()
}

func (b Bot) AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool) (bool, error) {
	return b.NewSendableAnswerPreCheckoutQuery(preCheckoutQueryId, ok).Send()
}
