package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

func newEventsChannel(s tcell.Screen, query string, entries []string) chan Event {
	out := make(chan Event)
	selected := 0

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(out)
					return
				case tcell.KeyEnter:
					close(out)
					return
				case tcell.KeyUp:
					filteredEntries := filterEntries(entries, query)
					if selected+1 < len(filteredEntries) {
						selected++
						out <- SelectedChangedEvent{
							oldSelected: selected - 1,
							state: SearchState{
								query:         query,
								filteredLines: filteredEntries,
								selected:      selected,
							},
						}
					}
				case tcell.KeyDown:
					filteredEntries := filterEntries(entries, query)
					if selected > 0 {
						selected--
						out <- SelectedChangedEvent{
							oldSelected: selected + 1,
							state: SearchState{
								query:         query,
								filteredLines: filteredEntries,
								selected:      selected,
							},
						}
					}
				case tcell.KeyDEL:
					if len(query) > 0 {
						filteredEntries := filterEntries(entries, query)
						oldQuery := query
						query = query[:len(query)-1]
						if len(filteredEntries) <= selected {
							selected = 0
						}
						out <- QueryChangedEvent{
							oldQuery: oldQuery,
							state: SearchState{
								query:         query,
								filteredLines: filteredEntries,
								selected:      selected,
							},
						}
					}
				case tcell.KeyRune:
					oldQuery := query
					query = fmt.Sprintf("%s%c", query, ev.Rune())
					filteredEntries := filterEntries(entries, query)
					if len(filteredEntries) <= selected {
						selected = 0
					}
					out <- QueryChangedEvent{
						oldQuery: oldQuery,
						state: SearchState{
							query:         query,
							filteredLines: filteredEntries,
							selected:      selected,
						},
					}
				}
			case *tcell.EventResize:

				s.Sync()
			}

		}

	}()
	return out
}

func filterEntries(lines []string, query string) []string {
	filteredLines := []string{}
	for _, l := range lines {
		if strings.Contains(l, query) {
			filteredLines = append(filteredLines, l)
		}
	}
	return filteredLines
}

//SearchState current state of the search
type SearchState struct {
	query         string
	filteredLines []string
	selected      int
}

//Event any event that can happen inside the application
type Event interface {
}

//ChangedEvent query has changed
type QueryChangedEvent struct {
	oldQuery string
	state    SearchState
}

type ScreenResizeEvent struct {
}

//SelectedChangedEvent selected item has changed
type SelectedChangedEvent struct {
	oldSelected int
	state       SearchState
}
