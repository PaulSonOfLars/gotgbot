package handlers

import (
	"gotgbot/ext"
	"gotgbot/types"
	"gotgbot"
)


type Message struct {
	filterFunc func(message *types.Message) bool
	response   func(b ext.Bot, u gotgbot.Update)

}

func NewMessage(filterFunc func(message *types.Message) bool,
				response func(b ext.Bot, u gotgbot.Update)) Message {
	h := Message{}
	h.filterFunc = filterFunc
	h.response = response
	return h
}

func (h Message) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	h.response(d.Bot, update)

}

func (h Message) CheckUpdate(update gotgbot.Update) (bool, error) {
	return update.Message != nil && h.filterFunc(update.Message), nil
}