package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot"
)

type Message struct {
	AllowEdited bool
	filterFunc  func(message *ext.Message) bool
	response    func(b ext.Bot, u gotgbot.Update) error
}

type FilterAble func(message *ext.Message) bool

func NewMessage(filterFunc FilterAble, response func(b ext.Bot, u gotgbot.Update) error) Message {
	return Message{
		AllowEdited: false,
		filterFunc:  filterFunc,
		response:    response,
	}
}

func (h Message) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) error {
	return h.response(d.Bot, update)
}

func (h Message) CheckUpdate(update gotgbot.Update) (bool, error) {
	return (update.Message != nil && h.filterFunc(update.Message)) ||
		(update.EditedMessage != nil && h.filterFunc(update.EditedMessage)), nil
}
