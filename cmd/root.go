package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/fnd/events"
	"github.com/txominpelu/fnd/screen"
	"github.com/txominpelu/fnd/search"
	"github.com/txominpelu/fnd/search/index"
)

var RootCmd = &cobra.Command{
	Use:   "fnd",
	Short: "Clone of fzf with extended features",
	Long:  `Clone of fzf with extended features`,
	Run:   runRoot,
}

var lineFormat string
var outputColumn string

func init() {
	RootCmd.PersistentFlags().StringVar(&lineFormat, "line_format", "plain", "fnd will parse the lines according to this format (plain,json,tabular)")
	RootCmd.PersistentFlags().StringVar(&outputColumn, "output_column", "$", "column that will be used as output when picking an element ($ means it outputs the whole row)")
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

	firstLine := ""
	var scanner *bufio.Scanner
	comesFromStdin := stdinHasPipe()
	if comesFromStdin {
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		firstLine = scanner.Text()
	}
	parser := search.FormatNameToParser(lineFormat, firstLine)
	lines := index.NewIndexedLines(
		index.CommandLineTokenizer(),
	)
	if comesFromStdin && lineFormat != "tabular" {
		lines.AddDocument(search.ParseLine(parser, firstLine))
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	go func() {
		if comesFromStdin {
			for scanner.Scan() {
				lines.AddDocument(search.ParseLine(parser, scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("Error: %s while reading stdin", err))
			}
		} else {
			filesChannel := listFiles()
			for line := range filesChannel {
				lines.AddDocument(search.ParseLine(parser, line))
			}
		}
	}()

	initialState := events.SearchState{Query: "", Selected: 0}
	printRows(s, initialState, &lines, parser.Headers())
	handleEvents(&lines, s, initialState, parser.Headers(), outputColumn)

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
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines
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

func handleEvents(lines *index.IndexedLines, s tcell.Screen, state events.SearchState, headers []string, outputColumn string) {
	eventChannel := events.NewEventsChannel(s, "", lines)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			printRows(s, state, lines, headers)
		case ev := <-eventChannel:
			state = ev.State()
			switch ev.(type) {
			case events.SearchStateChanged:
				qChangedEv := ev.(events.SearchStateChanged)
				printRows(s, qChangedEv.State(), lines, headers)
			case events.ScreenResizeEvent:
				s.Sync()
			case events.EntryFinalSelectEvent:
				finalSelectEvt := ev.(events.EntryFinalSelectEvent)
				fmt.Printf(finalSelectEvt.State().Entry(lines, outputColumn))
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

func printRows(s tcell.Screen, state events.SearchState, indexedLines *index.IndexedLines, headers []string) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.Query), 0, bold)

	filtered := state.FilteredLines(indexedLines)
	sc.AppendRow(fmt.Sprintf("  %d/%d ", len(filtered), indexedLines.Count()), 0, bold)

	t := screen.NewTable(headers)
	for _, l := range filtered {
		t.AddRow(l.ParsedLine)
	}
	t.WriteToScreen(&sc, state.Selected, plain, blink, bold)
	sc.PrintAll(s)

	s.Sync()
}
