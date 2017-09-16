package handlers

import (
	"bot/library"
	"regexp"
	"log"
)

type Regex struct {
	match string
	response func(b library.Bot, u library.Update)

}

func NewRegex(match string, response func(b library.Bot, u library.Update)) Regex {
	h := Regex{}
	h.match = match
	h.response = response
	return h
}

func (h Regex) HandleUpdate(update library.Update, d library.Dispatcher) {
	h.response(d.Bot, update)

}

func (h Regex) CheckUpdate(update library.Update) bool {
	res, err := regexp.Match(h.match, []byte(update.Message.Text))
	if err != nil {
		log.Fatal(err)
	}
	return res
}
