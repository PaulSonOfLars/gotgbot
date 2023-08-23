package ext

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type dummy struct {
	f    func(bot *gotgbot.Bot, ctx *Context) error
	name string
}

func (d dummy) CheckUpdate(b *gotgbot.Bot, ctx *Context) bool {
	return true
}

func (d dummy) HandleUpdate(b *gotgbot.Bot, ctx *Context) error {
	return d.f(b, ctx)
}

func (d dummy) Name() string {
	return "dummy" + d.name
}

// This test should demonstrate that once obtained, a list will not be changed by any additions/removals to that list by another call.
func Test_handlerMappings_getGroupsConcurrentSafe(t *testing.T) {
	m := handlerMappings{}
	firstHandler := dummy{name: "first"}
	secondHandler := dummy{name: "second"}

	// We expect 0 groups at the start
	startGroups := m.getGroups()
	if len(startGroups) != 0 {
		t.Errorf("failed predicate group layout")
	}

	// Add one handler.
	m.add(firstHandler, 0)
	currGroups := m.getGroups()
	if len(currGroups) != 1 && len(startGroups) != 0 {
		t.Errorf("Start groups should be 0, curr groups should be 1; got %d and %d", len(startGroups), len(currGroups))
	}
	if len(currGroups[0]) != 1 {
		t.Errorf("length of group 0 in the groups list should be 1; got %d", len(currGroups[0]))
	}

	// Add a second handler.
	m.add(secondHandler, 0)
	newGroups := m.getGroups()
	if len(currGroups[0]) != 1 {
		t.Errorf("length of group 0 in currGroups should be 1; got %d", len(currGroups[0]))
	}
	if len(newGroups[0]) != 2 {
		t.Errorf("length of group 0 in the newGroups 1; got %d", len(newGroups[0]))
	}

	// Remove second handler..
	ok := m.remove(secondHandler.Name(), 0)
	if !ok {
		t.Errorf("failed to remove second handler")
	}
	delGroups := m.getGroups()
	if len(currGroups[0]) != 1 {
		t.Errorf("length of group 0 in currGroups should be 1; got %d", len(currGroups[0]))
	}
	if len(newGroups[0]) != 2 {
		t.Errorf("length of group 0 in the newGroups 1; got %d", len(newGroups[0]))
	}
	if len(delGroups[0]) != 1 {
		t.Errorf("length of group 0 in delGroups should be 1; got %d", len(delGroups[0]))
	}

	// Re-add second handler.
	m.add(secondHandler, 0)
	reAddedGroups := m.getGroups()
	if len(currGroups[0]) != 1 &&
		currGroups[0][0].Name() == firstHandler.Name() {
		t.Errorf("length of group 0 in currGroups should be 1; got %d", len(currGroups[0]))
	}
	if len(newGroups[0]) != 2 &&
		newGroups[0][1].Name() == secondHandler.Name() {
		t.Errorf("length of group 0 in the newGroups 1; got %d", len(newGroups[0]))
	}
	if len(delGroups[0]) != 1 &&
		delGroups[0][0].Name() == firstHandler.Name() {
		t.Errorf("length of group 0 in delGroups should be 1; got %d", len(delGroups[0]))
	}
	if len(reAddedGroups[0]) != 2 &&
		newGroups[0][0].Name() == firstHandler.Name() &&
		newGroups[0][1].Name() == secondHandler.Name() {
		t.Errorf("length of group 0 in reAddedGroups should be 2; got %d", len(reAddedGroups[0]))
	}

	// Remove first handler.
	ok = m.remove(firstHandler.Name(), 0)
	if !ok {
		t.Errorf("failed to remove second handler")
	}
	noFirstGroups := m.getGroups()
	if len(currGroups[0]) != 1 &&
		currGroups[0][0].Name() == firstHandler.Name() {
		t.Errorf("length of group 0 in currGroups should be 1; got %d", len(currGroups[0]))
	}
	if len(newGroups[0]) != 2 &&
		newGroups[0][1].Name() == secondHandler.Name() {
		t.Errorf("length of group 0 in the newGroups 1; got %d", len(newGroups[0]))
	}
	if len(delGroups[0]) != 1 &&
		delGroups[0][0].Name() == firstHandler.Name() {
		t.Errorf("length of group 0 in delGroups should be 1; got %d", len(delGroups[0]))
	}
	if len(reAddedGroups[0]) != 2 &&
		reAddedGroups[0][0].Name() == firstHandler.Name() &&
		reAddedGroups[0][1].Name() == secondHandler.Name() {
		t.Errorf("length of group 0 in reAddedGroups should be 2; got %d", len(reAddedGroups[0]))
	}
	if len(noFirstGroups[0]) != 1 &&
		noFirstGroups[0][0].Name() == secondHandler.Name() {
		t.Errorf("length of group 0 in noFirstGroups should be 2; got %d", len(noFirstGroups[0]))
	}
}
