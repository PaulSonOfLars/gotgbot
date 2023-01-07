package ext

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDispatcherStop(t *testing.T) {
	d := NewDispatcher(nil)

	go d.Start(nil, make(chan json.RawMessage))

	waited := false
	d.waitGroup.Add(1)
	d.limiter <- struct{}{}

	// imitate a "wait"
	go func() {
		time.Sleep(time.Second)
		waited = true

		<-d.limiter
		d.waitGroup.Done()
	}()

	d.Stop()
	if !waited {
		t.Errorf("Dispatcher was stopped before the updates were done being handled.")
	}

}

func TestLimitedDispatcherStop(t *testing.T) {
	d := NewDispatcher(&DispatcherOpts{
		MaxRoutines: DefaultMaxRoutines,
	})

	if cap(d.limiter) != DefaultMaxRoutines {
		t.Errorf("Expected limiter to be of max size %d, got %d", DefaultMaxRoutines, cap(d.limiter))
	}

	go d.Start(nil, make(chan json.RawMessage))
	d.Stop() // ensure no panics
}

func TestUnlimitedDispatcherStop(t *testing.T) {
	d := NewDispatcher(&DispatcherOpts{
		MaxRoutines: -1,
	})

	if d.limiter != nil {
		t.Errorf("Expected limiter to be nil for unlimited dispatcher")
	}

	go d.Start(nil, make(chan json.RawMessage))
	d.Stop() // ensure no panics
}
