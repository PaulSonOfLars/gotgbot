package ext

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// botData Keeps track of the necessary update channels for each gotgbot.Bot.
type botData struct {
	// bot represents the bot for which this data is relevant.
	bot *gotgbot.Bot

	// updateChan represents the incoming updates channel.
	updateChan chan json.RawMessage
	// updateWriterControl is used to count the number of current writers on the update channel.
	// This is required to ensure that we can safely close the channel, and thus stop processing incoming updates.
	// While this remains non-zero, it is unsafe to close the update channel.
	updateWriterControl *sync.WaitGroup
	// stopUpdates allows us to close the stopUpdates loop.
	stopUpdates chan struct{}

	// urlPath defines the incoming webhook URL path for this bot.
	urlPath string
	// webhookSecret stores the webhook secret for this bot.
	webhookSecret string
}

// botMapping Ensures that all botData is stored in a thread-safe manner.
type botMapping struct {
	// mapping keeps track of the data required for each bot. The key is the bot token.
	mapping map[string]botData
	// mux attempts to keep the botMapping data concurrency-safe.
	mux sync.RWMutex
	// urlMapping allows us to keep track of the webhook urls that are in use.
	urlMapping map[string]string

	// errFunc fills the same purpose as Updater.UnhandledErrFunc.
	errFunc ErrorFunc
	// errorLog fills the same purpose as Updater.ErrorLog.
	errorLog *log.Logger
}

var ErrBotAlreadyExists = errors.New("bot already exists in bot mapping")
var ErrBotUrlPathAlreadyExists = errors.New("url path already exists in bot mapping")

// addBot Adds a new bot to the botMapping structure.
// Pass an empty urlPath/webhookSecret if using polling instead of webhooks.
func (m *botMapping) addBot(b *gotgbot.Bot, urlPath string, webhookSecret string) (*botData, error) {
	// Clean up the URLPath such that it remains consistent.
	urlPath = strings.TrimPrefix(urlPath, "/")

	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}
	if m.urlMapping == nil {
		m.urlMapping = make(map[string]string)
	}

	if _, ok := m.mapping[b.Token]; ok {
		return nil, ErrBotAlreadyExists
	}

	if _, ok := m.urlMapping[urlPath]; urlPath != "" && ok {
		return nil, ErrBotUrlPathAlreadyExists
	}

	bData := botData{
		bot:                 b,
		updateChan:          make(chan json.RawMessage),
		stopUpdates:         make(chan struct{}),
		updateWriterControl: &sync.WaitGroup{},
		urlPath:             urlPath,
		webhookSecret:       webhookSecret,
	}

	m.mapping[bData.bot.Token] = bData
	m.urlMapping[bData.urlPath] = bData.bot.Token
	return &bData, nil
}

func (m *botMapping) removeBot(token string) (botData, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()

	bData, ok := m.mapping[token]
	if !ok {
		return botData{}, false
	}

	delete(m.mapping, token)
	delete(m.urlMapping, bData.urlPath)
	return bData, true
}

func (m *botMapping) removeAllBots() []botData {
	m.mux.Lock()
	defer m.mux.Unlock()

	bots := make([]botData, 0, len(m.mapping))
	for key, bData := range m.mapping {
		bots = append(bots, bData)
		delete(m.mapping, key)
		delete(m.urlMapping, bData.urlPath)
	}
	return bots
}

func (m *botMapping) getBots() []botData {
	m.mux.RLock()
	defer m.mux.RUnlock()

	bots := make([]botData, 0, len(m.mapping))
	for _, bData := range m.mapping {
		bots = append(bots, bData)
	}
	return bots
}

func (m *botMapping) getBot(token string) (botData, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	bData, ok := m.mapping[token]
	return bData, ok
}

func (m *botMapping) getBotFromURL(urlPath string) (botData, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	token, ok := m.urlMapping[urlPath]
	if !ok {
		return botData{}, ok
	}

	bData, ok := m.mapping[token]
	return bData, ok
}

func (m *botMapping) getHandlerFunc(prefix string) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "*" {
			if r.ProtoAtLeast(1, 1) {
				w.Header().Set("Connection", "close")
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, ok := m.getBotFromURL(strings.TrimPrefix(r.URL.Path, prefix))
		if !ok {
			// If we don't recognise the URL, we return a 404.
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b.updateWriterControl.Add(1)
		defer b.updateWriterControl.Done()

		if b.shouldStopUpdates() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		headerSecret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if b.webhookSecret != "" && b.webhookSecret != headerSecret {
			// Drop any updates from invalid secret tokens.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			if m.errFunc != nil {
				m.errFunc(err)
			} else {
				m.logf("Failed to read incoming update contents: %s", err.Error())
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b.updateChan <- bytes
	}
}

func (m *botMapping) logf(format string, args ...interface{}) {
	if m.errorLog != nil {
		m.errorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (b *botData) stop() {
	// Close stopUpdates loops first, to ensure any updates currently being polled have the time to be sent to the updateChan.
	if b.stopUpdates != nil {
		close(b.stopUpdates)
	}

	// Wait for all writers to finish writing to the updateChannel
	b.updateWriterControl.Wait()

	// Then, close the updates channel.
	close(b.updateChan)
}

func (b *botData) shouldStopUpdates() bool {
	select {
	case <-b.stopUpdates:
		// if anything comes in on the closing channel, we know the channel is closed.
		return true
	default:
		// otherwise, continue as usual
		return false
	}
}
