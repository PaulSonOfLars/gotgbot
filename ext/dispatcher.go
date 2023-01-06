package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"strings"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	ErrPanicRecovered          = errors.New("panic recovered")
	ErrUnknownDispatcherAction = errors.New("unknown dispatcher action")
)

const DefaultMaxRoutines = 50

type (
	// DispatcherErrorHandler allows for handling the returned errors from matched handlers.
	// It takes the non-nil error returned by the handler.
	DispatcherErrorHandler func(b *gotgbot.Bot, ctx *Context, err error) DispatcherAction
	// DispatcherPanicHandler allows for handling goroutine panics, where the 'r' value contains the reason for the panic.
	DispatcherPanicHandler func(b *gotgbot.Bot, ctx *Context, r interface{})
)

type DispatcherAction string

const (
	// DispatcherActionNoop stops iteration of current group and moves to the next one.
	// This is the default action, and the same as would happen if the handler had completed successfully.
	DispatcherActionNoop DispatcherAction = "noop"
	// DispatcherActionContinueGroups continues iterating over current group as if the current handler did not match.
	// Functionally the same as returning ContinueGroups.
	DispatcherActionContinueGroups DispatcherAction = "continue-groups"
	// DispatcherActionEndGroups ends all group iteration.
	// Functionally the same as returning EndGroups.
	DispatcherActionEndGroups DispatcherAction = "end-groups"
)

var (
	EndGroups      = errors.New("group iteration ended")
	ContinueGroups = errors.New("group iteration continued")
)

type Dispatcher struct {
	// Error handles any errors that are returned by matched handlers.
	// The return type determines how to proceed with the current group iteration.
	// The default action is DispatcherActionNoop, which will simply move to next group as expected.
	Error DispatcherErrorHandler
	// Panic handles any panics that occur during handler execution.
	// Panics from handlers are automatically recovered to ensure bot stability. Once recovered, this method is called
	// and is left to determine how to log or handle the errors.
	// If this field is nil, the error will be passed to UnhandledErrFunc.
	Panic DispatcherPanicHandler

	// UnhandledErrFunc provides more flexibility for dealing with unhandled update processing errors.
	// This includes errors when unmarshalling updates, unhandled panics during handler executions, or unknown
	// dispatcher actions.
	// If nil, the error goes to ErrorLog.
	UnhandledErrFunc ErrorFunc
	// ErrorLog specifies an optional logger for unexpected behavior from handlers.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	// handlerGroups represents the list of available handler groups, numerically sorted.
	handlerGroups []int
	// handlers represents all available handles, split into groups (see handlerGroups).
	handlers map[int][]Handler

	// limiter is how we limit the maximum number of goroutines for handling updates.
	// if nil, this is a limitless dispatcher.
	limiter chan struct{}
	// waitGroup handles the number of running operations to allow for clean shutdowns.
	waitGroup sync.WaitGroup
}

type DispatcherOpts struct {
	// Error handles any errors that occur during handler execution.
	// More info at Dispatcher.Error.
	Error DispatcherErrorHandler
	// Panic handles any panics that occur during handler execution.
	// If no panic handlers are defined, the stack is logged to ErrorLog.
	// More info at Dispatcher.Panic.
	Panic DispatcherPanicHandler

	// UnhandledErrFunc provides more flexibility for dealing with unhandled update processing errors.
	// This includes errors when unmarshalling updates, unhandled panics during handler executions, or unknown
	// dispatcher actions.
	// If nil, the error goes to ErrorLog.
	UnhandledErrFunc ErrorFunc
	// ErrorLog specifies an optional logger for unexpected behavior from handlers.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	// MaxRoutines is used to decide how to limit the number of goroutines spawned by the dispatcher.
	// This defines how many updates can be processed at the same time.
	// If MaxRoutines == 0, DefaultMaxRoutines is used instead.
	// If MaxRoutines < 0, no limits are imposed.
	// If MaxRoutines > 0, that value is used.
	MaxRoutines int
}

// NewDispatcher creates a new dispatcher, which process and handles incoming updates from the updates channel.
func NewDispatcher(opts *DispatcherOpts) *Dispatcher {
	var errHandler DispatcherErrorHandler
	var panicHandler DispatcherPanicHandler
	var unhandledErrFunc ErrorFunc
	var errLog *log.Logger

	maxRoutines := DefaultMaxRoutines

	if opts != nil {
		if opts.MaxRoutines != 0 {
			maxRoutines = opts.MaxRoutines
		}

		errHandler = opts.Error
		panicHandler = opts.Panic
		unhandledErrFunc = opts.UnhandledErrFunc
		errLog = opts.ErrorLog
	}

	var limiter chan struct{}
	// if maxRoutines < 0, we use a limitless dispatcher. (limiter == nil)
	if maxRoutines >= 0 {
		if maxRoutines == 0 {
			maxRoutines = DefaultMaxRoutines
		}

		limiter = make(chan struct{}, maxRoutines)
	}

	return &Dispatcher{
		Error:            errHandler,
		Panic:            panicHandler,
		UnhandledErrFunc: unhandledErrFunc,
		ErrorLog:         errLog,
		handlers:         make(map[int][]Handler),
		limiter:          limiter,
		waitGroup:        sync.WaitGroup{},
	}
}

