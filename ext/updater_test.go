package ext_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func TestUpdaterThrowsErrorWhenSameWebhookAddedTwice(t *testing.T) {
	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: d})

	err := u.AddWebhook(b, "test", nil)
	if err != nil {
		t.Errorf("failed to add webhook: %v", err)
		return
	}

	// Adding a second time should throw an error
	err = u.AddWebhook(b, "test", nil)
	if err == nil {
		t.Errorf("should have failed to add the same webhook twice, but didnt")
		return
	}
}

func TestUpdaterSupportsReAdding(t *testing.T) {
	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: d})

	err := u.AddWebhook(b, "test", nil)
	if err != nil {
		t.Errorf("failed to add webhook: %v", err)
		return
	}

	ok := u.StopBot(b.Token)
	if !ok {
		t.Errorf("failed to stop bot: %v", err)
		return
	}

	// Should be able to re-add the bot now
	err = u.AddWebhook(b, "test", nil)
	if err != nil {
		t.Errorf("Failed to re-add a previously removed bot: %v", err)
		return
	}
}

func BenchmarkUpdaterMultibots(b *testing.B) {
	// Note: This benchmark skips the JSON marshal/unmarshal steps to get an accurate idea of how well the multibot
	//  features work.
	b.Run("single", func(b *testing.B) {
		benchmarkUpdaterWithNBots(b, 1)
	})
	b.Run("ten", func(b *testing.B) {
		benchmarkUpdaterWithNBots(b, 10)
	})
	b.Run("hundred", func(b *testing.B) {
		benchmarkUpdaterWithNBots(b, 100)
	})
	b.Run("thousand", func(b *testing.B) {
		benchmarkUpdaterWithNBots(b, 1000)
	})
}

func benchmarkUpdaterWithNBots(b *testing.B, numBot int) {
	d := ext.NewDispatcher(nil)
	u := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: d})

	wg := sync.WaitGroup{}
	d.AddHandler(ext.DummyHandler{F: func(b *gotgbot.Bot, ctx *ext.Context) error {
		wg.Done()
		return nil
	}})

	for i := 0; i < numBot; i++ {
		token := strconv.Itoa(i)
		err := u.AddWebhook(&gotgbot.Bot{
			Token:     token,
			BotClient: &gotgbot.BaseBotClient{},
		}, token, nil)
		if err != nil {
			b.Fatalf("failed to add webhook for bot: %s", err.Error())
		}
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		token := strconv.Itoa(i % numBot)

		go func() {
			err := u.InjectUpdate(token, gotgbot.Update{Message: &gotgbot.Message{Text: "test"}})
			if err != nil {
				b.Error("failed to send manual update: %w", err)
				b.Fail()
			}
		}()
	}

	wg.Wait()
	d.Stop()
}
