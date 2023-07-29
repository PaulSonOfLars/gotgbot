package ext

import (
	"encoding/json"
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

// addBot Adds a new bot to the botMapping structure.
func (m *botMapping) addBot(b *gotgbot.Bot, updateChan chan json.RawMessage, pollChan chan struct{}, urlPath string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}

	m.mapping[b.Token] = botData{
		bot:        b,
		updateChan: updateChan,
		polling:    pollChan,
		urlPath:    urlPath,
	}
}

// stopAllBots Stops all bots.
func (m *botMapping) stopAllBots() {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Close all the update channels and polling loops
	for key, bData := range m.mapping {
		stopBot(bData)
		delete(m.mapping, key)
	}
}

func (m *botMapping) removeBot(b *gotgbot.Bot) (botData, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()

	bData, ok := m.mapping[b.Token]
	if !ok {
		return botData{}, false
	}

	delete(m.mapping, b.Token)
	return bData, true
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

func (m *botMapping) getBot(b *gotgbot.Bot) (botData, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	bData, ok := m.mapping[b.Token]
	return bData, ok
}
