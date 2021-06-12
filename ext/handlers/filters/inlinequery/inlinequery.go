package inlinequery

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.InlineQuery) bool {
	return true
}

func FromUserID(id int64) filters.InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return iq.From.Id == id
	}
}

func Query(q string) filters.InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return iq.Query == q
	}
}

func QueryPrefix(prefix string) filters.InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return strings.HasPrefix(iq.Query, prefix)
	}
}

func QuerySuffix(suffix string) filters.InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return strings.HasSuffix(iq.Query, suffix)
	}
}

func Sender(iq *gotgbot.InlineQuery) bool {
	return iq.ChatType == "sender"
}

func Private(iq *gotgbot.InlineQuery) bool {
	return iq.ChatType == "private"
}

func Group(iq *gotgbot.InlineQuery) bool {
	return iq.ChatType == "group"
}

func Supergroup(iq *gotgbot.InlineQuery) bool {
	return iq.ChatType == "supergroup"
}

func Channel(iq *gotgbot.InlineQuery) bool {
	return iq.ChatType == "channel"
}

func Location(iq *gotgbot.InlineQuery) bool {
	return iq.Location != nil
}
