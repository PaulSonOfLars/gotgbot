package callbackquery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.CallbackQuery) bool {
	return true
}

func Prefix(prefix string) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasPrefix(cq.Data, prefix)
	}
}

func Suffix(suffix string) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasSuffix(cq.Data, suffix)
	}
}

func Equal(match string) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.Data == match
	}
}

func FromUserID(id int64) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.From.Id == id
	}
}

func GameName(name string) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.GameShortName == name
	}
}

func Inline(cq *gotgbot.CallbackQuery) bool {
	return cq.InlineMessageId != ""
}

func ChatInstance(instance string) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.ChatInstance == instance
	}
}
