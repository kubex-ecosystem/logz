// Package logz provides a global logging utility with configurable settings.
package logz

import (
	"fmt"
	"os"

	C "github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

var InitArgs *kbx.InitArgs

type LogzLogger = C.Logger
type LogzLoggerZ = C.LoggerZ[kbx.Entry]

type LogzConfig = C.LoggerConfig
type LogzOptions = C.LoggerOptionsImpl
type LogzAdvancedOptions = C.LogzAdvancedOptions
type LogzGeneralOptions = kbx.LogzGeneralOptions
type LogzBufferingOptions = kbx.LogzBufferingOptions
type LogzRotatingOptions = kbx.LogzRotatingOptions
type LogzFormatOptions = kbx.LogzFormatOptions
type LogzOutputOptions = kbx.LogzOutputOptions

type LogzJSONFormatter = formatter.JSONFormatter
type LogzTextFormatter = formatter.TextFormatter
type LogzPrettyFormatter = formatter.PrettyFormatter
type LogzFormatter = formatter.Formatter

type LoggerZ = LogzLoggerZ
type Entry = kbx.Entry
type Level = kbx.Level

// defaultLoggerOptions initializes and returns a pointer to a LogzOptions struct
// with default configuration values for logging.
func defaultLoggerOptions() *LogzOptions {
	opts := &LogzOptions{
		LoggerConfig: &LogzConfig{
			ID: kbx.LoggerArgs.ID,
			LogzGeneralOptions: &LogzGeneralOptions{
				Prefix: "",
			},
			LogzFormatOptions: &LogzFormatOptions{
				Output:   os.Stdout,
				Level:    kbx.LevelInfo,
				MinLevel: kbx.LevelInfo,
				MaxLevel: kbx.LevelFatal,
			},
			LogzOutputOptions:    &LogzOutputOptions{},
			LogzRotatingOptions:  &LogzRotatingOptions{},
			LogzBufferingOptions: &LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &LogzAdvancedOptions{},
	}
	opts.Formatter = formatter.NewTextFormatter(false)
	return opts
}

// defaultLogger creates a default logger configured for global use.
func defaultLogger() *LogzLogger {
	return C.NewLogger(
		"",
		defaultLoggerOptions(),
		false,
	)
}

// Logger is the global default logger instance.
var Logger = defaultLogger()

// loggerZ is the global default logger with field support.
var loggerZ *LogzLoggerZ

// NewEntry creates a new log entry with the specified level.
func NewEntry(level Level) (Entry, error) {
	return C.NewEntryImpl(level)
}

// NewGlobalLogger creates a new global logger with the specified prefix.
func NewGlobalLogger(prefix string) *LogzLogger {
	return C.NewLogger(
		prefix,
		defaultLoggerOptions(),
		false,
	)
}

// NewLogger creates a new logger with the specified prefix.
func NewLogger(prefix string) *LoggerZ {
	return C.NewLoggerZ[Entry](
		prefix,
		defaultLoggerOptions(),
		false,
	)
}

// NewLoggerZ creates a new LoggerZ with the given prefix, options, and default settings.
func NewLoggerZ(prefix string, opts *LogzOptions, withDefaults bool) *LogzLoggerZ {
	return C.NewLoggerZ[Entry](prefix, opts, withDefaults)
}

// GetLogger returns the global logger instance, initializing it if necessary.
func GetLogger(prefix string) *LogzLogger {
	if Logger == nil {
		Logger = defaultLogger()
	}
	return Logger
}

// GetLoggerZ returns the global LoggerZ instance, initializing it if necessary.
func GetLoggerZ(prefix string) *LogzLoggerZ {
	if loggerZ == nil {
		loggerZ = NewLoggerZ(prefix, nil, false)
	}
	return loggerZ
}

// Log is the simplest global logging function.
// Accepts a level as string and variadic messages.
func Log(level string, msg ...any) error {
	if Logger == nil {
		return nil
	}
	lvl := kbx.ParseLevel(level)
	return Logger.Log(lvl, msg...)
}

// LogAny is a variant that accepts any type as message.
func LogAny(level string, msg any) error {
	if Logger == nil {
		return nil
	}
	lvl := kbx.ParseLevel(level)
	return Logger.LogAny(lvl, msg)
}

// SetDebugMode enables or disables debug mode for the global logger.
func SetDebugMode(debug bool) {
	if Logger == nil {
		return
	}
	if debug {
		Logger.SetMinLevel(kbx.LevelDebug)
	} else {
		Logger.SetMinLevel(kbx.LevelInfo)
	}
}

// Debug logs a debug message.
func Debug(msg ...any) {
	Log("debug", msg...)
}

// Notice logs a notice message.
func Notice(msg ...any) {
	Log("notice", msg...)
}

// Info logs an informational message.
func Info(msg ...any) {
	Log("info", msg...)
}

// Success logs a success message.
func Success(msg ...any) {
	Log("success", msg...)
}

// Warn logs a warning.
func Warn(msg ...any) {
	Log("warn", msg...)
}

// Error logs an error and returns error.
func Error(msg ...any) error {
	return Log("error", msg...)
}

// Fatal logs a fatal message and exits the program with exit code 1.
func Fatal(msg ...any) {
	Log("fatal", msg...)
	os.Exit(1)
}

func Trace(msg ...any) {
	Log("trace", msg...)
}

func Critical(msg ...any) {
	Log("critical", msg...)
}

func Answer(msg ...any) {
	Log("answer", msg...)
}

func Alert(msg ...any) {
	Log("alert", msg...)
}

func Bug(msg ...any) {
	Log("bug", msg...)
}

func Panic(msg ...any) {
	Log("panic", msg...)
}

func Println(msg ...any) {
	Log("info", fmt.Sprintf("%s", msg...))
}

func Printf(format string, args ...any) {
	Log("info", fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
	Log("debug", fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	Log("info", fmt.Sprintf(format, args...))
}
func Noticef(format string, args ...any) {
	Log("notice", fmt.Sprintf(format, args...))
}

func Successf(format string, args ...any) {
	Log("success", fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	Log("warn", fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) error {
	return Log("error", fmt.Sprintf(format, args...))
}
func Fatalf(format string, args ...any) {
	Log("fatal", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Tracef(format string, args ...any) {
	Log("trace", fmt.Sprintf(format, args...))
}

func Criticalf(format string, args ...any) {
	Log("critical", fmt.Sprintf(format, args...))
}

func Answerf(format string, args ...any) {
	Log("answer", fmt.Sprintf(format, args...))
}

func Alertf(format string, args ...any) {
	Log("alert", fmt.Sprintf(format, args...))
}

func Bugf(format string, args ...any) {
	Log("bug", fmt.Sprintf(format, args...))
}

func Panicf(format string, args ...any) {
	Log("panic", fmt.Sprintf(format, args...))
}

// SetGlobalLogger allows setting a custom global logger instance.
func SetGlobalLogger(logger *LogzLogger) {
	Logger = logger
}

// SetGlobalLoggerZ allows setting a custom global LoggerZ instance.
func SetGlobalLoggerZ(logger *LogzLoggerZ) {
	loggerZ = logger
}

func init() {
	kbx.ParseLoggerArgs(
		"info",
		"notice",
		"fatal",
		"stdout",
	)
	InitArgs = kbx.LoggerArgs
}

func NewLogzFormatter(args *LogzFormatOptions, format string) LogzFormatter {
	switch format {
	case "json":
		return formatter.NewJSONFormatter(true)
	case "pretty":
		return formatter.NewPrettyFormatter()
	default:
		return formatter.NewTextFormatter(true)
	}
}
