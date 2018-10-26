package drawing

import (
	"fmt"
	"strings"
)

// Table handles the drawing of tabular data.
type Table struct {
	Rows [][]string
	cols int
}

func (t *Table) colWidth(col int) int {
	width := 0

	for _, row := range t.Rows {
		if len(row[col]) > width {
			width = len(row[col])
		}
	}

	return width
}

func (t *Table) addRow(cols ...string) {
	t.Rows = append(t.Rows, cols)

	if len(t.Rows) == 1 {
		t.cols = len(cols)
	}

	if len(cols) != t.cols {
		panic("table cannot have variable width rows")
	}
}

// AddRow adds a row to the table, inserting multiple if a column contains newlines.
func (t *Table) AddRow(cols ...string) {
	var prefixes []string
	var rows int

	for i, col := range cols {
		if strings.Contains(col, "\n") {
			prefixes = cols[:i]
			rows = len(strings.Split(col, "\n"))
		}
	}

	if rows == 0 {
		t.addRow(cols...)
		return
	}

	for i := 0; i < rows; i++ {
		c := make([]string, 0)

		if i == 0 {
			c = append(c, prefixes...)
		} else {
			c = append(c, make([]string, len(prefixes))...)
		}

		for _, col := range cols[len(prefixes):] {
			c = append(c, strings.Split(col, "\n")[i])
		}

		t.addRow(c...)
	}
}

func (t *Table) String() string {
	for ri, row := range t.Rows {
		for ci, col := range row {
			format := fmt.Sprintf("%%-%dv ", t.colWidth(ci))
			fmt.Printf(format, col)
		}

		if ri != len(t.Rows)-1 {
			fmt.Println()
		}
	}

	return ""
}
