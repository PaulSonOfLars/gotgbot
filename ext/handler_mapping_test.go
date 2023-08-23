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
	checkList(t, "currGroups", currGroups[0], firstHandler)

	// Add a second handler.
	m.add(secondHandler, 0)
	newGroups := m.getGroups()
	checkList(t, "newgroups;currGroups", currGroups[0], firstHandler)
	checkList(t, "newgroups;newGroups", newGroups[0], firstHandler, secondHandler)

	// Remove second handler..
	ok := m.remove(secondHandler.Name(), 0)
	if !ok {
		t.Errorf("failed to remove second handler")
	}
	delGroups := m.getGroups()
	checkList(t, "delgroups;currGroups", currGroups[0], firstHandler)
	checkList(t, "delgroups;newGroups", newGroups[0], firstHandler, secondHandler)
	checkList(t, "delgroups;delGroups", delGroups[0], firstHandler)

	// Re-add second handler.
	m.add(secondHandler, 0)
	reAddedGroups := m.getGroups()
	checkList(t, "readded;currGroups", currGroups[0], firstHandler)
	checkList(t, "readded;newGroups", newGroups[0], firstHandler, secondHandler)
	checkList(t, "readded;delGroups", delGroups[0], firstHandler)
	checkList(t, "readded;reAddedGroups", reAddedGroups[0], firstHandler, secondHandler)

	// Remove first handler.
	ok = m.remove(firstHandler.Name(), 0)
	if !ok {
		t.Errorf("failed to remove second handler")
	}
	noFirstGroups := m.getGroups()
	checkList(t, "nofirst;currGroups", currGroups[0], firstHandler)
	checkList(t, "nofirst;newGroups", newGroups[0], firstHandler, secondHandler)
	checkList(t, "nofirst;delGroups", delGroups[0], firstHandler)
	checkList(t, "nofirst;reAddedGroups", reAddedGroups[0], firstHandler, secondHandler)
	checkList(t, "nofirst;noFirstGroups", noFirstGroups[0], secondHandler)
}

func checkList(t *testing.T, name string, got []Handler, expected ...Handler) {
	if len(got) != len(expected) {
		t.Errorf("mismatch on length of expected outputs for %s - got %d, expected %d", name, len(got), len(expected))
	}
	for idx, v := range got {
		if v.Name() != expected[idx].Name() {
			t.Errorf("unexpected output name for %s - IDX %d got %s, expected %s", name, idx, v.Name(), expected[idx].Name())
		}
	}
}
