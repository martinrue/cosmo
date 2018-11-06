package table

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

// AddRow adds a row of arbitrary columns to the table.
func (t *Table) AddRow(cols ...string) {
	t.Rows = append(t.Rows, cols)

	if len(t.Rows) == 1 {
		t.cols = len(cols)
	}

	if len(cols) != t.cols {
		panic("table cannot have variable width rows")
	}
}

func (t *Table) String() string {
	if len(t.Rows) == 0 {
		return ""
	}

	padding := 0

	if len(t.Rows) > 1 {
		padding = 1
	}

	table := strings.Builder{}

	for ri, row := range t.Rows {
		cols := make([]string, 0)

		for ci, col := range row {
			format := fmt.Sprintf("%%-%dv ", t.colWidth(ci)+padding)
			cols = append(cols, fmt.Sprintf(format, col))
		}

		table.WriteString(strings.TrimSpace(strings.Join(cols, "")))

		if ri != len(t.Rows)-1 {
			table.WriteString("\n")
		}
	}

	return table.String()
}
