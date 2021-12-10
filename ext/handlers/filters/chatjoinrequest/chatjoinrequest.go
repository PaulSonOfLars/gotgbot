package chatjoinrequest

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.ChatJoinRequest) bool {
	return true
}

func ChatID(id int64) filters.ChatJoinRequest {
	return func(r *gotgbot.ChatJoinRequest) bool {
		return r.Chat.Id == id
	}
}
