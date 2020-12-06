package ext

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: extend to be used as a generic cancel context
type Context struct {
	Bot    *gotgbot.Bot
	Update *gotgbot.Update
	Data   map[string]string

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
	}

	return &Context{
		Bot:              b,
		Update:           update,
		Data:             make(map[string]string),
		EffectiveMessage: msg,
		EffectiveChat:    chat,
		EffectiveUser:    user,
	}
}

func (c *Context) Args() []string {
	if c.Update.Message != nil {
		return strings.Fields(c.Update.Message.Text)
	}

	if c.Update.EditedMessage != nil {
		return strings.Fields(c.Update.EditedMessage.Text)
	}

	if c.Update.ChannelPost != nil {
		return strings.Fields(c.Update.ChannelPost.Text)
	}

	if c.Update.EditedChannelPost != nil {
		return strings.Fields(c.Update.EditedChannelPost.Text)
	}

	if c.Update.CallbackQuery.Message != nil {
		return strings.Fields(c.Update.CallbackQuery.Message.Text)
	}

	return nil
}
