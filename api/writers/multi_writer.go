// Package writters provides functionality for managing writers.
package writers

import "github.com/kubex-ecosystem/logz/internal/loggerz"

// LogMultiWriter defines the contract for writing logs to multiple writers.
type LogMultiWriter[T any] interface {
	Write(entry []byte) (int, error)
	AddWriter(w LogzWriter[T, WriteLogz[T]])
	GetWriters() []LogzWriter[T, WriteLogz[T]]
}

// NewLogzMultiWriter creates a new instance of LogzMultiWriter.
// It initializes a multi-writer that can handle multiple log writers.
// This allows for writing log entries to multiple outputs simultaneously.
func NewLogzMultiWriter[T any]() LogMultiWriter[T] {
	return loggerz.NewMultiWriter()
}
