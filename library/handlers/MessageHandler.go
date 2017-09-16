package handlers

import (
	"bot/library"
)


type Message struct {
	filterFunc func(message *library.Message) bool
	response   func(b library.Bot, u library.Update)

}

func NewMessage(filterFunc func(message *library.Message) bool, response func(b library.Bot, u library.Update)) Message {
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