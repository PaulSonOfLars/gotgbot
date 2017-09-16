package library


type Handler interface {
	HandleUpdate(update Update, d Dispatcher)
	CheckUpdate(update Update) bool
}


type Update struct {
	Update_id            int
	Message              *Message
	Edited_message       *Message
	Channel_post         *Message
	Edited_channel_post  *Message
	Inline_query         *InlineQuery
	Chosen_inline_result *ChosenInlineResult
	Callback_query       *CallbackQuery
	Shipping_query       *ShippingQuery
	Pre_checkout_query   *PreCheckoutQuery

}

