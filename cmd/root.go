package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/rjobs/events"
	"github.com/txominpelu/rjobs/index"
	"github.com/txominpelu/rjobs/screen"
)

var RootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Clone of fzf with extended features",
	Long:  `Clone of fzf with extended features`,
	Run:   runRoot,
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

	lines := index.IndexedLines{}
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines.AddLine(scanner.Text())
		}
	}()

	initialState := events.SearchState{Query: "", Selected: 0}
	printRows(s, initialState, &lines)

	handleEvents(&lines, s, initialState)

	s.Fini()
}

func handleEvents(lines *index.IndexedLines, s tcell.Screen, state events.SearchState) {
	eventChannel := events.NewEventsChannel(s, "", lines)
	ticker := time.NewTicker(500 * time.Millisecond)
	count := lines.Count()
	for {
		select {
		case <-ticker.C:
			if lines.Count() != count {
				count = lines.Count()
				printRows(s, state, lines)
			}
		case ev := <-eventChannel:
			state = ev.State()
			switch ev.(type) {
			case events.SearchStateChanged:
				qChangedEv := ev.(events.SearchStateChanged)
				printRows(s, qChangedEv.State(), lines)
			case events.ScreenResizeEvent:
				s.Sync()
			case events.EntryFinalSelectEvent:
				finalSelectEvt := ev.(events.EntryFinalSelectEvent)
				fmt.Printf(finalSelectEvt.State().Entry(lines))
				return
			case events.EscapeEvent:
				return
			}
		}
	}
}

// template:
// >
//  {{#lines}}
//	{{if .selected}}{{style=highlight}}{{.}}{{^style}}{{else}}
//	{{fi}}
//  {{^lines}}

func printRows(s tcell.Screen, state events.SearchState, indexedLines *index.IndexedLines) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.Query), 0, bold)

	filtered := state.FilteredLines(indexedLines)
	sc.AppendRow(fmt.Sprintf("  %d/%d ", len(filtered), indexedLines.Count()), 0, bold)

	for i, l := range filtered {
		if i == state.Selected {
			sc.AppendRow(fmt.Sprintf(">  %s", l), 0, blink)
		} else {
			sc.AppendRow(fmt.Sprintf("   %s", l), 0, plain)
		}
	}
	sc.PrintAll(s)

	s.Sync()
}
