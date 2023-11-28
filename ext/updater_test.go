package ext_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	u := ext.NewUpdater(d, nil)

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

func TestUpdaterSupportsWebhookReAdding(t *testing.T) {
	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

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

func TestUpdaterDisallowsEmptyWebhooks(t *testing.T) {
	b := &gotgbot.Bot{
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

	err := u.AddWebhook(b, "", nil)
	if !errors.Is(err, ext.ErrEmptyPath) {
		t.Errorf("Expected an empty path error trying to add an empty webhook : %v", err)
		return
	}
}

func TestUpdaterAllowsWebhookDeletion(t *testing.T) {
	server := basicTestServer(t, map[string]string{
		"getUpdates":    `{}`,
		"deleteWebhook": `{"ok": true, "result": true}`,
	})
	defer server.Close()

	reqOpts := &gotgbot.RequestOpts{
		APIURL: server.URL,
	}

	b := &gotgbot.Bot{
		Token: "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{
			DefaultRequestOpts: reqOpts,
		},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

	err := u.StartPolling(b, &ext.PollingOpts{
		EnableWebhookDeletion: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err != nil {
		t.Errorf("failed to start long poll on first bot: %v", err)
		return
	}
}

func TestUpdaterSupportsTwoPollingBots(t *testing.T) {
	server := basicTestServer(t, map[string]string{
		"getUpdates": `{"ok": true, "result": []}`,
	})
	defer server.Close()

	reqOpts := &gotgbot.RequestOpts{
		APIURL: server.URL,
	}

	b1 := &gotgbot.Bot{
		Token: "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{

			DefaultRequestOpts: reqOpts,
		},
	}
	b2 := &gotgbot.Bot{
		Token: "SOME_OTHER_TOKEN",
		BotClient: &gotgbot.BaseBotClient{
			DefaultRequestOpts: reqOpts,
		},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

	err := u.StartPolling(b1, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err != nil {
		t.Errorf("failed to start long poll on first bot: %v", err)
		return
	}

	// Adding a second time should throw an error
	err = u.StartPolling(b2, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err != nil {
		t.Errorf("should be able to add two different polling bots")
		return
	}
}

func TestUpdaterThrowsErrorWhenSameLongPollAddedTwice(t *testing.T) {
	server := basicTestServer(t, map[string]string{
		"getUpdates": `{"ok": true, "result": []}`,
	})
	defer server.Close()

	reqOpts := &gotgbot.RequestOpts{
		APIURL: server.URL,
	}

	b := &gotgbot.Bot{
		Token: "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{
			DefaultRequestOpts: reqOpts,
		},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

	err := u.StartPolling(b, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err != nil {
		t.Errorf("failed to start long poll: %v", err)
		return
	}

	// Adding a second time should throw an error
	err = u.StartPolling(b, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err == nil {
		t.Errorf("should have failed to start the same long poll twice, but didnt")
		return
	}
}

func TestUpdaterSupportsLongPollReAdding(t *testing.T) {
	server := basicTestServer(t, map[string]string{
		"getUpdates": `{"ok": true, "result": []}`,
	})
	defer server.Close()

	reqOpts := &gotgbot.RequestOpts{
		APIURL: server.URL,
	}

	b := &gotgbot.Bot{
		User:  gotgbot.User{},
		Token: "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{
			DefaultRequestOpts: reqOpts,
		},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{})
	u := ext.NewUpdater(d, nil)

	err := u.StartPolling(b, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{RequestOpts: reqOpts},
	})
	if err != nil {
		t.Errorf("failed to start longpoll: %v", err)
		return
	}

	ok := u.StopBot(b.Token)
	if !ok {
		t.Errorf("failed to stop bot: %v", err)
		return
	}

	// Should be able to re-add the bot now
	err = u.StartPolling(b, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{RequestOpts: reqOpts},
	})
	if err != nil {
		t.Errorf("Failed to re-start a previously removed bot: %v", err)
		return
	}
}

func basicTestServer(t *testing.T, methods map[string]string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathItems := strings.Split(r.URL.Path, "/")
		lastItem := pathItems[len(pathItems)-1]
		t.Logf("Received API call to '%s'", lastItem)

		out, ok := methods[lastItem]
		if ok {
			fmt.Fprint(w, out)
			return
		}

		t.Errorf("Unknown API endpoint: '%s'", lastItem)
		bs, err := json.Marshal(gotgbot.Response{
			Ok:          false,
			ErrorCode:   400,
			Description: "unknown test method: " + r.URL.Path,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
	}))

	return srv
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
	u := ext.NewUpdater(d, nil)

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
