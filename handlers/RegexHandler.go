package handlers

import (
	"regexp"
	"gotgbot/ext"
	"gotgbot"
	"github.com/pkg/errors"
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
	h.response(d.Bot, update)

}

func (h Regex) CheckUpdate(update gotgbot.Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}
	res, err := regexp.Match(h.match, []byte(update.Message.Text))
	if err != nil {
		return false, errors.Wrapf(err, "Could not match regexp")
	}
	return res, nil
}
