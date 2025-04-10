package logger

import (
	il "github.com/faelmori/logz/internal/core"
	"log"
	"os"
)

// LogLevel represents the level of the log entry.
type LogLevel = il.LogLevel

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

type LogzCore interface{ il.LogzCore }

type LogzLogger interface{ il.LogzLogger }

type Logger struct {
	// il.LogzCoreImpl is the logz core logger
	il.LogzCoreImpl
	// il.LogzLogger is the logz logger
	il.LogzLogger
}

// logzLogger is the implementation of the LoggerInterface, unifying the new LogzCoreImpl and the old one.
type logzLogger struct {
	il.LogzLogger

	// log.Logger is the standard Go logger.
	log.Logger

	// logger is the logz logger.
	//logger LogzLogger

	// coreLogger is the LogzCore logger.
	//coreLogger LogzCore
}

// NewLogger creates a new instance of logzLogger with an optional prefix.
func NewLogger(prefix string) LogzLogger {
	lgz := &logzLogger{
		il.NewLogger(prefix),
		*log.New(
			os.Stdout,
			prefix,
			log.LstdFlags,
		),
	}
	lgz.SetPrefix(prefix)
	lgz.SetFlags(log.LstdFlags)
	lgz.SetOutput(os.Stdout)
	lgz.SetLevel(il.INFO)
	return lgz
}

func NewDefaultWriter(out *os.File, formatter LogFormatter) *il.DefaultWriter[any] {
	return il.NewDefaultWriter[any](out, formatter)
}
