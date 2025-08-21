// Package readers provides functionality to read and process data from various sources.
package readers

import (
	il "github.com/rafa-mori/logz/internal/core"
)

// LogzReader represents a reader for log entries.
type LogzReader = il.LogReader

// LogzFileReader represents a file reader for log entries.
type LogzFileReader = il.FileLogReader

// NewLogzReader creates a new instance of LogzReader.
// It initializes a reader that can read log entries from a source.
// This allows for reading log entries from various formats and sources.
func NewLogzReader() LogzReader {
	return il.NewFileLogReader()
}
