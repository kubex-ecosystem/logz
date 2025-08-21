// Package writters provides functionality for managing writers.
package writters

import (
	il "github.com/rafa-mori/logz/internal/core"
)

// LogzMultiWriter represents a multi-writer for log entries.
type LogzMultiWriter[T any] = il.LogMultiWriter[T]

// NewLogzMultiWriter creates a new instance of LogzMultiWriter.
// It initializes a multi-writer that can handle multiple log writers.
// This allows for writing log entries to multiple outputs simultaneously.
func NewLogzMultiWriter[T any]() LogzMultiWriter[T] {
	return il.NewMultiWriter[T]()
}
