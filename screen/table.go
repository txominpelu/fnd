package screen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell"
)

type Table struct {
	columns []string
	rows    []map[string]string
}

func NewTable(headers []string) Table {
	return Table{
		columns: headers,
		rows:    []map[string]string{},
	}
}

func (t *Table) AddRow(row map[string]string) {
	t.rows = append(t.rows, row)
}

func (t Table) computeWidths(width int) map[string]int {
	columnToWidth := map[string]int{}
	for _, c := range t.columns {
		columnToWidth[c] = t.maxWidth(c)
	}
	return columnToWidth
}

func (t Table) maxWidth(column string) int {
	max := utf8.RuneCountInString(column)
	for _, r := range t.rows {
		if max < utf8.RuneCountInString(r[column]) {
			max = utf8.RuneCountInString(r[column])
		}
	}
	return max
}

func (t Table) WriteToScreen(sc *Screen, selected int, plainStyle tcell.Style, selectedStyle tcell.Style, boldStyle tcell.Style) {
	// leftPaddingLength is require to have a space when listing elements to do '>' for the selected one
	leftPaddingLength := 2
	//TODO: allow trimming if all columns together get out of screen
	var columnToWidth map[string]int = t.computeWidths(sc.width - leftPaddingLength)
	for i, row := range t.rows {
		rowString := t.buildRowString(row, columnToWidth)
		if i == selected {
			sc.AppendRow(fmt.Sprintf("> %s", rowString), 0, selectedStyle)
		} else {
			sc.AppendRow(fmt.Sprintf("  %s", rowString), 0, plainStyle)
		}
		// 4 = headers line + query line + counter line + initial line
		if i+4 >= sc.height {
			break
		}
	}
	columns := map[string]string{}
	for _, column := range t.columns {
		columns[column] = column
	}
	headersString := t.buildRowString(columns, columnToWidth)
	sc.AppendRow(fmt.Sprintf("  %s", headersString), 0, boldStyle)
}

func (t Table) buildRowString(row map[string]string, columnToWidth map[string]int) string {
	sBuilder := strings.Builder{}
	for _, column := range t.columns {
		val := row[column]
		rightPaddingLen := (columnToWidth[column] + 1) - utf8.RuneCountInString(val)
		fieldValue := strings.Join([]string{val, strings.Repeat(" ", rightPaddingLen)}, "")
		sBuilder.WriteString(fieldValue)
	}
	return sBuilder.String()
}
