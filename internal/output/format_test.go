package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"table", FormatTable, false},
		{"yaml", FormatYAML, false},
		{"csv", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	type item struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	tests := []struct {
		name  string
		input any
	}{
		{"struct", item{Name: "Alice", Age: 30}},
		{"map", map[string]string{"key": "value"}},
		{"slice", []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewPrinter(&buf, FormatJSON)
			if err := p.PrintJSON(tt.input); err != nil {
				t.Fatalf("PrintJSON: %v", err)
			}
			output := strings.TrimSpace(buf.String())
			if !json.Valid([]byte(output)) {
				t.Errorf("output is not valid JSON: %s", output)
			}
		})
	}
}

func TestPrintTable(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, FormatTable)

	headers := []string{"ID", "Name"}
	rows := [][]string{
		{"1", "Alice"},
		{"2", "Bob"},
	}
	if err := p.PrintTable(headers, rows); err != nil {
		t.Fatalf("PrintTable: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 4 { // header + separator + 2 data rows
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), output)
	}
	if !strings.Contains(lines[0], "ID") {
		t.Errorf("header line should contain ID: %q", lines[0])
	}
	if !strings.Contains(lines[0], "Name") {
		t.Errorf("header line should contain Name: %q", lines[0])
	}
	if !strings.Contains(lines[1], "--") {
		t.Errorf("separator line should contain dashes: %q", lines[1])
	}
	if !strings.Contains(lines[2], "Alice") {
		t.Errorf("first data row should contain Alice: %q", lines[2])
	}
}

func TestPrinterFormat(t *testing.T) {
	p := NewPrinter(nil, FormatYAML)
	if p.Format() != FormatYAML {
		t.Errorf("Format() = %q, want %q", p.Format(), FormatYAML)
	}
}
