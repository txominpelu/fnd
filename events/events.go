package events

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/txominpelu/fnd/search"
)

func NewEventsChannel(s tcell.Screen, query string, searcher search.TextSearcher, sorter search.Compare) chan Event {
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
					break
				case tcell.KeyEnter:
					notifier.triggerSelect()
					break
				case tcell.KeyUp:
					if notifier.currentState.Selected+1 < len(notifier.currentState.FilteredLines(searcher, sorter)) {
						notifier.setSelected(notifier.currentState.Selected + 1)
					}
				case tcell.KeyDown:
					if notifier.currentState.Selected > 0 {
						notifier.setSelected(notifier.currentState.Selected - 1)
					}
				case tcell.KeyDEL:
					if len(notifier.currentState.Query) > 0 {
						notifier.setQuery(notifier.currentState.Query[:len(notifier.currentState.Query)-1], searcher, sorter)
					}
				case tcell.KeyBS:
					if len(notifier.currentState.Query) > 0 {
						notifier.setQuery(notifier.currentState.Query[:len(notifier.currentState.Query)-1], searcher, sorter)
					}
				case tcell.KeyRune:
					notifier.setQuery(fmt.Sprintf("%s%c", notifier.currentState.Query, ev.Rune()), searcher, sorter)
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

func (s *StateChangeNotifier) setQuery(query string, searcher search.TextSearcher, sorter search.Compare) {
	if s.currentState.Query != query {
		s.change(func(newState *SearchState) {
			(*newState).Query = query
		})
		filteredEntries := s.currentState.FilteredLines(searcher, sorter)
		if len(filteredEntries) <= s.currentState.Selected {
			s.setSelected(0)
		}
	}
}

func (s *StateChangeNotifier) triggerResize() {
	s.notifyChan <- ScreenResizeEvent{s.currentState}
}

func (s *StateChangeNotifier) triggerEscape() {
	s.notifyChan <- EscapeEvent{s.currentState}
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
		state: newState,
	}
	s.currentState = newState

}

//SearchState current state of the search
type SearchState struct {
	Query    string
	Selected int
}

func (state SearchState) FilteredLines(searcher search.TextSearcher, sorter search.Compare) []search.Document {
	return search.SortDocuments(
		searcher.FilterEntries(search.ParseQuery(state.Query)),
		searcher,
		sorter,
	)
}

func (state SearchState) Entry(searcher search.TextSearcher, sorter search.Compare) search.Document {
	filtered := state.FilteredLines(searcher, sorter)
	if state.Selected < len(filtered) {
		return filtered[state.Selected]
	} else {
		return search.Document{}
	}
}

//Event any event that can happen inside the application
type Event interface {
	State() SearchState
}

type ScreenResizeEvent struct {
	state SearchState
}

func (e ScreenResizeEvent) State() SearchState {
	return e.state
}

type EntryFinalSelectEvent struct {
	state SearchState
}

func (e EntryFinalSelectEvent) State() SearchState {
	return e.state
}

type EscapeEvent struct {
	state SearchState
}

func (e EscapeEvent) State() SearchState {
	return e.state
}

type SearchStateChanged struct {
	oldState SearchState
	state    SearchState
}

func (e SearchStateChanged) State() SearchState {
	return e.state
}
