package ext

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: extend to be used as a generic cancel context
type Context struct {
	*gotgbot.Update
	Data map[string]interface{}

	EffectiveMessage *gotgbot.Message
	EffectiveChat    *gotgbot.Chat
	EffectiveUser    *gotgbot.User
}

func NewContext(b *gotgbot.Bot, update *gotgbot.Update) *Context {
	var msg *gotgbot.Message
	var chat *gotgbot.Chat
	var user *gotgbot.User

	if update.Message != nil {
		msg = update.Message
		chat = &update.Message.Chat
		user = update.Message.From

	} else if update.EditedMessage != nil {
		msg = update.EditedMessage
		chat = &update.EditedMessage.Chat
		user = update.EditedMessage.From

	} else if update.ChannelPost != nil {
		msg = update.ChannelPost
		chat = &update.ChannelPost.Chat

	} else if update.EditedChannelPost != nil {
		msg = update.EditedChannelPost
		chat = &update.EditedChannelPost.Chat

	} else if update.InlineQuery != nil {
		user = &update.InlineQuery.From
	} else if update.CallbackQuery != nil {
		user = &update.CallbackQuery.From

		if update.CallbackQuery.Message != nil {
			msg = update.CallbackQuery.Message
			chat = &update.CallbackQuery.Message.Chat
		}

	} else if update.ChosenInlineResult != nil {
		user = &update.ChosenInlineResult.From

	} else if update.ShippingQuery != nil {
		user = &update.ShippingQuery.From

	} else if update.PreCheckoutQuery != nil {
		user = &update.PreCheckoutQuery.From

	} else if update.MyChatMember != nil {
		user = &update.MyChatMember.From
		chat = &update.MyChatMember.Chat

	} else if update.ChatMember != nil {
		user = &update.ChatMember.From
		chat = &update.ChatMember.Chat

	}

	return &Context{
		Update:           update,
		Data:             make(map[string]interface{}),
		EffectiveMessage: msg,
		EffectiveChat:    chat,
		EffectiveUser:    user,
	}
}

// Args gets the list of whitespace-separated arguments of the message text.
func (c *Context) Args() []string {
	var msg *gotgbot.Message

	if c.Update.Message != nil {
		msg = c.Update.Message
	} else if c.Update.EditedMessage != nil {
		msg = c.Update.EditedMessage
	} else if c.Update.ChannelPost != nil {
		msg = c.Update.ChannelPost
	} else if c.Update.EditedChannelPost != nil {
		msg = c.Update.EditedChannelPost
	} else if c.Update.CallbackQuery != nil && c.Update.CallbackQuery.Message != nil {
		msg = c.Update.CallbackQuery.Message
	}

	if msg == nil {
		return nil
	}

	if msg.Text != "" {
		return strings.Fields(c.Update.Message.Text)
	} else if msg.Caption != "" {
		return strings.Fields(c.Update.Message.Caption)
	}

	return nil
}
