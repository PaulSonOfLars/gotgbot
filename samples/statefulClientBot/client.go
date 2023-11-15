package main

import (
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// A basic handler client to share state across executions.
// Note: This is a very simple layout which uses a shared mutex.
// It is all in-memory, and so will not persist data across restarts.
type client struct {
	// Use a mutex to avoid concurrency issues.
	// If you use multiple maps, you may want to use a new mutex for each one.
	rwMux sync.RWMutex

	// We use a double map to:
	// - map once for the user id
	// - map a second time for the keys a user can have
	// The second map has values of type "any" so anything can be stored in them, for the purpose of this example.
	// This could be improved by using a struct with typed fields, though this would need some additional handling to
	// ensure concurrent safety.
	userData map[int64]map[string]any

	// This struct could also contain:
	// - pointers to database connections
	// - pointers cache connections
	// - localised strings
	// - helper methods for retrieving/caching chat settings
}

func (c *client) getUserData(ctx *ext.Context, key string) (any, bool) {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()

	if c.userData == nil {
		return nil, false
	}

	userData, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		return nil, false
	}

	v, ok := userData[key]
	return v, ok
}

func (c *client) setUserData(ctx *ext.Context, key string, val any) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userData == nil {
		c.userData = map[int64]map[string]any{}
	}

	_, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		c.userData[ctx.EffectiveUser.Id] = map[string]any{}
	}
	c.userData[ctx.EffectiveUser.Id][key] = val
}
