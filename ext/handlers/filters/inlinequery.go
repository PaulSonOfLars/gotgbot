package filters

import "github.com/PaulSonOfLars/gotgbot/v2"

type InlineQuery func(iq *gotgbot.InlineQuery) bool

func (f InlineQuery) And(f2 InlineQuery) InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return f(iq) && f2(iq)
	}
}

func (f InlineQuery) Or(f2 InlineQuery) InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return f(iq) || f2(iq)
	}
}

func (f InlineQuery) Not() InlineQuery {
	return func(iq *gotgbot.InlineQuery) bool {
		return !f(iq)
	}
}

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
