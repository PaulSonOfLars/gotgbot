package filters

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type CallbackQuery func(cq *gotgbot.CallbackQuery) bool

func Prefix(prefix string) CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasPrefix(cq.Data, prefix)
	}
}

func Suffix(suffix string) CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasSuffix(cq.Data, suffix)
	}
}

func Equal(match string) CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.Data == match
	}
}

func CallbackUserID(id int64) CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.From.Id == id
	}
}

func GameName(name string) CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return cq.GameShortName == name
	}
}
