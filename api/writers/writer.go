// Package writters provides functionality for managing writers.
package writters

import (
	"io"

	il "github.com/rafa-mori/logz/internal/core"
)

// LogzWriter represents the writer of the log entry.
type LogzWriter[T any] = il.LogWriter[T]

// LogzFormatter represents the formatter for the log entry.
type LogzFormatter = il.LogFormatter

// NewLogzWriter creates a new instance of LogzWriter with the given writer.
func NewLogzWriter[T any](out io.Writer, formatter LogzFormatter) LogzWriter[T] {
	return il.NewDefaultWriter[T](out, formatter)
}
