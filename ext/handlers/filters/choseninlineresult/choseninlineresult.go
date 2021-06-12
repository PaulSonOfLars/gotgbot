package choseninlineresult

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.ChosenInlineResult) bool {
	return true
}

func FromUserID(id int64) filters.ChosenInlineResult {
	return func(cir *gotgbot.ChosenInlineResult) bool {
		return cir.From.Id == id
	}
}

func Query(q string) filters.ChosenInlineResult {
	return func(cir *gotgbot.ChosenInlineResult) bool {
		return cir.Query == q
	}
}

func QueryPrefix(prefix string) filters.ChosenInlineResult {
	return func(cir *gotgbot.ChosenInlineResult) bool {
		return strings.HasPrefix(cir.Query, prefix)
	}
}

func QuerySuffix(suffix string) filters.ChosenInlineResult {
	return func(cir *gotgbot.ChosenInlineResult) bool {
		return strings.HasSuffix(cir.Query, suffix)
	}
}

func InlineMessageId(msgId string) filters.ChosenInlineResult {
	return func(cir *gotgbot.ChosenInlineResult) bool {
		return cir.InlineMessageId == msgId
	}
}

func Location(cir *gotgbot.ChosenInlineResult) bool {
	return cir.Location != nil
}
