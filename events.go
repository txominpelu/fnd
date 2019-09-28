package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/jinzhu/copier"
)

func newEventsChannel(s tcell.Screen, query string, entries []string) chan Event {
	out := make(chan Event)
	notifier := StateChangeNotifier{currentState: SearchState{query, 0}, notifyChan: out}

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
					if notifier.currentState.selected+1 < len(notifier.currentState.filteredLines(entries)) {
						notifier.setSelected(notifier.currentState.selected + 1)
					}
				case tcell.KeyDown:
					if notifier.currentState.selected > 0 {
						notifier.setSelected(notifier.currentState.selected - 1)
					}
				case tcell.KeyDEL:
					if len(notifier.currentState.query) > 0 {
						notifier.setQuery(notifier.currentState.query[:len(notifier.currentState.query)-1])
						filteredEntries := notifier.currentState.filteredLines(entries)
						if len(filteredEntries) <= notifier.currentState.selected {
							notifier.setSelected(0)
						}
					}
				case tcell.KeyRune:
					notifier.setQuery(fmt.Sprintf("%s%c", notifier.currentState.query, ev.Rune()))
				}
			case *tcell.EventResize:
				notifier.triggerResize()
			}
		}

	}()
	return out
}

func filterEntries(lines []string, query string) []string {
	fLines := []string{}
	for _, l := range lines {
		if strings.Contains(l, query) {
			fLines = append(fLines, l)
		}
	}
	return fLines
}

type StateChangeNotifier struct {
	notifyChan   chan Event
	currentState SearchState
}

func (s *StateChangeNotifier) setSelected(selected int) {
	if s.currentState.selected != selected {
		s.change(s.currentState, func(newState *SearchState) { (*newState).selected = selected })
	}
}

func (s *StateChangeNotifier) setQuery(query string) {
	if s.currentState.query != query {
		s.change(s.currentState, func(newState *SearchState) { (*newState).query = query })
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

func (s *StateChangeNotifier) change(currentState SearchState, updateState func(*SearchState)) {
	newState := SearchState{}
	copier.Copy(&s.currentState, &newState)
	updateState(&newState)
	s.notifyChan <- SearchStateChanged{
		oldState: s.currentState,
		state:    newState,
	}
	s.currentState = newState

}

//SearchState current state of the search
type SearchState struct {
	query    string
	selected int
}

func (state SearchState) filteredLines(entries []string) []string {
	return filterEntries(entries, state.query)
}

func (state SearchState) entry(entries []string) string {
	return filterEntries(entries, state.query)[state.selected]
}

//Event any event that can happen inside the application
type Event interface {
}

type ScreenResizeEvent struct {
	state SearchState
}

type EntryFinalSelectEvent struct {
	state SearchState
}

type EscapeEvent struct {
}

type SearchStateChanged struct {
	oldState SearchState
	state    SearchState
}
