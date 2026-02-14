// Package output provides formatted output for CLI results.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Format specifies the output serialization format.
type Format string

const (
	// FormatJSON outputs pretty-printed JSON.
	FormatJSON Format = "json"
	// FormatTable outputs aligned columns with headers.
	FormatTable Format = "table"
	// FormatYAML outputs YAML.
	FormatYAML Format = "yaml"
)

// Printer writes structured data in the configured format.
type Printer struct {
	w      io.Writer
	format Format
}

// NewPrinter creates a Printer that writes to w using the given format.
func NewPrinter(w io.Writer, format Format) *Printer {
	return &Printer{w: w, format: format}
}

// Format returns the printer's configured output format.
func (p *Printer) Format() Format {
	return p.format
}

// Print dispatches to the appropriate formatter based on the configured format.
func (p *Printer) Print(v any) error {
	switch p.format {
	case FormatJSON:
		return p.PrintJSON(v)
	case FormatYAML:
		return p.PrintYAML(v)
	case FormatTable:
		return fmt.Errorf("use PrintTable for table output")
	default:
		return p.PrintJSON(v)
	}
}

// PrintJSON writes v as pretty-printed JSON with two-space indentation.
func (p *Printer) PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	_, err = fmt.Fprintf(p.w, "%s\n", data)
	return err
}

// PrintYAML writes v as YAML.
func (p *Printer) PrintYAML(v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}
	_, err = p.w.Write(data)
	return err
}

// PrintTable renders a table with the given headers and rows.
func (p *Printer) PrintTable(headers []string, rows [][]string) error {
	t := NewTable(headers...)
	for _, row := range rows {
		t.AddRow(row...)
	}
	t.Render(p.w)
	return nil
}

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "json":
		return FormatJSON, nil
	case "table":
		return FormatTable, nil
	case "yaml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("unknown format: %q (use json, table, or yaml)", s)
	}
}
