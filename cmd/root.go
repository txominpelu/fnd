package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/fnd/events"
	"github.com/txominpelu/fnd/index"
	"github.com/txominpelu/fnd/screen"
)

var RootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Clone of fzf with extended features",
	Long:  `Clone of fzf with extended features`,
	Run:   runRoot,
}

var lineOutput string

func init() {
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "line_output", "", "what will be the output of choosing a line (jq format)")
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

	lines := index.NewIndexedLines(index.CommandLineTokenizer(), index.PlainTextParser())
	go func() {
		if stdinHasPipe() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				lines.AddLine(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("Error: %s while reading stdin", err))
			}
		} else {
			filesChannel := listFiles()
			for line := range filesChannel {
				lines.AddLine(line)
			}
		}
	}()

	initialState := events.SearchState{Query: "", Selected: 0}
	printRows(s, initialState, &lines)
	handleEvents(&lines, s, initialState)

	s.Fini()
}

func listFiles() chan string {
	out := make(chan string)
	go func() {
		if isGitFolder() {
			for _, f := range gitLs() {
				out <- f
			}
			close(out)
		} else {
			filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				out <- path
				return nil
			})
			close(out)
		}
		//FIXME: log error if err != nil
	}()
	return out
}

func gitLs() []string {
	args := strings.Fields("git ls-files")
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		//FIXME: log error
		return []string{}
	}
	return strings.Split(string(out), "\n")
}

func isGitFolder() bool {
	args := strings.Fields("git rev-parse --is-inside-work-tree")
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func stdinHasPipe() bool {

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}

func handleEvents(lines *index.IndexedLines, s tcell.Screen, state events.SearchState) {
	eventChannel := events.NewEventsChannel(s, "", lines)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			printRows(s, state, lines)
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
				close(eventChannel)
				return
			case events.EscapeEvent:
				close(eventChannel)
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
