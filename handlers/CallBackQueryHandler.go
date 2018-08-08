package handlers

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/pkg/errors"
	"regexp"
)

type CallBack struct {
	pattern  string
	response func(b ext.Bot, u gotgbot.Update)
}

func NewCallback(pattern string, response func(b ext.Bot, u gotgbot.Update)) CallBack {
	return CallBack{
		pattern:  pattern,
		response: response,
	}
}

func (cb CallBack) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	cb.response(d.Bot, update)
}

func (cb CallBack) CheckUpdate(update gotgbot.Update) (bool, error) {
	if update.CallbackQuery == nil {
		return false, nil
	}
	if cb.pattern != "" {
		res, err := regexp.MatchString(cb.pattern, update.CallbackQuery.Data)
		if err != nil {
			return false, errors.Wrapf(err, "Could not match regexp")
		}
		return res, nil
	}
	return true, nil
}
