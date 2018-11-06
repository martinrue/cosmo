package table_test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/martinrue/cosmo/table"
)

func TestPrintPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected variable width rows panic, but no panic")
		}
	}()

	table := &table.Table{}
	table.AddRow("col 1")
	table.AddRow("col 1", "col 2")
}

func TestPrint(t *testing.T) {
	stringMatchesGoldenFile := func(t *testing.T, str string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", filename))
		if err != nil {
			t.Fatal(err)
		}

		return str == string(golden)
	}

	tests := []struct {
		Name       string
		Rows       [][]string
		GoldenFile string
	}{
		{"no rows", [][]string{}, "no-rows.golden"},
		{"1 row", [][]string{{"one", "two", "three"}}, "1-row.golden"},
		{"2 rows", [][]string{{"one", "two"}, {"black", "blue"}}, "2-rows.golden"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			table := &table.Table{}

			for _, row := range test.Rows {
				table.AddRow(row...)
			}

			if !stringMatchesGoldenFile(t, table.String(), test.GoldenFile) {
				t.Fatalf("table output does not match golden file")
			}
		})
	}
}
