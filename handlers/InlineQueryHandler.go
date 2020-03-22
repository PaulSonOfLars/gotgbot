package handlers

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type InlineQuery struct {
	baseHandler
	Match    string
	Response func(b ext.Bot, u *gotgbot.Update) error
}

func NewInlineQuery(match string, response func(b ext.Bot, u *gotgbot.Update) error) InlineQuery {
	return InlineQuery{
		baseHandler: baseHandler{
			Name: "unnamedInlineQueryHandler",
		},
		Match:    match,
		Response: response,
	}
}

func (h InlineQuery) HandleUpdate(u *gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(*d.Bot, u)
}

func (h InlineQuery) CheckUpdate(u *gotgbot.Update) (bool, error) {
	if u.InlineQuery == nil {
		return false, nil
	}

	res, err := regexp.MatchString(h.Match, u.InlineQuery.Query)
	if err != nil {
		return false, errors.Wrapf(err, "Could not match regexp")
	}
	return res, nil
}
