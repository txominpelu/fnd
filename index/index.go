package index

import (
	"strings"
)

type IndexedLines struct {
	lines []string
	count int
}

func (i *IndexedLines) AddLine(line string) {
	if i.lines == nil {
		i.lines = []string{}
	}
	i.lines = append(i.lines, line)
	i.count++
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

func (i IndexedLines) Count() int {
	return i.count
}
