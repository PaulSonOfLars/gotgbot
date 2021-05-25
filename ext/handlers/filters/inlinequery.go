package filters

import "github.com/PaulSonOfLars/gotgbot/v2"

type InlineQuery func(iq *gotgbot.InlineQuery) bool

func InlineUserID(id int64) InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return iq.From.Id == id
	}
}

func Query(q string) InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return iq.Query == q
	}
}
