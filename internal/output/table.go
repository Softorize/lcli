package output

import (
	"fmt"
	"io"
	"strings"
)

// Table holds tabular data for aligned column rendering.
type Table struct {
	// Headers are the column names displayed at the top of the table.
	Headers []string
	// Rows holds the data rows, each being a slice of column values.
	Rows [][]string
}

// NewTable creates a Table with the given column headers.
func NewTable(headers ...string) *Table {
	return &Table{Headers: headers}
}

// AddRow appends a row of column values to the table.
func (t *Table) AddRow(cols ...string) {
	t.Rows = append(t.Rows, cols)
}

// Render writes the table to w with aligned columns and a header separator.
func (t *Table) Render(w io.Writer) {
	widths := t.columnWidths()
	t.renderRow(w, t.Headers, widths)
	t.renderSeparator(w, widths)
	for _, row := range t.Rows {
		t.renderRow(w, row, widths)
	}
}

// columnWidths returns the maximum width needed for each column.
func (t *Table) columnWidths() []int {
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, col := range row {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}
	return widths
}

// renderRow writes a single row with padded columns.
func (t *Table) renderRow(w io.Writer, cols []string, widths []int) {
	parts := make([]string, len(widths))
	for i := range widths {
		val := ""
		if i < len(cols) {
			val = cols[i]
		}
		parts[i] = fmt.Sprintf("%-*s", widths[i], val)
	}
	fmt.Fprintf(w, "%s\n", strings.Join(parts, "  "))
}

// renderSeparator writes a dashed line under the header row.
func (t *Table) renderSeparator(w io.Writer, widths []int) {
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("-", width)
	}
	fmt.Fprintf(w, "%s\n", strings.Join(parts, "  "))
}
