package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/spf13/cobra"
	"github.com/txominpelu/fnd/events"
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

func init() {
	RootCmd.PersistentFlags().StringVar(&lineFormat, "line_format", "plain", "fnd will parse the lines according to this format (plain,json,tabular)")
	RootCmd.PersistentFlags().StringVar(&outputColumn, "output_column", "", "column that will be used as output when picking an element ($ means it outputs the whole row)")
	RootCmd.PersistentFlags().StringVar(&outputTemplate, "output_template", "", "golang template for the output: e.g {{.PID}} means return PID field")
	RootCmd.PersistentFlags().StringVar(&searchType, "search_type", "fuzzy", "type of search (indexed, fuzzy). Indexed is faster for bigger input, fuzzy for finding more matches")
	RootCmd.PersistentFlags().StringSliceVar(&displayColumns, "display_columns", []string{}, "comma separated list of columns to display in order")
}

func runRoot(cmd *cobra.Command, args []string) {

	firstLine := ""
	var scanner *bufio.Scanner
	comesFromStdin := stdinHasPipe()
	if comesFromStdin {
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		firstLine = scanner.Text()
	}
	parser := search.FormatNameToParser(lineFormat, firstLine, displayColumns)
	var searcher search.TextSearcher
	if searchType == "indexed" {
		searcher = index.NewIndexedLines(index.CommandLineTokenizer())
	} else if searchType == "fuzzy" {
		searcher = fuzzy.NewFuzzySearcher()
	} else {
		panic(fmt.Sprintf("search_type should be one of (indexed / fuzzy) it was '%s'", searchType))
	}
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
			filesChannel := listFiles()
			for line := range filesChannel {
				searcher.AddDocument(search.ParseLine(parser, line))
			}
		}
	}()

	initialState := events.SearchState{Query: "", Selected: 0}
	outputCompiledTmpl := compileTemplate(outputColumn, outputTemplate)
	s := initScreen()
	printRows(s, initialState, &searcher, parser.Headers())
	handleEvents(&searcher, s, initialState, parser.Headers(), outputCompiledTmpl)

	s.Fini()
}

func compileTemplate(outputColumn string, outputTemplate string) *template.Template {
	if outputColumn != "" {
		outputTemplate = fmt.Sprintf("{{.%s}}", outputColumn)
	}
	tmpl, err := template.New("test").Parse(outputTemplate)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing output template: %s", err))
	}
	return tmpl
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

func handleEvents(searcher *search.TextSearcher, s tcell.Screen, state events.SearchState, headers []string, outputCompiledTmpl *template.Template) {
	eventChannel := events.NewEventsChannel(s, "", *searcher)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			printRows(s, state, searcher, headers)
		case ev := <-eventChannel:
			state = ev.State()
			switch ev.(type) {
			case events.SearchStateChanged:
				qChangedEv := ev.(events.SearchStateChanged)
				printRows(s, qChangedEv.State(), searcher, headers)
			case events.ScreenResizeEvent:
				s.Sync()
			case events.EntryFinalSelectEvent:
				finalSelectEvt := ev.(events.EntryFinalSelectEvent)
				fmt.Printf(finalSelectEvt.State().Entry(*searcher, outputCompiledTmpl))
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

func printRows(s tcell.Screen, state events.SearchState, searcher *search.TextSearcher, headers []string) {
	s.Clear()
	w, h := s.Size()
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)

	sc := screen.NewScreen(w, h)
	sc.AppendRow(fmt.Sprintf("> %s", state.Query), 0, bold)

	filtered := state.FilteredLines(*searcher)
	sc.AppendRow(fmt.Sprintf("  %d/%d ", len(filtered), (*searcher).Count()), 0, bold)

	t := screen.NewTable(headers)
	for _, l := range filtered {
		t.AddRow(l.ParsedLine)
	}
	t.WriteToScreen(&sc, state.Selected, plain, blink, bold)
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

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()
	return s
}
