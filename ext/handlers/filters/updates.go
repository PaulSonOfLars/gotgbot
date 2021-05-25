package filters

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type ChatMember func(u *gotgbot.ChatMemberUpdated) bool

func ChatMemberUserId(id int64) ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.NewChatMember.User.Id == id
	}
}

func ChatMemberFromUserId(id int64) ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.From.Id == id
	}
}

func ChatMemberChatId(id int64) ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.Chat.Id == id
	}
}

func ChatMemberAll() ChatMember {
	return func(_ *gotgbot.ChatMemberUpdated) bool {
		return true
	}
}

func ChatMemberGroup() ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.Chat.Type == "group"
	}
}

func ChatMemberSupergroup() ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.Chat.Type == "supergroup"
	}
}

func ChatMemberChannel() ChatMember {
	return func(cm *gotgbot.ChatMemberUpdated) bool {
		return cm.Chat.Type == "channel"
	}
}
