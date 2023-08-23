package ext

import (
	"sort"
	"sync"
)

type handlerMappings struct {
	// mutex is used to ensure everything threadsafe.
	mutex sync.RWMutex

	// handlerGroups represents the list of available handler groups, numerically sorted.
	handlerGroups []int
	// handlers represents all available handlers, split into groups (see handlerGroups).
	handlers map[int][]Handler
}

func (m *handlerMappings) add(h Handler, group int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.handlers == nil {
		m.handlers = map[int][]Handler{}
	}
	currHandlers, ok := m.handlers[group]
	if !ok {
		m.handlerGroups = append(m.handlerGroups, group)
		sort.Ints(m.handlerGroups)
	}
	m.handlers[group] = append(currHandlers, h)
}

func (m *handlerMappings) remove(name string, group int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currHandlers, ok := m.handlers[group]
	if !ok {
		// group does not exist; removal failed.
		return false
	}

	for i, handler := range currHandlers {
		if handler.Name() != name {
			continue
		}

		// Only one item left, so just delete the group entirely.
		if len(currHandlers) == 1 {
			m.handlerGroups = append(m.handlerGroups[:group], m.handlerGroups[group+1:]...)
			delete(m.handlers, group)
			return true
		}

		m.handlers[group] = append(currHandlers[:i], currHandlers[i+1:]...)
		return true
	}
	// handler not found - removal failed.
	return false
}

func (m *handlerMappings) removeGroup(group int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.handlers[group]; !ok {
		// Group doesn't exist in map, so already removed.
		return false
	}

	for j, handlerGroup := range m.handlerGroups {
		if handlerGroup != group {
			continue
		}

		m.handlerGroups = append(m.handlerGroups[:j], m.handlerGroups[j+1:]...)
		delete(m.handlers, group)
		// Group found, and deleted. Success!
		return true
	}
	// Group not found in list - so already removed.
	return false
}

func (m *handlerMappings) getGroups() [][]Handler {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var handlers [][]Handler
	for _, num := range m.handlerGroups {
		handlers = append(handlers, m.handlers[num])
	}
	return handlers
}
