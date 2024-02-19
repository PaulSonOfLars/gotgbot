package ext_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
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

// This test is a bit strange, as it tries to test concurrency events. Which means it relies on awkward timing
// situations. To try and mitigate this, we run it multiple times in parallel; hoping this will help catch any issues.
// The idea is that we want to be able to test getUpdates both BEFORE as well as AFTER the bot has been stopped.
// Execution flow is:
// - polling is started
// - getUpdates receives "stop" message. Keeps running on loop with a 1s delay.
// - dispatcher processes stop message; calls updater.stopBot, thus closing channels, removing from bot map, and stopping dispatcher.
// - getUpdates receives message again, 1s later; but updater+dispatcher channels are already closed by then.
// - IF not concurrently safe; we get a panic. Else, works fine!
// - Since updater channels are closed, messages should not be processed.
//
// NOTE: IF THIS FAILS, IT IS NOT A FLAKE! Investigate!
func TestUpdaterPollingConcurrency(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("run_%d", i), func(t *testing.T) {
			concurrentTest(t)
		})
	}
}

func concurrentTest(t *testing.T) {
	// we run it in parallel so that we get all the perks
	t.Parallel()

	delay := time.Second
	server := basicTestServer(t, map[string]*testEndpoint{
		"getUpdates": {
			delay: delay,
			replies: []string{
				`{"ok": true, "result": [{"message": {"text": "stop"}}]}`,
			},
			reply: `{"ok": true, "result": []}`,
		},
		"deleteWebhook": {reply: `{"ok": true, "result": true}`},
	})
	defer server.Close()

	reqOpts := &gotgbot.RequestOpts{
		APIURL:  server.URL,
		Timeout: delay * 2,
	}

	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	d := ext.NewDispatcher(&ext.DispatcherOpts{MaxRoutines: 1})
	u := ext.NewUpdater(d, nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	d.AddHandler(handlers.NewMessage(message.Contains("stop"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		if !u.StopBot(b.Token) {
			t.Errorf("Could not stop bot!")
		}
		// Make sure we mark this method as having run only once.
		// (if run twice, we get a panic)
		wg.Done()
		return nil
	}))

	err := u.StartPolling(b, &ext.PollingOpts{
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: reqOpts,
		},
	})
	if err != nil {
		t.Errorf("failed to start polling: %v", err)
		return
	}

	wg.Wait()
	time.Sleep(delay * 2)
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

func TestUpdater_GetHandlerFunc(t *testing.T) {
	b := &gotgbot.Bot{
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	type args struct {
		urlPath       string
		opts          *ext.AddWebhookOpts
		httpResponse  int
		handlerPrefix string
		requestPath   string // Should start with '/'
		headers       map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple path",
			args: args{
				urlPath:       "123:hello",
				httpResponse:  http.StatusOK,
				handlerPrefix: "/",
				requestPath:   "/123:hello",
			},
		}, {
			name: "slash prefixed path",
			args: args{
				urlPath:       "/123:hello",
				httpResponse:  http.StatusOK,
				handlerPrefix: "/",
				requestPath:   "/123:hello",
			},
		}, {
			name: "using subpath",
			args: args{
				urlPath:       "123:hello",
				httpResponse:  http.StatusOK,
				handlerPrefix: "/test/",
				requestPath:   "/test/123:hello",
			},
		}, {
			name: "unknown path",
			args: args{
				urlPath:       "123:hello",
				httpResponse:  http.StatusNotFound,
				handlerPrefix: "/",
				requestPath:   "/this-path-doesnt-exist",
			},
		}, {
			name: "missing secret token",
			args: args{
				urlPath: "123:hello",
				opts: &ext.AddWebhookOpts{
					SecretToken: "secret",
				},
				httpResponse:  http.StatusUnauthorized,
				handlerPrefix: "/",
				requestPath:   "/123:hello",
			},
		}, {
			name: "matching secret token",
			args: args{
				urlPath: "123:hello",
				opts: &ext.AddWebhookOpts{
					SecretToken: "secret",
				},
				httpResponse:  http.StatusOK,
				handlerPrefix: "/",
				requestPath:   "/123:hello",
				headers: map[string]string{
					"X-Telegram-Bot-Api-Secret-Token": "secret",
				},
			},
		}, {
			name: "invalid secret token",
			args: args{
				urlPath: "123:hello",
				opts: &ext.AddWebhookOpts{
					SecretToken: "secret",
				},
				httpResponse:  http.StatusUnauthorized,
				handlerPrefix: "/",
				requestPath:   "/123:hello",
				headers: map[string]string{
					"X-Telegram-Bot-Api-Secret-Token": "wrong",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ext.NewDispatcher(nil)
			u := ext.NewUpdater(d, nil)

			if err := u.AddWebhook(b, tt.args.urlPath, tt.args.opts); err != nil {
				t.Errorf("failed to add webhook: %v", err)
				return
			}

			s := httptest.NewServer(u.GetHandlerFunc(tt.args.handlerPrefix))
			url := s.URL + tt.args.requestPath
			// We pass {} to satisfy JSON handling
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, strings.NewReader("{}"))
			if err != nil {
				t.Errorf("Failed to build request, should have worked: %v", err.Error())
				return
			}

			for k, v := range tt.args.headers {
				req.Header.Set(k, v)
			}

			r, err := s.Client().Do(req)
			if err != nil {
				t.Fatal()
			}

			defer r.Body.Close()
			if r.StatusCode != tt.args.httpResponse {
				t.Errorf("Expected code %d, got %d", tt.args.httpResponse, r.StatusCode)
				return
			}
		})
	}
}

func TestUpdaterAllowsWebhookDeletion(t *testing.T) {
	server := basicTestServer(t, map[string]*testEndpoint{
		"getUpdates":    {reply: `{"ok": true}`},
		"deleteWebhook": {reply: `{"ok": true, "result": true}`},
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

	err = u.Stop()
	if err != nil {
		t.Errorf("failed to stop updater: %v", err)
		return
	}
}

func TestUpdaterSupportsTwoPollingBots(t *testing.T) {
	server := basicTestServer(t, map[string]*testEndpoint{
		"getUpdates": {reply: `{"ok": true, "result": []}`},
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

	err = u.Stop()
	if err != nil {
		t.Errorf("failed to stop updater: %v", err)
		return
	}
}

func TestUpdaterThrowsErrorWhenSameLongPollAddedTwice(t *testing.T) {
	server := basicTestServer(t, map[string]*testEndpoint{
		"getUpdates": {reply: `{"ok": true, "result": []}`},
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

	err = u.Stop()
	if err != nil {
		t.Errorf("failed to stop updater: %v", err)
		return
	}
}

func TestUpdaterSupportsLongPollReAdding(t *testing.T) {
	server := basicTestServer(t, map[string]*testEndpoint{
		"getUpdates": {reply: `{"ok": true, "result": []}`},
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

	err = u.Stop()
	if err != nil {
		t.Errorf("failed to stop updater: %v", err)
		return
	}
}

type testEndpoint struct {
	delay time.Duration
	// Will reply these until we run out of replies, at which point we repeat "reply"
	replies []string
	idx     atomic.Int32
	// default reply
	reply string
}

func basicTestServer(t *testing.T, methods map[string]*testEndpoint) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathItems := strings.Split(r.URL.Path, "/")
		lastItem := pathItems[len(pathItems)-1]
		t.Logf("Received API call to '%s'", lastItem)

		out, ok := methods[lastItem]
		if ok {
			if out.delay != 0 {
				time.Sleep(out.delay)
			}
			count := int(out.idx.Add(1) - 1)
			if len(out.replies) != 0 && len(out.replies) > count {
				fmt.Fprint(w, out.replies[count])
			} else {
				fmt.Fprint(w, out.reply)
			}
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
