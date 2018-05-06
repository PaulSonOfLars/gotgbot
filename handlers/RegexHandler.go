package handlers

import (
	"regexp"
	"log"
	"gotgbot/ext"
	"gotgbot"
)

type Regex struct {
	match string
	response func(b ext.Bot, u gotgbot.Update)

}

func NewRegex(match string, response func(b ext.Bot, u gotgbot.Update)) Regex {
	h := Regex{}
	h.match = match
	h.response = response
	return h
}

func (h Regex) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	go h.response(d.Bot, update)

}

func (h Regex) CheckUpdate(update gotgbot.Update) bool {
	if update.Message != nil {
		res, err := regexp.Match(h.match, []byte(update.Message.Text))
		if err != nil {
			log.Fatal(err)
		}
		return res
	} else {
		return false
	}
}
