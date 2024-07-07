package businessconnection

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.BusinessConnection) bool {
	return true
}

func BusinessId(businessId string) filters.BusinessConnection {
	return func(bc *gotgbot.BusinessConnection) bool {
		return bc.Id == businessId
	}
}

func FromUserID(userId int64) filters.BusinessConnection {
	return func(bc *gotgbot.BusinessConnection) bool {
		return bc.User.Id == userId
	}
}

func ChatID(bc *gotgbot.BusinessConnection, chatId int64) filters.BusinessConnection {
	return func(bc *gotgbot.BusinessConnection) bool {
		return bc.UserChatId == chatId
	}
}

func CanReply(bc *gotgbot.BusinessConnection) bool {
	return bc.CanReply
}
