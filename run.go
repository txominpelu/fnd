package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/txominpelu/rjobs/screen"
)

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	printRows(s, SearchState{query: "", selected: 0}, lines)

	handleEvents(lines, s)

	s.Fini()
}

func handleEvents(lines []string, s tcell.Screen) {
	eventChannel := newEventsChannel(s, "", lines)
	for ev := range eventChannel {
		switch ev.(type) {
		case SearchStateChanged:
			qChangedEv := ev.(SearchStateChanged)
			printRows(s, qChangedEv.state, lines)
		case ScreenResizeEvent:
			s.Sync()
		case EntryFinalSelectEvent:
			finalSelectEvt := ev.(EntryFinalSelectEvent)
			fmt.Printf(finalSelectEvt.state.entry(lines))
			break
		}
	}
}

// template:
// >
//  {{#lines}}
//	{{if .selected}}{{style=highlight}}{{.}}{{^style}}{{else}}
//	{{fi}}
//  {{^lines}}

func printRows(s tcell.Screen, state SearchState, entries []string) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.query), 0, bold)

	for i, l := range state.filteredLines(entries) {
		if i == state.selected {
			sc.AppendRow(fmt.Sprintf(">  %s", l), 0, blink)
		} else {
			sc.AppendRow(fmt.Sprintf("   %s", l), 0, plain)
		}
	}
	sc.PrintAll(s)

	s.Sync()
}
