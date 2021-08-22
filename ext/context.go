package ext

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: extend to be used as a generic cancel context?
type Context struct {
	*gotgbot.Update
	Data map[string]interface{}

	// EffectiveMessage is the message which triggered the update, if possible
	EffectiveMessage *gotgbot.Message
	// EffectiveChat is the chat the update was triggered in, if possible
	EffectiveChat *gotgbot.Chat
	// EffectiveUser is the user who triggered the update, if possible.
	// Note: when adding a user, the user who ADDED should be the EffectiveUser;
	// they caused the update. If a user joins naturally, then they are the EffectiveUser.
	EffectiveUser *gotgbot.User
}

// NewContext populates a context with the relevant fields from the current update.
// It takes a data field in the case where custom data needs to be passed.
func NewContext(update *gotgbot.Update, data map[string]interface{}) *Context {
	var msg *gotgbot.Message
	var chat *gotgbot.Chat
	var user *gotgbot.User

	switch {
	case update.Message != nil:
		msg = update.Message
		chat = &update.Message.Chat
		user = update.Message.From

	case update.EditedMessage != nil:
		msg = update.EditedMessage
		chat = &update.EditedMessage.Chat
		user = update.EditedMessage.From

	case update.ChannelPost != nil:
		msg = update.ChannelPost
		chat = &update.ChannelPost.Chat

	case update.EditedChannelPost != nil:
		msg = update.EditedChannelPost
		chat = &update.EditedChannelPost.Chat

	case update.InlineQuery != nil:
		user = &update.InlineQuery.From

	case update.CallbackQuery != nil:
		user = &update.CallbackQuery.From

		if update.CallbackQuery.Message != nil {
			msg = update.CallbackQuery.Message
			chat = &update.CallbackQuery.Message.Chat
		}

	case update.ChosenInlineResult != nil:
		user = &update.ChosenInlineResult.From

	case update.ShippingQuery != nil:
		user = &update.ShippingQuery.From

	case update.PreCheckoutQuery != nil:
		user = &update.PreCheckoutQuery.From

	case update.MyChatMember != nil:
		user = &update.MyChatMember.From
		chat = &update.MyChatMember.Chat

	case update.ChatMember != nil:
		user = &update.ChatMember.From
		chat = &update.ChatMember.Chat
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	return &Context{
		Update:           update,
		Data:             data,
		EffectiveMessage: msg,
		EffectiveChat:    chat,
		EffectiveUser:    user,
	}
}

// Args gets the list of whitespace-separated arguments of the message text.
func (c *Context) Args() []string {
	var msg *gotgbot.Message

	switch {
	case c.Update.Message != nil:
		msg = c.Update.Message

	case c.Update.EditedMessage != nil:
		msg = c.Update.EditedMessage

	case c.Update.ChannelPost != nil:
		msg = c.Update.ChannelPost

	case c.Update.EditedChannelPost != nil:
		msg = c.Update.EditedChannelPost

	case c.Update.CallbackQuery != nil && c.Update.CallbackQuery.Message != nil:
		msg = c.Update.CallbackQuery.Message
	}

	if msg == nil {
		return nil
	}

	if msg.Text != "" {
		return strings.Fields(msg.Text)
	} else if msg.Caption != "" {
		return strings.Fields(msg.Caption)
	}

	return nil
}
