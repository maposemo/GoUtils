package dispatcher

import (
	"testing"
)

type (
	cleanDish struct{ what string }
	cleanRoom struct{ what string }
)

type testListener struct {
	workCh chan string
}

func (u testListener) Listen(event interface{}) {
	switch evt := event.(type) {
	case cleanDish:
		u.workCh <- "Done " + evt.what
	default:
		u.workCh <- "Not My Work"
	}
}

func TestDispatcher(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatchCases := []struct {
		to   string
		what any
		exp  string
	}{
		{to: "me", what: cleanDish{what: "Fripan"}, exp: "Done Fripan"},
		{to: "wife", what: cleanDish{what: "Spoon"}, exp: "Done Spoon"},
		{to: "me", what: cleanRoom{what: "cleanRoom"}, exp: "Not My Work"},
		{to: "friend", what: cleanRoom{what: "cleanRoom"}, exp: "'friend' address is not registered"},
		{to: "noWhere", what: cleanDish{what: "nothing"}, exp: "'noWhere' address is not registered"},
	}

	// Test Register
	tester := testListener{workCh: make(chan string)}
	if err := dispatcher.Register(tester, "me", "wife"); err != nil {
		t.Errorf("Register failed(%v)", err)
	}

	// Test Dispatch
	for _, dc := range dispatchCases {
		if err := dispatcher.Dispatch(dc.to, dc.what); err != nil {
			if err.Error() != dc.exp {
				t.Errorf("Handle failed(%v)", err)
			}
			continue
		}

		select {
		case work := <-tester.workCh:
			if work != dc.exp {
				t.Errorf("Handle failed(%v)", work)
			}
		}
	}
}