func (d *Dispatcher) logf(format string, args ...interface{}) {
	if d.ErrorLog != nil {
		d.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// CurrentUsage returns the current number of concurrently processing updates.
func (d *Dispatcher) CurrentUsage() int {
	return len(d.limiter)
}

// MaxUsage returns the maximum number of concurrently processing updates.
func (d *Dispatcher) MaxUsage() int {
	return cap(d.limiter)
}

// Start to handle incoming updates.
func (d *Dispatcher) Start(b *gotgbot.Bot, updates chan json.RawMessage) {
	// Listen to updates as they come in from the updater.
	for upd := range updates {
		d.waitGroup.Add(1)

		// If a limiter has been set, we use it to control the number of concurrent updates being processed.
		if d.limiter != nil {
			// Send data to limiter.
			// If limiter buffer is full, this will block, until another update finishes processing.
			d.limiter <- struct{}{}
		}

		go func(upd json.RawMessage) {
			// We defer here so that whatever happens, we can clean up the dispatcher.
			defer func() {
				if d.limiter != nil {
					// Pop an item from the limiter, allowing another update to process.
					<-d.limiter
				}
				d.waitGroup.Done()
			}()

			err := d.ProcessRawUpdate(b, upd)
			if err != nil {
				if d.UnhandledErrFunc != nil {
					d.UnhandledErrFunc(err)
				} else {
					d.logf("Failed to process update: %s", err.Error())
				}
			}

		}(upd)
	}
}

// Stop waits for all currently processing updates to finish, and then returns.
func (d *Dispatcher) Stop() {
	d.waitGroup.Wait()
	close(d.limiter)
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

func (d *Dispatcher) ProcessRawUpdate(b *gotgbot.Bot, r json.RawMessage) error {
	var upd gotgbot.Update
	if err := json.Unmarshal(r, &upd); err != nil {
		return fmt.Errorf("failed to unmarshal update: %w", err)
	}

	return d.ProcessUpdate(b, &upd, nil)
}

// ProcessUpdate iterates over the list of groups to execute the matching handlers.
func (d *Dispatcher) ProcessUpdate(b *gotgbot.Bot, update *gotgbot.Update, data map[string]interface{}) (err error) {
	ctx := NewContext(update, data)

	defer func() {
		if r := recover(); r != nil {
			// If a panic handler is defined, handle the error.
			if d.Panic != nil {
				d.Panic(b, ctx, r)
				return

			} else {
				// Otherwise, create an error from the panic, and return it.
				err = fmt.Errorf("%w: %v\n%s", ErrPanicRecovered, r, cleanedStack())
				return
			}
		}
	}()

	err = d.iterateOverHandlerGroups(b, ctx)
	// We don't inline this, because we want to make sure that the defer function can override the error in the case of
	// a panic.
	return err
}

func (d *Dispatcher) iterateOverHandlerGroups(b *gotgbot.Bot, ctx *Context) error {
	for _, groupNum := range d.handlerGroups {
		for _, handler := range d.handlers[groupNum] {
			if !handler.CheckUpdate(b, ctx) {
				// Handler filter doesn't match this update; continue.
				continue
			}

			err := handler.HandleUpdate(b, ctx)
			if err != nil {
				if errors.Is(err, ContinueGroups) {
					// Continue handling current group.
					continue

				} else if errors.Is(err, EndGroups) {
					// Stop all group handling.
					return nil

				} else {
					action := DispatcherActionNoop
					if d.Error != nil {
						action = d.Error(b, ctx, err)
					}

					switch action {
					case DispatcherActionNoop:
						// Move on to next group; same action as if group had been successful.
					case DispatcherActionContinueGroups:
						// Continue handling current group.
						continue
					case DispatcherActionEndGroups:
						// Stop all group handling.
						return nil
					default:
						return fmt.Errorf("%w: '%s', ending groups here", ErrUnknownDispatcherAction, action)
					}
				}
			}

			// Handler matched this update, move to next group by default.
			break
		}
	}
	return nil
}

// cleanedStack obtains a "cleaned" version of the stack trace which doesn't point the last few lines to the library.
// This is because historically, people see the library in the stack trace, and immediately blame it, when in fact it is
// recovering their errors.
func cleanedStack() string {
	return strings.Join(strings.Split(string(debug.Stack()), "\n")[4:], "\n")
}
