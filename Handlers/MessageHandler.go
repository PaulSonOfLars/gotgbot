package Handlers

import (
	"gotgbot/Ext"
	"gotgbot/Types"
	"gotgbot"
)


type Message struct {
	filterFunc func(message *Types.Message) bool
	response   func(b Ext.Bot, u gotgbot.Update)

}

func NewMessage(filterFunc func(message *Types.Message) bool,
				response func(b Ext.Bot, u gotgbot.Update)) Message {
	h := Message{}
	h.filterFunc = filterFunc
	h.response = response
	return h
}

func (h Message) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	go h.response(d.Bot, update)

}

func (h Message) CheckUpdate(update gotgbot.Update) bool {
	return update.Message != nil && h.filterFunc(update.Message)
}