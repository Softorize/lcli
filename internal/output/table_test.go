package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewTable(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
	}{
		{"single header", []string{"ID"}},
		{"multiple headers", []string{"ID", "Name", "Email"}},
		{"empty headers", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := NewTable(tt.headers...)
			if len(tbl.Headers) != len(tt.headers) {
				t.Errorf("header count = %d, want %d", len(tbl.Headers), len(tt.headers))
			}
			if len(tbl.Rows) != 0 {
				t.Errorf("initial rows = %d, want 0", len(tbl.Rows))
			}
		})
	}
}

func TestAddRow(t *testing.T) {
	tbl := NewTable("A", "B")
	tbl.AddRow("1", "2")
	tbl.AddRow("3", "4")

	if len(tbl.Rows) != 2 {
		t.Fatalf("row count = %d, want 2", len(tbl.Rows))
	}
	if tbl.Rows[0][0] != "1" || tbl.Rows[0][1] != "2" {
		t.Errorf("row 0 = %v, want [1 2]", tbl.Rows[0])
	}
	if tbl.Rows[1][0] != "3" || tbl.Rows[1][1] != "4" {
		t.Errorf("row 1 = %v, want [3 4]", tbl.Rows[1])
	}
}

func TestRenderAlignedColumns(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		rows    [][]string
		checks  []string
	}{
		{
			name:    "basic alignment",
			headers: []string{"ID", "Name"},
			rows:    [][]string{{"1", "Alice"}, {"200", "Bo"}},
			checks:  []string{"ID ", "Name", "---", "Alice", "200"},
		},
		{
			name:    "longer data than header",
			headers: []string{"X"},
			rows:    [][]string{{"LongValue"}},
			checks:  []string{"X        ", "LongValue"},
		},
		{
			name:    "empty table",
			headers: []string{"Col"},
			rows:    [][]string{},
			checks:  []string{"Col", "---"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := NewTable(tt.headers...)
			for _, row := range tt.rows {
				tbl.AddRow(row...)
			}
			var buf bytes.Buffer
			tbl.Render(&buf)
			output := buf.String()
			for _, want := range tt.checks {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestRenderSeparatorLength(t *testing.T) {
	tbl := NewTable("Name", "Email")
	tbl.AddRow("Alice", "alice@example.com")

	var buf bytes.Buffer
	tbl.Render(&buf)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	// Separator line (index 1) should have same length as header line
	headerLen := len(lines[0])
	sepLen := len(lines[1])
	if headerLen != sepLen {
		t.Errorf("header len = %d, separator len = %d", headerLen, sepLen)
	}
}
