package events

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

func filterEntries(lines []string, query string) []string {
	fLines := []string{}
	for _, l := range lines {
		if strings.Contains(l, query) {
			fLines = append(fLines, l)
		}
	}
	return fLines
}

func NewEventsChannel(s tcell.Screen, query string, entries []string) chan Event {
	out := make(chan Event)
	st := SearchState{query, 0}
	notifier := StateChangeNotifier{currentState: st, notifyChan: out}

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					notifier.triggerEscape()
					close(out)
				case tcell.KeyEnter:
					notifier.triggerSelect()
					close(out)
				case tcell.KeyUp:
					if notifier.currentState.Selected+1 < len(notifier.currentState.FilteredLines(entries)) {
						notifier.setSelected(notifier.currentState.Selected + 1)
					}
				case tcell.KeyDown:
					if notifier.currentState.Selected > 0 {
						notifier.setSelected(notifier.currentState.Selected - 1)
					}
				case tcell.KeyDEL:
					if len(notifier.currentState.Query) > 0 {
						notifier.setQuery(notifier.currentState.Query[:len(notifier.currentState.Query)-1], entries)
					}
				case tcell.KeyRune:
					notifier.setQuery(fmt.Sprintf("%s%c", notifier.currentState.Query, ev.Rune()), entries)
				}
			case *tcell.EventResize:
				notifier.triggerResize()
			}
		}

	}()
	return out
}

type StateChangeNotifier struct {
	notifyChan   chan Event
	currentState SearchState
}

func (s *StateChangeNotifier) setSelected(selected int) {
	if s.currentState.Selected != selected {
		s.change(func(newState *SearchState) {
			(*newState).Selected = selected
			(*newState).Query = s.currentState.Query
		})
	}
}

func (s *StateChangeNotifier) setQuery(query string, entries []string) {
	if s.currentState.Query != query {
		s.change(func(newState *SearchState) {
			(*newState).Query = query
		})
		filteredEntries := s.currentState.FilteredLines(entries)
		if len(filteredEntries) <= s.currentState.Selected {
			s.setSelected(0)
		}
	}
}

func (s *StateChangeNotifier) triggerResize() {
	s.notifyChan <- ScreenResizeEvent{s.currentState}
}

func (s *StateChangeNotifier) triggerEscape() {
	s.notifyChan <- EscapeEvent{}
}

func (s *StateChangeNotifier) triggerSelect() {
	s.notifyChan <- EntryFinalSelectEvent{s.currentState}
}

func (s *StateChangeNotifier) change(updateState func(*SearchState)) {
	newState := SearchState{
		Query:    s.currentState.Query,
		Selected: s.currentState.Selected,
	}
	updateState(&newState)
	s.notifyChan <- SearchStateChanged{
		oldState: SearchState{
			Query:    s.currentState.Query,
			Selected: s.currentState.Selected,
		},
		State: newState,
	}
	s.currentState = newState

}

//SearchState current state of the search
type SearchState struct {
	Query    string
	Selected int
}

func (state SearchState) FilteredLines(entries []string) []string {
	return filterEntries(entries, state.Query)
}

func (state SearchState) Entry(entries []string) string {
	return filterEntries(entries, state.Query)[state.Selected]
}

//Event any event that can happen inside the application
type Event interface {
}

type ScreenResizeEvent struct {
	State SearchState
}

type EntryFinalSelectEvent struct {
	State SearchState
}

type EscapeEvent struct {
}

type SearchStateChanged struct {
	oldState SearchState
	State    SearchState
}
