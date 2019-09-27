package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func putln(s tcell.Screen, str string, row int, style tcell.Style) {
	puts(s, style, 1, row, str)
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

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

	printRows(s, SearchState{query: "", filteredLines: lines, selected: 0})

	handleEvents(lines, s)

	s.Fini()
}

func handleEvents(lines []string, s tcell.Screen) {
	eventChannel := newEventsChannel(s, "", lines)
	for ev := range eventChannel {
		switch ev.(type) {
		case QueryChangedEvent:
			qChangedEv := ev.(QueryChangedEvent)
			printRows(s, qChangedEv.state)
		case SelectedChangedEvent:
			sChangedEv := ev.(SelectedChangedEvent)
			printRows(s, sChangedEv.state)
		case ScreenResizeEvent:
			s.Sync()
		case EntryFinalSelectEvent:
			finalSelectEvt := ev.(EntryFinalSelectEvent)
			fmt.Printf(finalSelectEvt.entry)
			break
		}
	}
}

func printRows(s tcell.Screen, state SearchState) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen{}
	sc.width = w
	sc.height = h
	sc.appendRow(fmt.Sprintf("> %s", state.query), 0, bold)

	for i, l := range state.filteredLines {
		if i == state.selected {
			sc.appendRow(fmt.Sprintf(">  %s", l), 0, blink)
		} else {
			sc.appendRow(fmt.Sprintf("   %s", l), 0, plain)
		}
	}
	sc.printAll(s)

	s.Sync()
}

// template:
// >
//  {{#lines}}
//	{{if .selected}}{{style=highlight}}{{.}}{{^style}}{{else}}
//	{{fi}}
//  {{^lines}}

type ContentBlock struct {
	r     rune
	style tcell.Style
}

type Row struct {
	blocks []ContentBlock
	width  int
}

func (row *Row) writeRune(r rune, x int, style tcell.Style) {
	row.blocks[x].r = r
	row.blocks[x].style = style
}

func (row *Row) writeString(s string, x int, style tcell.Style) {
	i := 0
	for _, char := range s {
		if i+x >= row.width {
			break
		}
		row.writeRune(char, i+x, style)
		i++
	}
}

func newRow(width int) Row {
	return Row{
		width:  width,
		blocks: make([]ContentBlock, width),
	}
}

type screen struct {
	rows   []Row
	width  int
	height int
}

func (sc *screen) setRune(x int, y int, r rune, style tcell.Style) {
	sc.rows[y].blocks[x].r = r
	sc.rows[y].blocks[x].style = style
}

// appends a row at height = current_max_height + 1
func (sc *screen) appendRow(s string, x int, style tcell.Style) {
	r := newRow(sc.width)
	r.writeString(s, x, style)
	sc.rows = append(sc.rows, r)
}

func (sc *screen) printAll(s tcell.Screen) {
	for y, r := range sc.rows {
		for x, b := range r.blocks {
			s.SetContent(x, sc.height-(y+1), b.r, []rune{}, b.style)
		}
	}
}
