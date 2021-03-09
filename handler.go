package gotgbot

import (
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/ext"
)

type Handler interface {
	// HandleUpdate processes the update. The error return can be used to return either EndGroups{} or ContinueGroups{}
	// errors, which allows for either continuing execution of the current group, or stopping all further execution.
	HandleUpdate(u *Update, d Dispatcher) error
	// CheckUpdate checks whether the update should be processes or not.
	CheckUpdate(u *Update) (bool, error)
	// GetName gets the handler name; used to differentiate handlers programmatically.
	GetName() string
}

// Update Incoming updates from telegram servers (including extra, internally used fields added for convenience)
type Update struct {
	UpdateId           int                     `json:"update_id"`
	Message            *ext.Message            `json:"message"`
	EditedMessage      *ext.Message            `json:"edited_message"`
	ChannelPost        *ext.Message            `json:"channel_post"`
	EditedChannelPost  *ext.Message            `json:"edited_channel_post"`
	InlineQuery        *ext.InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ext.ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *ext.CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ext.ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *ext.PreCheckoutQuery   `json:"pre_checkout_query"`
	Poll               *ext.Poll               `json:"poll"`
	PollAnswer         *ext.PollAnswer         `json:"poll_answer"`
	MyChatMember       *ext.ChatMemberUpdated  `json:"my_chat_member"`
	ChatMember         *ext.ChatMemberUpdated  `json:"chat_member"`

	// Self added type
	EffectiveMessage *ext.Message `json:"effective_message"`
	EffectiveChat    *ext.Chat    `json:"effective_chat"`
	EffectiveUser    *ext.User    `json:"effective_user"`
	Data             map[string]string
}

func initUpdate(data RawUpdate, bot ext.Bot) (*Update, error) {
	var upd Update
	if err := json.Unmarshal(data, &upd); err != nil {
		return nil, err
	}
	if upd.Message != nil {
		upd.EffectiveMessage = upd.Message
		upd.EffectiveChat = upd.Message.Chat
		upd.EffectiveUser = upd.Message.From

	} else if upd.EditedMessage != nil {
		upd.EffectiveMessage = upd.EditedMessage
		upd.EffectiveChat = upd.EditedMessage.Chat
		upd.EffectiveUser = upd.EditedMessage.From

	} else if upd.ChannelPost != nil {
		upd.EffectiveMessage = upd.ChannelPost
		upd.EffectiveChat = upd.ChannelPost.Chat

	} else if upd.EditedChannelPost != nil {
		upd.EffectiveMessage = upd.EditedChannelPost
		upd.EffectiveChat = upd.EditedChannelPost.Chat

	} else if upd.InlineQuery != nil {
		upd.EffectiveUser = upd.InlineQuery.From

	} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
		upd.CallbackQuery.Bot = bot
		upd.EffectiveMessage = upd.CallbackQuery.Message
		upd.EffectiveChat = upd.CallbackQuery.Message.Chat
		upd.EffectiveUser = upd.CallbackQuery.From

	} else if upd.ChosenInlineResult != nil {
		upd.EffectiveUser = upd.ChosenInlineResult.From

	} else if upd.ShippingQuery != nil {
		upd.EffectiveUser = upd.ShippingQuery.From

	} else if upd.PreCheckoutQuery != nil {
		upd.EffectiveUser = upd.PreCheckoutQuery.From

	} else if upd.MyChatMember != nil {
		upd.EffectiveUser = upd.MyChatMember.NewChatMember.User
		upd.EffectiveChat = &upd.MyChatMember.Chat

		upd.MyChatMember.From.Bot = bot
		upd.MyChatMember.OldChatMember.User.Bot = bot

	} else if upd.ChatMember != nil {
		upd.EffectiveUser = upd.ChatMember.NewChatMember.User
		upd.EffectiveChat = &upd.ChatMember.Chat

		upd.ChatMember.From.Bot = bot
		upd.ChatMember.OldChatMember.User.Bot = bot
	}

	if upd.EffectiveMessage != nil {
		upd.EffectiveMessage.Bot = bot
		if upd.EffectiveMessage.ReplyToMessage != nil {
			upd.EffectiveMessage.ReplyToMessage.Bot = bot
			if upd.EffectiveMessage.ReplyToMessage.From != nil {
				upd.EffectiveMessage.ReplyToMessage.From.Bot = bot
			}
		}
	}
	if upd.EffectiveChat != nil {
		upd.EffectiveChat.Bot = bot
	}
	if upd.EffectiveUser != nil {
		upd.EffectiveUser.Bot = bot
	}
	upd.Data = make(map[string]string)
	return &upd, nil
}
