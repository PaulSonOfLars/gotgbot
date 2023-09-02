package ext

import (
	"encoding/json"
	"errors"
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
}

// botMapping Ensures that all botData is stored in a thread-safe manner.
type botMapping struct {
	// mapping keeps track of the data required for each bot. The key is the bot token.
	mapping map[string]botData
	// mux attempts to keep the botMapping data concurrency-safe.
	mux sync.RWMutex
}

var ErrBotAlreadyExists = errors.New("bot already exists in bot mapping")

// addBot Adds a new bot to the botMapping structure.
func (m *botMapping) addBot(b *gotgbot.Bot, updateChan chan json.RawMessage, pollChan chan struct{}, urlPath string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}

	if _, ok := m.mapping[b.Token]; ok {
		return ErrBotAlreadyExists
	}

	m.mapping[b.Token] = botData{
		bot:        b,
		updateChan: updateChan,
		polling:    pollChan,
		urlPath:    urlPath,
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
