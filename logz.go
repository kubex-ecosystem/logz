// Package logz provides a global logging utility with configurable settings.
package logz

import (
	"fmt"
	"io"
	"os"

	// "strings"

	"github.com/kubex-ecosystem/logz/interfaces"
	C "github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/kubex-ecosystem/logz/internal/writer"
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
type EntryImpl = C.Entry
type Entry = kbx.Entry
type Level = kbx.Level

type Writer = writer.Writer
type LogzWriter = writer.LogzWriter
type LogzIOWriter = writer.IOWriter
type LogzMultiWriter = writer.MultiWriter
type LogzEntry = kbx.LogzEntry

type LogzHooks[T any] = interfaces.LHook[T]

func ParseLevel(level string) Level {
	return kbx.ParseLevel(level)
}

func ParseWriter(output string) io.Writer {
	return writer.ParseWriter(output)
}

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
				Level:    kbx.LevelDebug,
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

func defaultLoggerZ() *LogzLoggerZ {
	return C.NewLoggerZ[Entry](
		"",
		defaultLoggerOptions(),
		false,
	)
}

// Logger is the global default logger instance.
var Logger = defaultLogger()

// LoggerLogz is the global default logger with field support.
var LoggerLogz = defaultLoggerZ()

// NewEntry creates a new log entry with the specified level.
func NewEntry(level Level) Entry {
	entry, err := C.NewKbxEntry(level)
	if err != nil {
		// Handle error by returning a default entry with level Info
		defaultEntry, _ := C.NewKbxEntry(kbx.LevelInfo)
		return defaultEntry
	}
	return entry
}

// NewLogzEntry creates a new log entry with the specified level.
func NewLogzEntry(level Level) kbx.LogzEntry {
	return C.NewLogzEntry(level)
}

// NewEntryStrict creates a new log entry with the specified level. (returns error on failure)
func NewEntryStrict(level Level) (Entry, error) {
	return C.NewKbxEntry(level)
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
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
	return LoggerLogz
}

func SetLogzConfig(opts *LogzConfig) {
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
	if opts == nil {
		opts = defaultLoggerOptions().LoggerConfig
	}
	lgrArgs := kbx.ParseLoggerArgs(
		opts.Level.String(),
		opts.MinLevel.String(),
		opts.MaxLevel.String(),
		"",
	)

	cfg := LoggerLogz.GetConfig()

	lgrArgs.ShowColor = opts.ShowColor
	lgrArgs.ShowIcons = opts.ShowIcons
	lgrArgs.ShowTraceID = opts.ShowTraceID
	lgrArgs.ShowFields = opts.ShowFields
	lgrArgs.ShowStack = opts.ShowStack

	lgrArgs.ID = opts.ID
	lgrArgs.LogzGeneralOptions = opts.LogzGeneralOptions
	lgrArgs.LogzFormatOptions = opts.LogzFormatOptions
	lgrArgs.LogzOutputOptions = opts.LogzOutputOptions
	lgrArgs.LogzRotatingOptions = opts.LogzRotatingOptions
	lgrArgs.LogzBufferingOptions = opts.LogzBufferingOptions

	cfg.LoggerConfig = lgrArgs

	LoggerLogz.SetConfig(cfg.LoggerConfig)
}

// Log is the simplest global logging function.
// Accepts a level as string and variadic messages.
func Log(level string, msg ...any) error {
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
	lvl := kbx.ParseLevel(level)
	if lvl.Severity() >= 40 {
		LoggerLogz.Log(lvl, msg...)
		return fmt.Errorf("%v", msg...)
	}
	if LoggerLogz.Enabled(lvl) {
		return LoggerLogz.Log(lvl, msg...)
	}
	return nil
}

// LogAny is a variant that accepts any type as message.
func LogAny(level string, msg any) error {
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
	lvl := kbx.ParseLevel(level)
	if lvl.Severity() >= 40 {
		LoggerLogz.LogAny(lvl, msg)
		return fmt.Errorf("%v", msg)
	}
	if LoggerLogz.Enabled(lvl) {
		return LoggerLogz.LogAny(lvl, msg)
	}
	return nil
}

// SetDebugMode enables or disables debug mode for the global logger.
func SetDebugMode(debug bool) {
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
	if debug {
		LoggerLogz.SetMinLevel(kbx.LevelDebug)
	} else {
		LoggerLogz.SetMinLevel(kbx.LevelInfo)
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
	m := fmt.Sprintln(msg...)
	if len(msg) > 1 && msg[len(msg)-1] == "%log=true%" {
		Log("println", fmt.Sprintf("%s", msg...))
	}
	fmt.Print(m)
}

func Sprintln(msg ...any) string {
	m := fmt.Sprintln(msg...)
	if len(msg) > 1 && msg[len(msg)-1] == "%log=true%" {
		Log("sprintln", m)
	}
	return m
}

func Fprintf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("fprintf", m)
	}
	fmt.Print(m)
}

func Printf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("printf", m)
	}
	fmt.Print(m)
}

func Errorf(format string, args ...any) error {
	m := fmt.Errorf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("error", m)
	}
	return m
}

// func Printf(format string, args ...any) string {
// 	m := fmt.Sprintf(format, args...)
// 	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
// 		Log("printf", m)
// 	}
// 	return m
// }

func Sprintf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("sprintf", m)
	}
	return m
}

func Fatalf(format string, args ...any) {
	Log("fatal", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Tracef(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("trace", m)
}

func Criticalf(format string, args ...any) {
	Log("critical", fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("debug", m)
}

func Infof(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("info", m)
}

func Noticef(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("notice", m)
}

func Successf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("success", m)
}

func Warnf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("warn", m)
}
func Answerf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("answer", m)
}

func Alertf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("alert", m)
}

func Bugf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	Log("bug", m)
}

func Sdebugf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("debug", m)
	}
	return m
}

func Sinfof(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("info", m)
	}
	return m
}

func Snoticef(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("notice", m)
	}
	return m
}

func Ssuccessf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("success", m)
	}
	return m
}

func Swarnf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("warn", m)
	}
	return m
}

func Sanswerf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("answer", m)
	}
	return m
}

func Salertf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("alert", m)
	}
	return m
}

func Sbugf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		Log("bug", m)
	}
	return m
}

func Panicf(format string, args ...any) {
	Log("panic", fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

// SetGlobalLogger allows setting a custom global logger instance.
func SetGlobalLogger(logger *LogzLogger) {
	Logger = logger
}

// SetGlobalLoggerZ allows setting a custom global LoggerZ instance.
func SetGlobalLoggerZ(logger *LogzLoggerZ) {
	LoggerLogz = logger
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
		return formatter.NewPrettyFormatter(true)
	default:
		return formatter.NewTextFormatter(true)
	}
}

func NewLogzWriter(output string, w io.Writer) LogzWriter {
	if w == nil {
		w = writer.ParseWriter(output)
	}
	return writer.NewLogzWriter(w)
}

func NewLogzMultiWriter(outputs ...writer.Writer) LogzWriter {
	return writer.NewMultiWriter(outputs...)
}

func NewLogzIOWriter(w io.Writer) LogzWriter {
	if w == nil {
		w = os.Stdout
	}
	if wrt, ok := w.(writer.LogzWriter); ok {
		return writer.NewDynamicWriter(wrt)
	}
	return writer.NewDynamicWriter(writer.NewLogzWriter(w))
}
