package handlers

import (
	"regexp"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/pkg/errors"
)

type Regex struct {
	baseHandler
	Match    string
	Response func(b ext.Bot, u *gotgbot.Update) error
}

func NewRegex(match string, response func(b ext.Bot, u *gotgbot.Update) error) Regex {
	return Regex{
		baseHandler: baseHandler{
			Name: match,
		},
		Match:       match,
		Response:    response,
	}
}

func (h Regex) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(d.Bot, u)
}

func (h Regex) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.EffectiveMessage == nil {
		return false, nil
	}
	text := u.EffectiveMessage.Text
	if text == "" {
		text = u.EffectiveMessage.Caption
	}
	res, err := regexp.Match(h.Match, []byte(text))
	if err != nil {
		return false, errors.Wrapf(err, "Could not match regexp")
	}
	return res, nil
}
