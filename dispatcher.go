package gotgbot

import (
	"encoding/json"
	"runtime/debug"
	"sort"

	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
)

type RawUpdate json.RawMessage

type Dispatcher struct {
	Bot           *ext.Bot
	MaxRoutines   int
	updates       chan *RawUpdate
	handlers      map[int][]Handler
	handlerGroups *[]int
}

const DefaultMaxDispatcherRoutines = 50

func NewDispatcher(bot *ext.Bot, updates chan *RawUpdate) *Dispatcher {
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
		select {
		case limiter <- struct{}{}:
		default:
			logrus.Debugf("update dispatcher has reached limit of %d", d.MaxRoutines)
			limiter <- struct{}{} // make sure to send anyway
		}
		go func(upd *RawUpdate) {
			d.processUpdate(upd)
			<-limiter
		}(upd)
	}
}

type EndGroups struct{}
type ContinueGroups struct{}

func (eg EndGroups) Error() string      { return "Group iteration ended" }
func (eg ContinueGroups) Error() string { return "Group iteration has continued" }

func (d Dispatcher) processUpdate(upd *RawUpdate) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r)
			debug.PrintStack()
		}
	}()

	update, err := initUpdate(*upd, *d.Bot)
	if err != nil {
		logrus.WithError(err).Error("failed to init update while processing")
		return
	}
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
				logrus.WithError(err).Error("failed to check update while processing")
				return
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
