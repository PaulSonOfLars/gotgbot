package handlers

import (
	"regexp"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/pkg/errors"
)

type CallBack struct {
	baseHandler
	Pattern  string
	Response func(b ext.Bot, u gotgbot.Update) error
}

func NewCallback(pattern string, response func(b ext.Bot, u gotgbot.Update) error) CallBack {
	return CallBack{
		baseHandler: baseHandler{
			Name: pattern,
		},
		Pattern:     pattern,
		Response:    response,
	}
}

func (cb CallBack) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) error {
	return cb.Response(d.Bot, update)
}

func (cb CallBack) CheckUpdate(update gotgbot.Update) (bool, error) {
	if update.CallbackQuery == nil {
		return false, nil
	}
	if cb.Pattern != "" {
		res, err := regexp.MatchString(cb.Pattern, update.CallbackQuery.Data)
		if err != nil {
			return false, errors.Wrapf(err, "Could not match regexp")
		}
		return res, nil
	}
	return true, nil
}
