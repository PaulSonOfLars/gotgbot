package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot"

)


type Message struct {
	filterFunc func(message *ext.Message) bool
	response   func(b ext.Bot, u gotgbot.Update)

}

func NewMessage(filterFunc func(message *ext.Message) bool,
				response func(b ext.Bot, u gotgbot.Update)) Message {
	return Message{
		filterFunc: filterFunc,
		response:   response,
	}
}

func (h Message) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	h.response(d.Bot, update)

}

func (h Message) CheckUpdate(update gotgbot.Update) (bool, error) {
	return update.Message != nil && h.filterFunc(update.Message), nil
}
