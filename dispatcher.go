package gotgbot

import (
	"runtime/debug"
	"sort"

	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	Bot           ext.Bot
	MaxRoutines   int
	updates       chan *Update
	handlers      map[int][]Handler
	handlerGroups *[]int
}

const DefaultMaxDispatcherRoutines = 50

func NewDispatcher(bot ext.Bot, updates chan *Update) *Dispatcher {
	return &Dispatcher{
		Bot:           bot,
		MaxRoutines:   DefaultMaxDispatcherRoutines,
		updates:       updates,
		handlers:      map[int][]Handler{},
		handlerGroups: &[]int{},
	}
}

func (d Dispatcher) Start() {
	limiter := make(chan struct{}, d.MaxRoutines)
	for upd := range d.updates {
		limiter <- struct{}{}
		go func(upd *Update) {
			d.processUpdate(upd)
			<-limiter
		}(upd)
	}
}

type EndGroups struct{}
type ContinueGroups struct{}

func (eg EndGroups) Error() string      { return "Group iteration ended" }
func (eg ContinueGroups) Error() string { return "Group iteration has continued" }

func (d Dispatcher) processUpdate(update *Update) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r)
			debug.PrintStack()
		}
	}()

	for _, groupNum := range *d.handlerGroups {
		for _, handler := range d.handlers[groupNum] {
			if res, err := handler.CheckUpdate(update); res {
				err := handler.HandleUpdate(update, d)
				if err != nil {
					switch err.(type) {
					case EndGroups:
						return
					case ContinueGroups:
						continue
					default:
						logrus.Warning(err.Error())
					}
				}
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
