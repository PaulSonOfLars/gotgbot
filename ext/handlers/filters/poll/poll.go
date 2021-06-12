package poll

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.Poll) bool {
	return true
}

func Id(id string) filters.Poll {
	return func(p *gotgbot.Poll) bool {
		return p.Id == id
	}
}

func Type(t string) filters.Poll {
	return func(p *gotgbot.Poll) bool {
		return p.Type == t
	}
}

func Regular(p *gotgbot.Poll) bool {
	return p.Type == "regular"
}

func Quiz(p *gotgbot.Poll) bool {
	return p.Type == "quiz"
}
