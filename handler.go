package gotgbot

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type Handler interface {
	HandleUpdate(u *Update, d Dispatcher) error
	CheckUpdate(u *Update) (bool, error)
	GetName() string
}

type Update struct {
	UpdateId           int                     `json:"update_id"`
	Message            *ext.Message            `json:"message"`
	EditedMessage      *ext.Message            `json:"edited_message"`
	ChannelPost        *ext.Message            `json:"channel_post"`
	EditedChannelPost  *ext.Message            `json:"edited_channel_post"`
	InlineQuery        *ext.Message            `json:"inline_query"`
	ChosenInlineResult *ext.ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *ext.CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ext.ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *ext.PreCheckoutQuery   `json:"pre_checkout_query"`

	// Self added type
	EffectiveMessage *ext.Message `json:"effective_message"`
	EffectiveChat    *ext.Chat    `json:"effective_chat"`
	EffectiveUser    *ext.User    `json:"effective_user"`
	data             map[string]string
}

func (u *Update) GetData() map[string]string {
	if u.data == nil {
		u.data = make(map[string]string)
	}
	return u.data
}
