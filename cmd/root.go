package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/fnd/events"
	"github.com/txominpelu/fnd/log"
	"github.com/txominpelu/fnd/screen"
	"github.com/txominpelu/fnd/search"
	"github.com/txominpelu/fnd/search/fuzzy"
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
var outputTemplate string
var searchType string
var displayColumns []string
var logFile string
var sorterName string

func init() {
	RootCmd.PersistentFlags().StringVar(&lineFormat, "line_format", "plain", "fnd will parse the lines according to this format (plain,json,tabular)")
	RootCmd.PersistentFlags().StringVar(&outputColumn, "output_column", "$", "column that will be used as output when picking an element ($ means it outputs the whole row)")
	RootCmd.PersistentFlags().StringVar(&outputTemplate, "output_template", "", "golang template for the output: e.g {{.PID}} means return PID field")
	RootCmd.PersistentFlags().StringVar(&searchType, "search_type", "fuzzy", "type of search (indexed, fuzzy). Indexed is faster for bigger input, fuzzy for finding more matches")
	RootCmd.PersistentFlags().StringVar(&logFile, "log_file", "", "errors will be logged to the given file")
	RootCmd.PersistentFlags().StringSliceVar(&displayColumns, "display_columns", []string{}, "comma separated list of columns to display in order")
	RootCmd.PersistentFlags().StringVar(&sorterName, "sorter", "default", " sorter (index/default) ")
}

func runRoot(cmd *cobra.Command, args []string) {

	s := initScreen()
	// make sure to always clean the screen
	defer func() {
		s.Fini()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	firstLine := ""
	var scanner *bufio.Scanner
	comesFromStdin := stdinHasPipe()
	if comesFromStdin {
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		firstLine = scanner.Text()
	}
	logger := log.NewLogger(logFile)
	searcher, err := getSearcher(searchType)
	sorter := getSorter(searcher, sorterName)
	logger.CheckError(err, "when parsing search_type flag")

	parser := search.FormatNameToParser(lineFormat, firstLine, displayColumns, logger)
	if comesFromStdin && lineFormat != "tabular" {
		searcher.AddDocument(search.ParseLine(parser, firstLine))
	}

	go func() {
		if comesFromStdin {
			for scanner.Scan() {
				searcher.AddDocument(search.ParseLine(parser, scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("Error: %s while reading stdin", err))
			}
		} else {
			filesChannel := listFiles(logger)
			for line := range filesChannel {
				searcher.AddDocument(search.ParseLine(parser, line))
			}
		}
	}()

	initialState := events.SearchState{Query: "", Selected: 0}
	renderer := getRenderer(outputColumn, outputTemplate, logger)
	printRows(s, initialState, &searcher, parser.Headers(), sorter)
	handleEvents(&searcher, s, initialState, parser.Headers(), renderer, sorter)

	s.Fini()
}

func getSorter(searcher search.TextSearcher, sorter string) search.Compare {
	if sorter == "index" {
		return func(d1 int, d2 int) bool {
			return d2 < d1
		}
	} else {
		return func(d1 int, d2 int) bool {
			if len(searcher.GetDocById(d1).RawText) != len(searcher.GetDocById(d2).RawText) {
				t1 := searcher.GetDocById(d1).RawText
				t2 := searcher.GetDocById(d2).RawText
				return len(t1) < len(t2)
			} else {
				return d1 < d2
			}
		}

	}

}

func getSearcher(searchType string) (search.TextSearcher, error) {
	if searchType == "indexed" {
		return index.NewIndexedLines(index.CommandLineTokenizer()), nil
	} else if searchType == "fuzzy" {
		return fuzzy.NewFuzzySearcher(), nil
	} else {
		return nil, fmt.Errorf("search_type should be one of (indexed / fuzzy) it was '%s'", searchType)
	}
}

func listFiles(logger *log.StandardLogger) chan string {
	out := make(chan string)
	go func() {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			out <- path
			return nil
		})
		close(out)
		logger.CheckError(err, "when iterating over files recursively")
	}()
	return out
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

func handleEvents(searcher *search.TextSearcher, s tcell.Screen, state events.SearchState, headers []string, renderer renderOutput, sorter search.Compare) {
	eventChannel := events.NewEventsChannel(s, "", *searcher, sorter)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			printRows(s, state, searcher, headers, sorter)
		case ev := <-eventChannel:
			state = ev.State()
			switch ev.(type) {
			case events.SearchStateChanged:
				qChangedEv := ev.(events.SearchStateChanged)
				printRows(s, qChangedEv.State(), searcher, headers, sorter)
			case events.ScreenResizeEvent:
				s.Sync()
			case events.EntryFinalSelectEvent:
				finalSelectEvt := ev.(events.EntryFinalSelectEvent)
				fmt.Printf(renderer(finalSelectEvt.State().Entry(*searcher, sorter).ParsedLine))
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

func printRows(s tcell.Screen, state events.SearchState, searcher *search.TextSearcher, headers []string, sorter search.Compare) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault.Normal()
	bold := tcell.StyleDefault.Normal().Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.Query), 0, bold)

	filtered := state.FilteredLines(*searcher, sorter)
	sc.AppendRow(fmt.Sprintf("  %d/%d ", len(filtered), (*searcher).Count()), 0, bold)

	t := screen.NewTable(headers)
	for _, l := range filtered {
		t.AddRow(l.ParsedLine)
	}
	t.WriteToScreen(&sc, state.Selected, plain, bold, bold)
	sc.PrintAll(s)

	s.Sync()
}

func initScreen() tcell.Screen {
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

	s.Clear()
	return s
}
