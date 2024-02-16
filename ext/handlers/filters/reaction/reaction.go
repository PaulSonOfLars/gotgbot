package reaction

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.MessageReactionUpdated) bool {
	return true
}

func FromPeer(id int64) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		if mru.User != nil {
			return mru.User.Id == id
		}

		if mru.ActorChat != nil {
			return mru.ActorChat.Id == id
		}

		return false
	}
}

func ChatID(id int64) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		return mru.Chat.Id == id
	}
}

func NewReactionEmoji(reaction string) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		for _, r := range mru.NewReaction {
			if r.MergeReactionType().Emoji == reaction {
				return true
			}
		}

		return false
	}
}

func OldReactionEmoji(reaction string) filters.Reaction {
	return func(mru *gotgbot.MessageReactionUpdated) bool {
		for _, r := range mru.OldReaction {
			if r.MergeReactionType().Emoji == reaction {
				return true
			}
		}

		return false
	}
}
