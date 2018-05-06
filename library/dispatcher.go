package library

import (
	"bot/library/Ext"
)

type Dispatcher struct {
	Bot      Ext.Bot
	updates  chan Update
	handlers *[]Handler
}

func NewDispatcher(bot Ext.Bot, updates chan Update) Dispatcher {
	d := Dispatcher{}
	d.Bot = bot
	d.updates = updates
	d.handlers = new([]Handler)
	return d
}

func (d Dispatcher) Start() {
	for upd := range d.updates {
		d.process_update(upd)
	}
}

func (d Dispatcher) process_update(update Update) {
	for _, handler := range *d.handlers {
		if handler.CheckUpdate(update) {
			handler.HandleUpdate(update, d)
			break
		}
	}
}

func (d Dispatcher) Add_handler(handler Handler) {
	*d.handlers = append(*d.handlers, handler)

}
