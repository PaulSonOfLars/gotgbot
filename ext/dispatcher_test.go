package ext

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
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

func BenchmarkDispatcher(b *testing.B) {
	d := NewDispatcher(nil)

	wg := sync.WaitGroup{}
	d.AddHandler(DummyHandler{F: func(b *gotgbot.Bot, ctx *Context) error {
		wg.Done()
		return nil
	}})

	updateChan := make(chan json.RawMessage)

	go d.Start(&gotgbot.Bot{}, updateChan)

	upd, err := json.Marshal(gotgbot.Update{Message: &gotgbot.Message{Text: "test"}})
	if err != nil {
		b.Fatalf("failed to marshal test msg: %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() { updateChan <- upd }()
	}

	wg.Wait()
	d.Stop()
}
