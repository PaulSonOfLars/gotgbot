package reaction

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.MessageReactionUpdated) bool {
	return true
}

func FromUserID(id int64) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		if mru.User != nil {
			return mru.User.Id == id
		}

		return false
	}
}

func FromAnonymousChatID(id int64) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		if mru.ActorChat != nil {
			return mru.ActorChat.Id == id
		}

		return false
	}
}

func FromChatID(id int64) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		return mru.Chat.Id == id
	}
}

func NewReactionIn(reaction string) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		for _, r := range mru.NewReaction {
			if r.MergeReactionType().Emoji == reaction {
				return true
			}
		}

		return false
	}
}

func OldReactionIn(reaction string) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		for _, r := range mru.OldReaction {
			if r.MergeReactionType().Emoji == reaction {
				return true
			}
		}

		return false
	}
}
