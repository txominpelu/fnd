package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

func newEventsChannel(s tcell.Screen, query string) chan Event {
	out := make(chan Event)

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
				case tcell.KeyDEL:
					if len(query) > 0 {
						oldQuery := query
						query = query[:len(query)-1]
						out <- QueryChangedEvent{
							newQuery: query,
							oldQuery: oldQuery,
						}
					}
				case tcell.KeyRune:
					oldQuery := query
					query = fmt.Sprintf("%s%c", query, ev.Rune())
					out <- QueryChangedEvent{
						newQuery: query,
						oldQuery: oldQuery,
					}
				}
			case *tcell.EventResize:

				s.Sync()
			}

		}

	}()
	return out
}

type Event interface {
}

type QueryChangedEvent struct {
	oldQuery string
	newQuery string
}

type ScreenResizeEvent struct {
}
