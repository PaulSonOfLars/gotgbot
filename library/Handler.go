package library

import (
	"bot/library/Ext"
	"bot/library/Types"
)

type Handler interface {
	HandleUpdate(update Update, d Dispatcher)
	CheckUpdate(update Update) bool
}

type Update struct {
	Update_id            int
	Message              *Types.Message
	Edited_message       *Types.Message
	Channel_post         *Types.Message
	Edited_channel_post  *Types.Message
	Inline_query         *Types.Message
	Chosen_inline_result *Types.ChosenInlineResult
	Callback_query       *Types.CallbackQuery
	Shipping_query       *Types.ShippingQuery
	Pre_checkout_query   *Types.PreCheckoutQuery

	// Self added type
	Effective_message *Ext.Message
}
