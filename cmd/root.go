package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/rjobs/events"
	"github.com/txominpelu/rjobs/screen"
)

var RootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: runRoot,
}

func init() {
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
}

func runRoot(cmd *cobra.Command, args []string) {

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

	printRows(s, events.SearchState{Query: "", Selected: 0}, lines)

	handleEvents(lines, s)

	s.Fini()
}

func handleEvents(lines []string, s tcell.Screen) {
	eventChannel := events.NewEventsChannel(s, "", lines)
	for ev := range eventChannel {
		switch ev.(type) {
		case events.SearchStateChanged:
			qChangedEv := ev.(events.SearchStateChanged)
			printRows(s, qChangedEv.State, lines)
		case events.ScreenResizeEvent:
			s.Sync()
		case events.EntryFinalSelectEvent:
			finalSelectEvt := ev.(events.EntryFinalSelectEvent)
			fmt.Printf(finalSelectEvt.State.Entry(lines))
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

func printRows(s tcell.Screen, state events.SearchState, entries []string) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.Query), 0, bold)

	for i, l := range state.FilteredLines(entries) {
		if i == state.Selected {
			sc.AppendRow(fmt.Sprintf(">  %s", l), 0, blink)
		} else {
			sc.AppendRow(fmt.Sprintf("   %s", l), 0, plain)
		}
	}
	sc.PrintAll(s)

	s.Sync()
}
