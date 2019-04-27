package gotgbot

import (
	"encoding/json"
	"runtime/debug"
	"sort"

	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/sirupsen/logrus"
)

//RawUpdate alias to json.RawMessage
type RawUpdate json.RawMessage

//Dispatcher Store data for the dispatcher to work as expected; such as the incoming update channel,
// handler mappings and maximum number of goroutines allowed to be run at once
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

//Start Begin dispatching updates
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

//AddHandler adds a new handler to the dispatcher. The dispatcher will call CheckUpdate() to see whether the handler
// should be executed, and then HandleUpdate() to execute it.
func (d Dispatcher) AddHandler(handler Handler) {
	//*d.handlers = append(*d.handlers, handler)
	d.AddHandlerToGroup(handler, 0)
}

//AddHandlerToGroup adds a handler to a specific group; lowest number will be processed first.
func (d Dispatcher) AddHandlerToGroup(handler Handler, group int) {
	currHandlers, ok := d.handlers[group]
	if !ok {
		*d.handlerGroups = append(*d.handlerGroups, group)
		sort.Ints(*d.handlerGroups)
	}
	d.handlers[group] = append(currHandlers, handler)
}
