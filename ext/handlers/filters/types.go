package filters

import "github.com/PaulSonOfLars/gotgbot/v2"

type (
	CallbackQuery      func(cq *gotgbot.CallbackQuery) bool
	ChatMember         func(u *gotgbot.ChatMemberUpdated) bool
	ChosenInlineResult func(cir *gotgbot.ChosenInlineResult) bool
	InlineQuery        func(iq *gotgbot.InlineQuery) bool
	Message            func(msg *gotgbot.Message) bool
	Poll               func(poll *gotgbot.Poll) bool
)
