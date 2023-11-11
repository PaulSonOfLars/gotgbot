package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
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

// The UpdateDispatcher interface is used to abstract away common Dispatcher implementations.
// It assumes that all incoming updates come through a JSON channel.
type UpdateDispatcher interface {
	Start(b *gotgbot.Bot, updates <-chan json.RawMessage)
	Stop()
}

// The Dispatcher struct is the default UpdateDispatcher implementation.
// It supports grouping of update handlers, allowing for powerful update handling flows.
// Customise the handling of updates by wrapping the Processor struct.
type Dispatcher struct {
	// Processor defines how to process the raw updates being handled by the Dispatcher.
	// This can be extended to include additional error handling, metrics, etc.
	Processor Processor

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

	// handlers represents all available handlers.
	handlers handlerMapping

	// limiter is how we limit the maximum number of goroutines for handling updates.
	// if nil, this is a limitless dispatcher.
	limiter chan struct{}
	// waitGroup handles the number of running operations to allow for clean shutdowns.
	waitGroup sync.WaitGroup
}

// Ensure compile-time type safety.
var _ UpdateDispatcher = &Dispatcher{}

// DispatcherOpts can be used to configure or override default Dispatcher behaviours.
type DispatcherOpts struct {
	// Processor allows for providing custom Processor interfaces with different behaviours.
	Processor Processor
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

// NewDispatcher creates a new Dispatcher, which process and handles incoming updates from the updates channel.
func NewDispatcher(opts *DispatcherOpts) *Dispatcher {
	var errHandler DispatcherErrorHandler
	var panicHandler DispatcherPanicHandler
	var unhandledErrFunc ErrorFunc
	var errLog *log.Logger

	maxRoutines := DefaultMaxRoutines
	processor := Processor(BaseProcessor{})

	if opts != nil {
		if opts.MaxRoutines != 0 {
			maxRoutines = opts.MaxRoutines
		}
		if opts.Processor != nil {
			processor = opts.Processor
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
		Processor:        processor,
		Error:            errHandler,
		Panic:            panicHandler,
		UnhandledErrFunc: unhandledErrFunc,
		ErrorLog:         errLog,
		handlers:         handlerMapping{},
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
// This is a blocking method; it should be called as a goroutine, such that it can receive incoming updates.
func (d *Dispatcher) Start(b *gotgbot.Bot, updates <-chan json.RawMessage) {
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

			err := d.processRawUpdate(b, upd)
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
	if d.limiter != nil {
		close(d.limiter)
	}
}

// AddHandler adds a new handler to the dispatcher. The dispatcher will call CheckUpdate() to see whether the handler
// should be executed, and then HandleUpdate() to execute it.
func (d *Dispatcher) AddHandler(handler Handler) {
	d.AddHandlerToGroup(handler, 0)
}

// AddHandlerToGroup adds a handler to a specific group; lowest number will be processed first.
func (d *Dispatcher) AddHandlerToGroup(h Handler, group int) {
	d.handlers.add(h, group)
}

// RemoveHandlerFromGroup removes a handler by name from the specified group.
// If multiple handlers have the same name, only the first one is removed.
// Returns true if the handler was successfully removed.
func (d *Dispatcher) RemoveHandlerFromGroup(handlerName string, group int) bool {
	return d.handlers.remove(handlerName, group)
}

// RemoveGroup removes an entire group from the dispatcher's processing.
// If group can't be found, this is a noop.
func (d *Dispatcher) RemoveGroup(group int) bool {
	return d.handlers.removeGroup(group)
}

// processRawUpdate takes a JSON update to be unmarshalled and processed by Dispatcher.ProcessUpdate.
func (d *Dispatcher) processRawUpdate(b *gotgbot.Bot, r json.RawMessage) error {
	var upd gotgbot.Update
	if err := json.Unmarshal(r, &upd); err != nil {
		return fmt.Errorf("failed to unmarshal update: %w", err)
	}

	return d.ProcessUpdate(b, &upd, nil)
}

// ProcessUpdate iterates over the list of groups to execute the matching handlers.
// This is also where we recover from any panics that are thrown by user code, to avoid taking down the bot.
func (d *Dispatcher) ProcessUpdate(b *gotgbot.Bot, u *gotgbot.Update, data map[string]interface{}) (err error) {
	ctx := NewContext(u, data)

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

	err = d.Processor.ProcessUpdate(d, b, ctx)
	// We don't inline this, because we want to make sure that the defer function can override the error in the case of
	// a panic.
	return err
}

func (d *Dispatcher) iterateOverHandlerGroups(b *gotgbot.Bot, ctx *Context) error {
	for _, groups := range d.handlers.getGroups() {
		for _, handler := range groups {
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
