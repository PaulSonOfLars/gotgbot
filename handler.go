package gotgbot

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/types"
)

type Handler interface {
	HandleUpdate(update Update, d Dispatcher)
	CheckUpdate(update Update) (bool, error)
}

type Update struct {
	UpdateId           int                       `json:"update_id"`
	Message            *types.Message            `json:"message"`
	EditedMessage      *types.Message            `json:"edited_message"`
	ChannelPost        *types.Message            `json:"channel_post"`
	EditedChannelPost  *types.Message            `json:"edited_channel_post"`
	InlineQuery        *types.Message            `json:"inline_query"`
	ChosenInlineResult *types.ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *types.CallbackQuery      `json:"callback_query"`
	ShippingQuery      *types.ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *types.PreCheckoutQuery   `json:"pre_checkout_query"`

	// Self added type
	EffectiveMessage *ext.Message `json:"effective_message"`
	EffectiveChat    *ext.Chat    `json:"effective_chat"`
	EffectiveUser    *ext.User    `json:"effective_user"`
}
