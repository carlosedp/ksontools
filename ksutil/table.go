package ksutil

import (
	"fmt"
	"io"
	"strings"
)

const (
	// sepChar is the character used to separate the header from the content in a table.
	sepChar = "="
)

// Table creates an output table.
type Table struct {
	w io.Writer

	header []string
	rows   [][]string
}

// NewTable creates an instance of table.
func NewTable(w io.Writer) *Table {
	return &Table{
		w: w,
	}
}

// SetHeader sets the header for the table.
func (t *Table) SetHeader(columns []string) {
	t.header = columns
}

// Append appends a row to the table.
func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

// AppendBulk appends multiple rows to the table.
func (t *Table) AppendBulk(rows [][]string) {
	t.rows = append(t.rows, rows...)
}

// Render writes the output to the table's writer.
func (t *Table) Render() {
	var output [][]string

	if len(t.header) > 0 {
		headerRow := make([]string, len(t.header), len(t.header))
		sepRow := make([]string, len(t.header), len(t.header))

		for i := range t.header {
			sepLen := len(t.header[i])
			headerRow[i] = strings.ToUpper(t.header[i])
			sepRow[i] = strings.Repeat(sepChar, sepLen)
		}

		output = append(output, headerRow, sepRow)
	}

	output = append(output, t.rows...)

	// count the number of columns
	colCount := 0
	for _, row := range output {
		if l := len(row); l > colCount {
			colCount = l
		}
	}

	// get the max len for each column
	counts := make([]int, colCount, colCount)
	for _, row := range output {
		for i := range row {
			if l := len(row[i]); l > counts[i] {
				counts[i] = l
			}
		}
	}

	// print rows
	for _, row := range output {
		var parts []string
		for i, col := range row {
			val := col
			if i < len(row)-1 {
				format := fmt.Sprintf("%%-%ds", counts[i])
				val = fmt.Sprintf(format, col)
			}
			parts = append(parts, val)

		}
		fmt.Fprintf(t.w, "%s\n", strings.Join(parts, " "))
	}
}
