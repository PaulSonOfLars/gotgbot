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
	// polling allows us to close the polling loop.
	polling chan struct{}
	// urlPath defines the incoming webhook URL path for this bot.
	urlPath string
	// webhookSecret stores the webhook secret for this bot
	webhookSecret string
}

// botMapping Ensures that all botData is stored in a thread-safe manner.
type botMapping struct {
	// mapping keeps track of the data required for each bot. The key is the bot token.
	mapping map[string]botData
	// mux attempts to keep the botMapping data concurrency-safe.
	mux sync.RWMutex

	// errFunc fills the same purpose as Updater.UnhandledErrFunc.
	errFunc ErrorFunc
	// errorLog fills the same purpose as Updater.ErrorLog.
	errorLog *log.Logger
}

var ErrBotAlreadyExists = errors.New("bot already exists in bot mapping")

// addBot Adds a new bot to the botMapping structure.
func (m *botMapping) addBot(b *gotgbot.Bot, updateChan chan json.RawMessage, pollChan chan struct{}, urlPath string, webhookSecret string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}

	if _, ok := m.mapping[b.Token]; ok {
		return ErrBotAlreadyExists
	}

	m.mapping[b.Token] = botData{
		bot:           b,
		updateChan:    updateChan,
		polling:       pollChan,
		urlPath:       urlPath,
		webhookSecret: webhookSecret,
	}
	return nil
}

func (m *botMapping) removeBot(token string) (botData, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()

	bData, ok := m.mapping[token]
	if !ok {
		return botData{}, false
	}

	delete(m.mapping, token)
	return bData, true
}

func (m *botMapping) removeAllBots() []botData {
	m.mux.Lock()
	defer m.mux.Unlock()

	bots := make([]botData, 0, len(m.mapping))
	for key, bData := range m.mapping {
		bots = append(bots, bData)
		delete(m.mapping, key)
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

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (m *botMapping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	potentialToken := strings.TrimPrefix(r.URL.Path, "/")
	b, ok := m.getBot(potentialToken)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if secret := b.webhookSecret; secret != "" && secret != r.Header.Get("X-Telegram-Bot-Api-Secret-Token") {
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

func (m *botMapping) logf(format string, args ...interface{}) {
	if m.errorLog != nil {
		m.errorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
