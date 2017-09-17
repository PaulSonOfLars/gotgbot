package Handlers

import (
	"bot/library"
	"bot/library/Ext"
	"bot/library/Types"
)


type Message struct {
	filterFunc func(message *Types.Message) bool
	response   func(b Ext.Bot, u library.Update)

}

func NewMessage(filterFunc func(message *Types.Message) bool,
				response func(b Ext.Bot, u library.Update)) Message {
	h := Message{}
	h.filterFunc = filterFunc
	h.response = response
	return h
}

func (h Message) HandleUpdate(update library.Update, d library.Dispatcher) {
	h.response(d.Bot, update)

}

func (h Message) CheckUpdate(update library.Update) bool {
	return h.filterFunc(update.Message)
}