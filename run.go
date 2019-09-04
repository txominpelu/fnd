package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	query := ""
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
	updateResults(s, lines, query)
	quit := make(chan struct{})

	handleEvents(lines, s, quit)

	<-quit

	s.Fini()
}

func handleEvents(lines []string, s tcell.Screen, quit chan struct{}) {
	query := ""
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyDEL:
					if len(query) > 0 {
						query = query[:len(query)-1]
						updateResults(s, lines, query)
					}
				case tcell.KeyRune:
					query = fmt.Sprintf("%s%c", query, ev.Rune())
					updateResults(s, lines, query)
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()
}

func updateResults(s tcell.Screen, lines []string, query string) {
	s.Clear()
	filteredLines := []string{}
	for _, l := range lines {
		if strings.Contains(l, query) {
			filteredLines = append(filteredLines, l)
		}
	}
	printRows(s, filteredLines, query)
}

func printRows(s tcell.Screen, filteredLines []string, query string) {
	_, h := s.Size()
	row := h
	plain := tcell.StyleDefault
	bold := tcell.StyleDefault.Bold(true)

	row--
	putln(s, fmt.Sprintf("Query: %s", query), row, bold)

	for _, l := range filteredLines {
		row--
		putln(s, l, row, plain)
	}

	s.Sync()
}

// stateMachine {
// 	delete -> {
// 		query = query[0:len(query)-1]
// 		updateResults(query)
// 	}
// 	letter -> {
// 		query = query + newLetter
// 		updateResults(query)
// 	}
// 	esc -> exit()
// }
//
// updateResults {
// 	filter(query)
// 	refresh(screen)
// }
