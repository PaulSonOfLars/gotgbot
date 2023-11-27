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
func (m *botMapping) addBot(bData botData) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}
	if m.urlMapping == nil {
		m.urlMapping = make(map[string]string)
	}

	if _, ok := m.mapping[bData.bot.Token]; ok {
		return ErrBotAlreadyExists
	}
	if _, ok := m.urlMapping[bData.urlPath]; bData.urlPath != "" && ok {
		return ErrBotUrlPathAlreadyExists
	}

	m.mapping[bData.bot.Token] = bData
	m.urlMapping[bData.urlPath] = bData.bot.Token
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
