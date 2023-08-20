package ext

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Processor interface {
	ProcessUpdate(d *Dispatcher, b *gotgbot.Bot, u *gotgbot.Update, data map[string]interface{}) error
}

var _ Processor = BaseProcessor{}

type BaseProcessor struct{}

// ProcessUpdate iterates over the list of groups to execute the matching handlers.
func (bp BaseProcessor) ProcessUpdate(d *Dispatcher, b *gotgbot.Bot, u *gotgbot.Update, data map[string]interface{}) (err error) {
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

	err = d.iterateOverHandlerGroups(b, ctx)
	// We don't inline this, because we want to make sure that the defer function can override the error in the case of
	// a panic.
	return err
}
