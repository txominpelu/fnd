package screen

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell"
)

func TestPrintBasic(t *testing.T) {
	table := Table{
		columns: []string{"column1", "column2"},
		rows: []map[string]string{
			map[string]string{
				"column1": "hello world",
				"column2": "this is a really long text",
			},
			map[string]string{
				"column1": "h",
				"column2": "this is really short",
			},
		},
	}
	sc := NewScreen(40, 3)
	plain := tcell.StyleDefault
	blink := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	bold := tcell.StyleDefault.Bold(true)
	table.WriteToScreen(&sc, 0, plain, blink, bold)
	fmt.Println("Screen:")
	fmt.Print(sc.toString())
	//if !reflect.DeepEqual(ev.State(), expected) {
	//	t.Errorf("expected: '%v' got: '%v'\n", expected, ev.State())
	//}

}
