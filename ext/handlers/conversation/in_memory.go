package conversation

import (
	"errors"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var ErrKeyNotFound = errors.New("conversation key not found")

// InMemoryStorage is a thread-safe in-memory implementation of the Storage interface.
type InMemoryStorage struct {
	// keyStrategy defines how to calculate keys for each conversation.
	keyStrategy KeyStrategy
	// conversations is a map of key -> state, which tracks at which point of each conversation a user/chat is.
	conversations map[string]State
	// lock allows us to ensure synchronous data access.
	lock sync.RWMutex
}

func NewInMemoryStorage(strategy KeyStrategy) *InMemoryStorage {
	return &InMemoryStorage{
		keyStrategy:   strategy,
		lock:          sync.RWMutex{},
		conversations: map[string]State{},
	}
}

func (c *InMemoryStorage) Get(ctx *ext.Context) (*State, error) {
	key, err := StateKey(ctx, c.keyStrategy)
	if err != nil {
		return nil, err
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.conversations == nil {
		return nil, ErrKeyNotFound
	}

	s, ok := c.conversations[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return &s, nil
}

func (c *InMemoryStorage) Set(ctx *ext.Context, state State) error {
	key, err := StateKey(ctx, c.keyStrategy)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conversations == nil {
		c.conversations = map[string]State{}
	}

	c.conversations[key] = state
	return nil
}

func (c *InMemoryStorage) Delete(ctx *ext.Context) error {
	key, err := StateKey(ctx, c.keyStrategy)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conversations == nil {
		return nil
	}

	delete(c.conversations, key)
	return nil
}
