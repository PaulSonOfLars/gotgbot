package managedbot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.ManagedBotUpdated) bool {
	return true
}

func OwnedByUserID(id int64) filters.ManagedBot {
	return func(mbu *gotgbot.ManagedBotUpdated) bool {
		return mbu.User.Id == id
	}
}

func AboutBotID(id int64) filters.ManagedBot {
	return func(mbu *gotgbot.ManagedBotUpdated) bool {
		return mbu.Bot.Id == id
	}
}
