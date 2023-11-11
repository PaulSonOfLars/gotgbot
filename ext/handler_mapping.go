package ext

import (
	"sort"
	"sync"
)

type handlerMapping struct {
	// mutex is used to ensure everything threadsafe.
	mutex sync.RWMutex

	// handlerGroups represents the list of available handler groups, numerically sorted.
	handlerGroups []int
	// handlers represents all available handlers, split into groups (see handlerGroups).
	handlers map[int][]Handler
}

func (m *handlerMapping) add(h Handler, group int) {
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

func (m *handlerMapping) remove(name string, group int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currHandlers, ok := m.handlers[group]
	if !ok {
		// group does not exist; removal failed.
		return false
	}

	for idx, handler := range currHandlers {
		if handler.Name() != name {
			continue
		}

		// Only one item left, so just delete the group entirely.
		if len(currHandlers) == 1 {
			// get index of the current group to remove it from the list of handlergroups
			gIdx := getIndex(group, m.handlerGroups)
			if gIdx != -1 {
				m.handlerGroups = append(m.handlerGroups[:gIdx], m.handlerGroups[gIdx+1:]...)
			}
			delete(m.handlers, group)
			return true
		}

		// Make sure to copy the handler list to ensure we don't change the values of the underlying arrays, which
		// could cause slice access issues when used concurrently.
		newHandlers := make([]Handler, len(m.handlers[group]))
		copy(newHandlers, m.handlers[group])

		m.handlers[group] = append(newHandlers[:idx], newHandlers[idx+1:]...)
		return true
	}
	// handler not found - removal failed.
	return false
}

func getIndex(find int, is []int) int {
	for i, v := range is {
		if v == find {
			return i
		}
	}
	return -1
}

func (m *handlerMapping) removeGroup(group int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.handlers[group]; !ok {
		// Group doesn't exist in map, so already removed.
		return false
	}

	for idx, handlerGroup := range m.handlerGroups {
		if handlerGroup != group {
			continue
		}

		m.handlerGroups = append(m.handlerGroups[:idx], m.handlerGroups[idx+1:]...)
		delete(m.handlers, group)
		// Group found, and deleted. Success!
		return true
	}
	// Group not found in list - so already removed.
	return false
}

func (m *handlerMapping) getGroups() [][]Handler {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	allHandlers := make([][]Handler, len(m.handlerGroups))
	for idx, num := range m.handlerGroups {
		allHandlers[idx] = m.handlers[num]
	}
	return allHandlers
}
