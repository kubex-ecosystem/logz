// Package logger provides a unified interface for logging with various configurations and formats.
package logger

import (
	"os"

	il "github.com/kubex-ecosystem/logz/internal/core"
)

// LogFormat represents the format of the log entry.
type LogFormat = il.LogFormat

// LogWriter represents the writer of the log entry.
type LogWriter[T any] interface{ il.LogWriter[T] }

// Config represents the configuration of the il.
type Config interface{ il.Config }

// LogzEntry represents a single log entry with various attributes.
type LogzEntry interface{ il.LogzEntry }

// LogFormatter defines the contract for formatting log entries.
type LogFormatter interface{ il.LogFormatter }

type logxLogger = il.LogzCoreImpl

type LogzLogger = il.LogzLogger

// NewLogger creates a new instance of logzLogger with an optional prefix.
func NewLogger(prefix string) LogzLogger {
	return il.NewLogger(prefix)
}

func NewDefaultWriter(out *os.File, formatter LogFormatter) *il.DefaultWriter[any] {
	return il.NewDefaultWriter[any](out, formatter)
}
