package gotgbot

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
	"sort"
)

type Dispatcher struct {
	Bot           ext.Bot
	updates       chan Update
	handlers      map[int][]Handler
	handlerGroups *[]int
}

func NewDispatcher(bot ext.Bot, updates chan Update) *Dispatcher {
	return &Dispatcher{
		Bot:           bot,
		updates:       updates,
		handlers:      map[int][]Handler{},
		handlerGroups: &[]int{},
	}
}

func (d Dispatcher) Start() {
	for upd := range d.updates {
		d.processUpdate(upd)
	}
}

func (d Dispatcher) processUpdate(update Update) {
	for _, handlerGroupNum := range *d.handlerGroups {
		for _, handler := range d.handlers[handlerGroupNum] {
			if res, err := handler.CheckUpdate(update); res {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							logrus.Error(r)
						}
					}()
					handler.HandleUpdate(update, d)
				}()
				break // move to next group
			} else if err != nil {
				logrus.WithError(err).Error("Failed to parse update")
			}
		}
	}
}

func (d Dispatcher) AddHandler(handler Handler) {
	//*d.handlers = append(*d.handlers, handler)
	d.AddHandlerToGroup(handler, 0)
}

func (d Dispatcher) AddHandlerToGroup(handler Handler, group int) {
	currHandlers, ok := d.handlers[group]
	if !ok {
		*d.handlerGroups = append(*d.handlerGroups, group)
		sort.Ints(*d.handlerGroups)
	}
	d.handlers[group] = append(currHandlers, handler)
}
