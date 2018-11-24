package handlers

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type FilterFunc func(message *ext.Message) bool

type Message struct {
	AllowEdited bool
	Filter      FilterFunc
	Response    func(b ext.Bot, u gotgbot.Update) error
}

func NewMessage(filterFunc FilterFunc, response func(b ext.Bot, u gotgbot.Update) error) Message {
	return Message{
		AllowEdited: false,
		Filter:      filterFunc,
		Response:    response,
	}
}

func (h Message) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.Response(d.Bot, update)
}

func (h Message) CheckUpdate(update gotgbot.Update) (bool, error) {
	return (update.Message != nil && h.Filter(update.Message)) ||
		(h.AllowEdited && update.EditedMessage != nil && h.Filter(update.EditedMessage)), nil
}
