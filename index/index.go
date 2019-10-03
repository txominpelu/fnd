package index

import (
	"strings"
)

type IndexedLines struct {
	lines          []string
	NewLineChannel chan string // Channel that notifies about new lines
}

func (i *IndexedLines) AddLine(line string) {
	if i.lines == nil {
		i.lines = []string{}
	}
	if i.NewLineChannel == nil {
		i.NewLineChannel = make(chan string)
	}
	i.lines = append(i.lines, line)
	i.NewLineChannel <- line
}

func (i IndexedLines) FilterEntries(query string) []string {
	fLines := []string{}
	for _, l := range i.lines {
		if strings.Contains(l, query) {
			fLines = append(fLines, l)
		}
	}
	return fLines
}
