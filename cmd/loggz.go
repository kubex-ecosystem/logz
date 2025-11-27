// Package logger provides a unified interface for logging with various configurations and formats.
package main

import (
	"os"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/writer"
)

// type ILogWriter[T any, W il.WriteLog[T]] interface{ il.LogWriter[T,W] }

// type LogFormat = LogFormat
type LogWriter = writer.IOWriter

type Writer = interfaces.Writer

// type Config = interfaces.Config
type LogzEntry = interfaces.Entry
type LogFormatter = interfaces.Formatter
type logxLogger = interfaces.Logger

func NewLogger(prefix string, opts *core.LoggerOptionsImpl, withDefaults bool) *core.Logger {
	return core.NewLogger(prefix, opts, withDefaults)
}
func NewDefaultWriter(out *os.File, formatter LogFormatter) writer.Writer {
	return writer.NewIOWriter(os.Stdout)
}
