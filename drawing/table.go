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

// AddRow adds a row to the table.
func (t *Table) AddRow(cols ...string) {
	t.Rows = append(t.Rows, cols)

	if len(t.Rows) == 1 {
		t.cols = len(cols)
	}

	if len(cols) != t.cols {
		panic("table cannot have variable width rows")
	}
}

// AddRows adds a row to the table for each line in the argument.
func (t *Table) AddRows(prefix string, cols ...string) {
	rows := len(strings.Split(cols[0], "\n"))

	for i := 0; i < rows; i++ {
		c := make([]string, 0)

		if i == 0 {
			c = append(c, prefix)
		} else {
			c = append(c, "")
		}

		for _, col := range cols {
			c = append(c, strings.Split(col, "\n")[i])
		}

		t.AddRow(c...)
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
