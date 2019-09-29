package events

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func TestQueryAndChangeSelect(t *testing.T) {
	s, _ := tcell.NewScreen()
	defer s.Fini()
	encoding.Register()
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}
	eventsChan := NewEventsChannel(s, "", []string{"hello", "hellos", "hellod", "helloll"})
	go func() {
		s.PostEvent(tcell.NewEventKey(tcell.KeyRune, 'h', tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyESC, 0, tcell.ModNone))
	}()
	events := []Event{}
	for ev := range eventsChan {
		events = append(events, ev)
	}
	expected := SearchState{
		Query:    "hel",
		Selected: 3,
	}
	ev := events[len(events)-2].(SearchStateChanged)
	if !reflect.DeepEqual(ev.State, expected) {
		t.Errorf("expected: '%v' got: '%v'\n", expected, ev.State)
	}

}

func TestSelectGoesZero(t *testing.T) {
	s, _ := tcell.NewScreen()
	defer s.Fini()
	encoding.Register()
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}
	eventsChan := NewEventsChannel(s, "", []string{"hello", "hellos", "hellod", "helloll"})
	go func() {
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyRune, 'h', tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone))
		s.PostEvent(tcell.NewEventKey(tcell.KeyESC, 0, tcell.ModNone))
	}()
	events := []Event{}
	for ev := range eventsChan {
		events = append(events, ev)
	}
	expected := SearchState{
		Query:    "ho",
		Selected: 0,
	}
	ev := events[len(events)-2].(SearchStateChanged)
	if !reflect.DeepEqual(ev.State, expected) {
		t.Errorf("expected: '%v' got: '%v'\n", expected, ev.State)
	}

}
