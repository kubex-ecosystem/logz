// Package writters provides functionality for managing writers.
package writers

import (
	"io"

	il "github.com/kubex-ecosystem/logz/internal/core"
)

type WriteLogz[E any] = func(entry E, ioWriter func([]byte) error) error

// LogzWriter defines the contract for writing logs.
type LogzWriter[E any, T il.WriteLogz[E]] interface {
	Write(entry []byte) (int, error)
	WriteLogz(entry E, ioWriter T) (int, error)
}

type Writer = il.Writer

// LogzFormatter represents the formatter for the log entry.
type LogzFormatter = il.LogFormatter

// NewLogzWriter creates a new instance of LogzWriter with the given writer.
func NewLogzWriter[E any, T il.WriteLogz[E]](typ E, output T, format string, writer io.Writer) LogzWriter[E, T] {
	return il.NewLogzWriter[E, T](typ, output, format, writer)
}

// NewDefaultWriter creates a new instance of LogzWriter with the given writer.
func NewDefaultWriter[T any](out io.Writer, formatter LogzFormatter) Writer {
	return il.NewDefaultWriter[[]byte](out, formatter)
}
