package ext

import (
	"fmt"
	"sync"
	"testing"
)

// This test should demonstrate that once obtained, a list will not be changed by any additions/removals to that list by another call.
func Test_handlerMappings_getGroupsConcurrentSafe(t *testing.T) {
	m := handlerMapping{}
	firstHandler := DummyHandler{N: "first"}
	secondHandler := DummyHandler{N: "second"}

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
	t.Helper()

	if len(got) != len(expected) {
		t.Errorf("mismatch on length of expected outputs for %s - got %d, expected %d", name, len(got), len(expected))
	}

	for idx, v := range got {
		//nolint:gosec
		if v.Name() != expected[idx].Name() {
			//nolint:gosec
			t.Errorf("unexpected output name for %s - IDX %d got %s, expected %s", name, idx, v.Name(), expected[idx].Name())
		}
	}
}

func Test_handlerMappings_remove(t *testing.T) {
	m := &handlerMapping{}
	handler := DummyHandler{N: "test"}

	t.Run("nonExistent", func(t *testing.T) {
		// removing an item that doesnt exist returns "false"
		if got := m.remove(handler.Name(), 0); got {
			t.Errorf("remove() = %v, want false", got)
		}
	})

	t.Run("removalSuccess", func(t *testing.T) {
		m.add(handler, 0)
		// removing an item that DOES exist, returns true
		if got := m.remove(handler.Name(), 0); !got {
			t.Errorf("remove() = %v, want true", got)
		}
		// And so the second time, it returns false
		if got := m.remove(handler.Name(), 0); got {
			t.Errorf("remove() = %v, want false", got)
		}
	})

	t.Run("removalSuccess", func(t *testing.T) {
		m.add(handler, 0)
		// removing an item that DOES exist, returns true
		if got := m.remove(handler.Name(), 0); !got {
			t.Errorf("remove() = %v, want true", got)
		}
		// And so the second time, it returns false
		if got := m.remove(handler.Name(), 0); got {
			t.Errorf("remove() = %v, want false", got)
		}
	})

	t.Run("removalDifferentIndexes", func(t *testing.T) {
		m.add(handler, 1)
		m.add(handler, 2)
		if got := m.remove(handler.Name(), 2); !got {
			t.Errorf("remove() = %v, want true", got)
		}
	})
}

// TestHandlerMapping_Add_DoesNotMutatePreviousSlice guards against add()
// being used to do `append(currHandlers, h)` directly. When currHandlers
// had spare capacity, that write landed in the same backing array referenced
// by any slice a caller had already obtained from getGroups(), silently
// corrupting it. The fix always copies into a freshly allocated slice.
func TestHandlerMapping_Add_DoesNotMutatePreviousSlice(t *testing.T) {
	m := &handlerMapping{}

	m.add(DummyHandler{N: "first"}, 0)

	// Snapshot the slice as observed *before* the second add().
	before := m.getGroups()[0]

	m.add(DummyHandler{N: "second"}, 0)

	// The previously observed slice must be untouched by the later add().
	if len(before) != 1 || before[0].Name() != "dummy_first" {
		t.Fatalf("previously observed slice was mutated: got %+v", before)
	}

	after := m.getGroups()[0]
	if len(after) != 2 || after[0].Name() != "dummy_first" || after[1].Name() != "dummy_second" {
		t.Fatalf("unexpected handlers after add: got %+v", after)
	}
}

// TestHandlerMapping_ConcurrentAddAndRead exercises add() and getGroups()
// concurrently.
func TestHandlerMapping_ConcurrentAddAndRead(t *testing.T) {
	m := &handlerMapping{}
	const group = 0
	const n = 500

	var wg sync.WaitGroup
	wg.Add(2)

	// Writer: keeps appending handlers to the same group.
	go func() {
		defer wg.Done()
		for i := range n {
			m.add(DummyHandler{N: fmt.Sprintf("h%d", i)}, group)
		}
	}()

	// Reader: repeatedly reads and iterates over the current handlers.
	go func() {
		defer wg.Done()
		for range n {
			for _, handlers := range m.getGroups() {
				for _, h := range handlers {
					_ = h.Name()
				}
			}
		}
	}()

	wg.Wait()

	final := m.getGroups()
	if len(final) != 1 || len(final[0]) != n {
		t.Fatalf("expected %d handlers in group %d, got %d", n, group, len(final[0]))
	}
}

// TestHandlerMapping_ConcurrentAddRemoveAndRead exercises add(), remove(), and
// getGroups() all at once under `go test -race`. remove() already copied its
// slice before mutating it, so this mainly guards against a regression in that
// pattern and against any new interaction between add() and remove() sharing
// backing arrays.
//
// To keep the expected final count deterministic despite the concurrency, the
// writer signals each handler's name over a channel right after adding it, and
// the remover only ever removes handlers it knows have already been added.
func TestHandlerMapping_ConcurrentAddRemoveAndRead(t *testing.T) {
	m := &handlerMapping{}
	const group = 0
	const n = 500

	added := make(chan string, n)

	var wg sync.WaitGroup
	wg.Add(3)

	// Writer: adds n handlers to the group, signalling each name once added.
	go func() {
		defer wg.Done()
		defer close(added)
		for i := range n {
			name := fmt.Sprintf("h%d", i)
			m.add(DummyHandler{N: name}, group)
			added <- name
		}
	}()

	// Remover: removes every other handler as soon as it's confirmed added.
	removed := 0
	go func() {
		defer wg.Done()
		i := 0
		for name := range added {
			if i%2 == 0 {
				if m.remove(name, group) {
					removed++
				}
			}
			i++
		}
	}()

	// Reader: repeatedly reads and iterates over the current handlers while
	// add() and remove() are both mutating the map concurrently.
	go func() {
		defer wg.Done()
		for range n {
			for _, handlers := range m.getGroups() {
				for _, h := range handlers {
					_ = h.Name()
				}
			}
		}
	}()

	wg.Wait()

	wantCount := n - removed
	final := m.getGroups()

	var finalCount int
	if wantCount > 0 {
		if len(final) != 1 {
			t.Fatalf("expected group %d to still be present, got groups: %+v", group, final)
		}
		finalCount = len(final[0])
	} else if len(final) != 0 {
		t.Fatalf("expected group %d to have been removed entirely, got groups: %+v", group, final)
	}

	if finalCount != wantCount {
		t.Fatalf("expected %d handlers remaining, got %d (removed %d of %d)", wantCount, finalCount, removed, n)
	}
}
