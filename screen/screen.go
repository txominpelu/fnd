package screen

import (
	"strings"

	"github.com/gdamore/tcell"
)

type ContentBlock struct {
	r     rune
	style tcell.Style
}

type Row struct {
	blocks []ContentBlock
	width  int
}

func (row *Row) writeRune(r rune, x int, style tcell.Style) {
	row.blocks[x].r = r
	row.blocks[x].style = style
}

func (row *Row) writeString(s string, x int, style tcell.Style) {
	i := 0
	for _, char := range s {
		if i+x >= row.width {
			break
		}
		row.writeRune(char, i+x, style)
		i++
	}
}

func newRow(width int) Row {
	return Row{
		width:  width,
		blocks: make([]ContentBlock, width),
	}
}

type Screen struct {
	rows   []Row
	width  int
	height int
}

func NewScreen(width int, height int) Screen {
	return Screen{
		width:  width,
		height: height,
	}
}

func (sc *Screen) setRune(x int, y int, r rune, style tcell.Style) {
	sc.rows[y].blocks[x].r = r
	sc.rows[y].blocks[x].style = style
}

// appends a row at height = current_max_height + 1
func (sc *Screen) AppendRow(s string, x int, style tcell.Style) {
	r := newRow(sc.width)
	r.writeString(s, x, style)
	sc.rows = append(sc.rows, r)
}

func (sc *Screen) PrintAll(s tcell.Screen) {
	for y, r := range sc.rows {
		for x, b := range r.blocks {
			s.SetContent(x, sc.height-(y+1), b.r, []rune{}, b.style)
		}
	}
}

//To facilitate debug
func (sc *Screen) toString() string {
	all := [][]string{}
	for i := 0; i < sc.height; i++ {
		emptyRow := make([]string, sc.width)
		for j := 0; j < sc.width; j++ {
			emptyRow[j] = " "
		}
		all = append(all, emptyRow)
	}
	for y, r := range sc.rows {
		for x, b := range r.blocks {
			all[sc.height-(y+1)][x] = string(b.r)
		}
	}
	s := strings.Builder{}
	for _, l := range all {
		s.WriteString(strings.Join(l, ""))
		s.WriteString("\n")
	}
	return s.String()
}
