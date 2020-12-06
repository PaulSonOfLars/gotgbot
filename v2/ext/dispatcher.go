package ext

import (
	"encoding/json"
	"errors"
	"os"
	"runtime/debug"
	"sort"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const DefaultMaxRoutines = 50

type Dispatcher struct {
	// Updates receives all incoming updates from the updater'
	Updates chan json.RawMessage
	// MaxRoutines limits the maximum number of goroutines available for update handling.
	// If value is < 0, no limiting is applied.
	MaxRoutines int
	// Error handles any errors that occur during handler execution.
	Error func(ctx *Context, err error)
	// Panic handles any panics that occur during handler execution.
	//If no panic is defined, the stack is printed to stderr.
	Panic func(ctx *Context, stack []byte)

	// handlerGroups represents the list of available handler groups, numerically sorted.
	handlerGroups []int // TODO: should be pointer?
	// handlers represents all available handles, split into groups (see handlerGroups).
	handlers map[int][]Handler
}

func NewDispatcher(updates chan json.RawMessage, maxRoutines int) *Dispatcher {
	return &Dispatcher{
		Updates:     updates,
		MaxRoutines: maxRoutines,
		handlers:    make(map[int][]Handler),
	}
}

// Start to handle incoming updates
func (d *Dispatcher) Start(b *gotgbot.Bot) {
	if d.MaxRoutines < 0 {
		d.limitlessDispatcher(b)
		return
	}

	routines := d.MaxRoutines
	if routines == 0 {
		routines = DefaultMaxRoutines
	}
	d.limitedDispatcher(b, routines)
}

func (d *Dispatcher) limitedDispatcher(b *gotgbot.Bot, routines int) {
	limiter := make(chan struct{}, routines)
	for upd := range d.Updates {
		// Send empty data to limiter.
		// if limiter buffer is full, it block until it another update ends.
		limiter <- struct{}{}
		go func(upd json.RawMessage) {
			d.ProcessRawUpdate(b, upd)
			<-limiter
		}(upd)
	}
}

func (d *Dispatcher) limitlessDispatcher(b *gotgbot.Bot) {
	for upd := range d.Updates {
		go func(upd json.RawMessage) {
			d.ProcessRawUpdate(b, upd)
		}(upd)
	}
}

// AddHandler adds a new handler to the dispatcher. The dispatcher will call CheckUpdate() to see whether the handler
// should be executed, and then HandleUpdate() to execute it.
func (d *Dispatcher) AddHandler(handler Handler) {
	d.AddHandlerToGroup(handler, 0)
}

// AddHandlerToGroup adds a handler to a specific group; lowest number will be processed first.
func (d *Dispatcher) AddHandlerToGroup(handler Handler, group int) {
	currHandlers, ok := d.handlers[group]
	if !ok {
		d.handlerGroups = append(d.handlerGroups, group)
		sort.Ints(d.handlerGroups)
	}
	d.handlers[group] = append(currHandlers, handler)
}

var EndGroups = errors.New("group iteration ended")
var ContinueGroups = errors.New("group iteration continued")

func (d *Dispatcher) ProcessRawUpdate(b *gotgbot.Bot, r json.RawMessage) {
	var upd gotgbot.Update
	if err := json.Unmarshal(r, &upd); err != nil {
		// todo: improve logging
		os.Stderr.WriteString(err.Error())
		return
	}

	d.ProcessUpdate(b, &upd)
}

func (d *Dispatcher) ProcessUpdate(b *gotgbot.Bot, update *gotgbot.Update) {
	var ctx *Context

	defer func() {
		if r := recover(); r != nil {
			if d.Panic != nil {
				d.Panic(ctx, debug.Stack())
				return
			}

			debug.PrintStack()
		}
	}()

	for _, groupNum := range d.handlerGroups {
		for _, handler := range d.handlers[groupNum] {
			if !handler.CheckUpdate(b, update) {
				return
			}

			if ctx == nil {
				ctx = NewContext(b, update)
			}

			err := handler.HandleUpdate(ctx)
			if err != nil {
				switch err {
				case EndGroups:
					return
				case ContinueGroups:
					continue
				default:
					if d.Error != nil {
						d.Error(ctx, err)
					}
				}
			}
			break // move to next group
		}
	}

	return
}
