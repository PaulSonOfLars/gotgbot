package gotgbot

import (
	"gotgbot/ext"
	"gotgbot/types"
)

type Handler interface {
	HandleUpdate(update Update, d Dispatcher)
	CheckUpdate(update Update) bool
}

type Update struct {
	Update_id            int
	Message              *types.Message
	Edited_message       *types.Message
	Channel_post         *types.Message
	Edited_channel_post  *types.Message
	Inline_query         *types.Message
	Chosen_inline_result *types.ChosenInlineResult
	Callback_query       *types.CallbackQuery
	Shipping_query       *types.ShippingQuery
	Pre_checkout_query   *types.PreCheckoutQuery

	// Self added type
	Effective_message *ext.Message
}
