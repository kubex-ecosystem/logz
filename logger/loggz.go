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
	log.Logger
	// il.LogzLogger is the logz logger
	LogzLogger
}

// logzLogger is the implementation of the LoggerInterface, unifying the new LogzCoreImpl and the old one.
type logzLogger struct {
	// logger is the logz logger.
	//logger LogzLogger
	*Logger

	//il.LogzLogger

	// coreLogger is the LogzCore logger.
	//LogzCore
}

// NewLogger creates a new instance of logzLogger with an optional prefix.
func NewLogger(prefix string) LogzLogger {
	lgzR := &Logger{
		*log.New(
			os.Stdout,
			prefix,
			log.LstdFlags,
		),
		il.NewLogger(prefix),
	}
	lgz := &logzLogger{
		lgzR,
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
