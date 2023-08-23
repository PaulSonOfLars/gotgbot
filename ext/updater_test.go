package ext_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

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

	opts := ext.WebhookOpts{}
	for i := 0; i < numBot; i++ {
		token := strconv.Itoa(i)
		err := u.AddWebhook(&gotgbot.Bot{
			BotClient: &gotgbot.BaseBotClient{
				Token: token,
			},
		}, token, opts)
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
