// Package utils provides utility functions for handling tables.
package utils

import (
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"

	"os"
	"path/filepath"
)

// Table represents a simple table.
type Table struct {
	data [][]string
}

// NewTable creates a new simple table.
func NewTable(data [][]string) Table {
	return Table{data}
}

func getDefaultConfig() *tablewriter.Config {
	config := tablewriter.NewConfigBuilder().
		WithAutoHide(tw.On).
		WithHeaderAlignment(tw.AlignLeft).
		WithRowAlignment(tw.AlignLeft).
		Build()

	return &config
}

// PrintTable prints the simple table in the shell with side and vertical borders.
func (t Table) PrintTable() {
	config := getDefaultConfig()

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithConfig(*config),
	)

	for _, row := range t.data {
		table.Append(row)
	}

	table.Render()
}

// FormattedTable represents a formatted table.
type FormattedTable struct {
	data   [][]string
	header []string
}

// NewFormattedTable creates a new formatted table.
func NewFormattedTable(header []string, data [][]string) FormattedTable {
	return FormattedTable{header: header, data: data}
}

// SaveFormattedTable saves the formatted table to a file.
func (ft FormattedTable) SaveFormattedTable(filename string) error {
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	table := tablewriter.NewWriter(file)

	table.Header(ft.header)

	for _, row := range ft.data {
		table.Append(row)
	}
	table.Render()
	return nil
}
