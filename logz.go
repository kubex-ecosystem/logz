// Package logz provides a global logging utility with configurable settings.
package logz

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	// "strings"
	"github.com/google/uuid"
	C "github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/events"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/kubex-ecosystem/logz/internal/writer"
)

// Level é um alias para kbx.Level
type Level = kbx.Level

// Entry é a interface de entrada de log genérica
type Entry = kbx.Entry

// EntryImpl é uma implementação concreta de entrada de log
type EntryImpl = C.Entry

// InitArgs é usado para inicializar o logger.
var InitArgs *kbx.InitArgs

// Aliases para facilitar o uso pelo desenvolvedor, evitar declarações redundantes,
// reaproveitar tipos já existentes e manter retrocompatibilidade, convenção e boas práticas.

// LogzLogger representa o Logger padrão do Logz (Core).
type LogzLogger = C.Logger

// LogzLoggerZ representa o Logger completo do Logz (Core + KBX).
type LogzLoggerZ = C.LoggerZ[kbx.Entry] //nolint

// Configs
type LogzConfig = C.LoggerConfig

// LogzOptions
type LogzOptions = C.LoggerOptionsImpl

// LogzAdvancedOptions
type LogzAdvancedOptions = C.LogzAdvancedOptions

// LogzGeneralOptions
type LogzGeneralOptions = kbx.LogzGeneralOptions

// LogzBufferingOptions
type LogzBufferingOptions = kbx.LogzBufferingOptions

// LogzRotatingOptions
type LogzRotatingOptions = kbx.LogzRotatingOptions

// LogzFormatOptions
type LogzFormatOptions = kbx.LogzFormatOptions

// LogzOutputOptions
type LogzOutputOptions = kbx.LogzOutputOptions

// LogzJSONFormatter
type LogzJSONFormatter = formatter.JSONFormatter

// LogzTextFormatter
type LogzTextFormatter = formatter.TextFormatter

// LogzPrettyFormatter
type LogzPrettyFormatter = formatter.PrettyFormatter

// LogzHUDFormatter
type LogzHUDFormatter = formatter.HUDFormatter

// LogzFormatter
type LogzFormatter = formatter.Formatter

// LoggerZ
type LoggerZ = LogzLoggerZ

// LogzEntryImpl
type LogzEntryImpl = C.Entry

// LogzLevel
type LogzLevel = kbx.Level

// Writer
type Writer = writer.Writer

// LogzWriter
type LogzWriter = writer.LogzWriter

// LogzIOWriter
type LogzIOWriter = writer.IOWriter

// LogzMultiWriter
type LogzMultiWriter = writer.MultiWriter

// LogzEntry
type LogzEntry = kbx.LogzEntry

// LogzHooks
type LogzHooks[T any] = events.LHook[T]

// NewLogzOptions is a wrapper for C.NewLoggerOptions
// Creates a new LogzOptions struct.
// If withDefaults is true, it uses default values for all options.
// If withDefaults is false, it uses the values from kbx.LoggerArgs.
func NewLogzOptions(withDefaults bool) *LogzOptions {
	if withDefaults {
		return defaultLoggerOptions()
	}
	if kbx.LoggerArgs == nil {
		//kbx.NewConfig()
		kbx.LoggerArgs = &kbx.InitArgs{
			ID:                   uuid.New(),
			Metadata:             make(map[string]string),
			LogzGeneralOptions:   &LogzGeneralOptions{},
			LogzFormatOptions:    &LogzFormatOptions{},
			LogzOutputOptions:    &LogzOutputOptions{},
			LogzRotatingOptions:  &LogzRotatingOptions{},
			LogzBufferingOptions: &LogzBufferingOptions{},
		}
	}
	opts := C.NewLoggerOptions(kbx.LoggerArgs)

	opts.Level = ParseLevel(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_LEVEL", kbx.DefaultLogLevel))
	opts.MinLevel = ParseLevel(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_MIN_LEVEL", kbx.DefaultLogMinLevel))
	opts.MaxLevel = ParseLevel(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_MAX_LEVEL", kbx.DefaultLogMaxLevel))
	opts.Output = ParseWriter(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_OUTPUT", kbx.DefaultLogOutput))
	opts.ShowColor = kbx.BoolPtr(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_SHOW_COLOR", kbx.DefaultShowColor))
	opts.ShowIcons = kbx.BoolPtr(kbx.GetEnvOrDefaultWithType("LOGZ_LOG_SHOW_ICONS", kbx.DefaultShowIcons))
	opts.ShowTraceID = kbx.GetEnvOrDefaultWithType("LOGZ_LOG_SHOW_TRACE_ID", kbx.DefaultShowTraceID)
	opts.ShowFields = kbx.GetEnvOrDefaultWithType("LOGZ_LOG_SHOW_FIELDS", kbx.DefaultShowFields)
	opts.ShowStack = kbx.GetEnvOrDefaultWithType("LOGZ_LOG_SHOW_STACK", kbx.DefaultShowStack)

	return opts
}

// ParseLevel converts a string representation of a log level to a Level enum.
func ParseLevel(level string) Level {
	return kbx.ParseLevel(level)
}

// ParseWriter converts a string representation of a log output to an io.Writer.
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
				Output:   ParseWriter(kbx.DefaultLogOutput),
				Level:    ParseLevel(kbx.DefaultLogLevel),
				MinLevel: ParseLevel(kbx.DefaultLogMinLevel),
				MaxLevel: ParseLevel(kbx.DefaultLogMaxLevel),
			},
			LogzOutputOptions:    &LogzOutputOptions{},
			LogzRotatingOptions:  &LogzRotatingOptions{},
			LogzBufferingOptions: &LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &LogzAdvancedOptions{},
	}
	opts.Formatter = formatter.ParseFormatter(kbx.DefaultLogFormat, kbx.DefaultShowColor)
	return opts
}

// defaultLogger creates a default logger configured for global use.
func defaultLogger() *LogzLogger {
	l := logger.Load()
	if l != nil {
		return l
	}
	l = C.NewLogger(
		"",
		defaultLoggerOptions(),
		false,
	)
	logger.Store(l)
	return l
}

// defaultLoggerZ creates a default logger with field support configured for global use.
func defaultLoggerZ() *LogzLoggerZ {
	l := loggerLogz.Load()
	if l != nil {
		return l
	}
	l = C.NewLoggerZ[Entry](
		"",
		defaultLoggerOptions(),
		false,
	)
	loggerLogz.Store(l)
	return l
}

// Logger is the global default logger instance.
var logger atomic.Pointer[LogzLogger]

// Logger is a convenience field for accessing the global default logger.
var Logger *LogzLogger

// LoggerLogz is the global default logger with field support.
var loggerLogz atomic.Pointer[LogzLoggerZ]

// LoggerLogz is the global representation of the default logger (atomic pointer).
var LoggerLogz *LogzLoggerZ

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

// SetLogger provide an interface to use the go standard library log.
// It requires to set the output and flags for the logger.
func SetLogger(l *log.Logger, prefix string, opts *LogzOptions, withDefaults bool) {
	if l == nil {
		return
	}

	if LoggerLogz == nil {
		LoggerLogz = NewLoggerZ(prefix, opts, withDefaults)
	}

	ll := GetLogger(prefix)

	ll.Logger.SetOutput(l.Writer())
	ll.Logger.SetFlags(l.Flags())

	// Ele já faz o store no atomic logger
	// No entanto, a manipulação do ponteiro atomic.Pointer
	// não ocorre aqui, só ocorre no constructor do stdlog.
	// O que ocorre é a cópia de parâmetros do log standard
	// para o logz. É o suficiente por nós usarmos nossos próprios writers
	// o que permite que a propagação das mensagens ocorra, já que ele é um ponteiro
	// e está linkado com o writer do log standard.
	ll.Logger = l
}

// SetLogzConfig permite alterar a configuração do logger.
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
		// O método Log(level, msg...) já está ciente de que
		// se o nível for de erro (>= 40), ele deve retornar um erro.
		return LoggerLogz.Log(lvl, msg...)
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
		return LoggerLogz.LogAny(lvl, msg)
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
		LoggerLogz.SetMinLevel(kbx.LevelNotice)
	} else {
		LoggerLogz.SetMinLevel(kbx.LevelInfo)
	}
}

// Debug logs a debug message.
func Debug(msg ...any) {
	_ = Log("debug", msg...)
}

// Notice logs a notice message.
func Notice(msg ...any) {
	_ = Log("notice", msg...)
}

// Info logs an informational message.
func Info(msg ...any) {
	_ = Log("info", msg...)
}

// Success logs a success message.
func Success(msg ...any) {
	_ = Log("success", msg...)
}

// Warn logs a warning.
func Warn(msg ...any) {
	_ = Log("warn", msg...)
}

// Error logs an error and returns error.
func Error(msg ...any) error {
	return Log("error", msg...)
}

// Fatal logs a fatal message and exits the program with exit code 1.
func Fatal(msg ...any) {
	_ = Log("fatal", msg...)
	os.Exit(1)
}

// Trace logs a trace message.
func Trace(msg ...any) {
	// TODO: inserir o uso do runtime.Trace para adquirir métricas nativas do Go.
	_ = Log("trace", msg...)
}

// Critical logs a critical message.
func Critical(msg ...any) {
	_ = Log("critical", msg...)
}

// Answer logs an answer message.
func Answer(msg ...any) {
	_ = Log("answer", msg...)
}

// Alert logs an alert message.
func Alert(msg ...any) {
	_ = Log("alert", msg...)
}

// Bug logs a bug message.
func Bug(msg ...any) {
	_ = Log("bug", msg...)
}

// Panic logs a panic message and panics.
func Panic(msg ...any) {
	panic(Log("panic", msg...))
}

// Println logs a println message.
func Println(msg ...any) {
	m := fmt.Sprintln(msg...)
	if len(msg) > 1 && msg[len(msg)-1] == "%log=true%" {
		_ = Log("println", fmt.Sprintf("%s", msg...))
	}
	fmt.Print(m)
}

// Sprintln logs a sprintln message.
func Sprintln(msg ...any) string {
	m := fmt.Sprintln(msg...)
	if len(msg) > 1 && msg[len(msg)-1] == "%log=true%" {
		_ = Log("sprintln", m)
	}
	return m
}

// Fprintf logs a fprintf message.
func Fprintf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("fprintf", m)
	}
	fmt.Print(m)
}

// Printf logs a printf message.
func Printf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("printf", m)
	}
	fmt.Print(m)
}

// Errorf logs an errorf message.
func Errorf(format string, args ...any) error {
	m := fmt.Errorf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("error", m)
	}
	return m
}

// // Sprintf logs a sprintf message. (not used)
// func Sprintf(format string, args ...any) string {
// 	m := fmt.Sprintf(format, args...)
// 	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
// 		Log("printf", m)
// 	}
// 	return m
// }

// Sprintf logs a sprintf message.
func Sprintf(format string, args ...any) string {
	// m := fmt.Sprintf(format, args...)
	// if len(args) > 1 && args[len(args)-1] == "%log=true%" {
	// 	_ = Log("sprintf", m)
	// }
	// return m
	return fmt.Sprintf(format, args...)
}

// Fatalf logs a fatalf message.
func Fatalf(format string, args ...any) {
	_ = Log("fatal", fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Tracef logs a tracef message.
func Tracef(format string, args ...any) {
	// TODO: inserir o uso do runtime.Trace para adquirir métricas nativas do Go.
	_ = Log("trace", fmt.Sprintf(format, args...))
}

// Criticalf logs a criticalf message.
func Criticalf(format string, args ...any) {
	_ = Log("critical", fmt.Sprintf(format, args...))
}

// Debugf logs a debugf message.
func Debugf(format string, args ...any) {
	_ = Log("debug", fmt.Sprintf(format, args...))
}

// Infof logs an infof message.
func Infof(format string, args ...any) {
	_ = Log("info", fmt.Sprintf(format, args...))
}

// Noticef logs a noticef message.
func Noticef(format string, args ...any) {
	_ = Log("notice", fmt.Sprintf(format, args...))
}

// Successf logs a successf message.
func Successf(format string, args ...any) {
	_ = Log("success", fmt.Sprintf(format, args...))
}

// Warnf logs a warnf message.
func Warnf(format string, args ...any) {
	_ = Log("warn", fmt.Sprintf(format, args...))
}

// Answerf logs an answerf message.
func Answerf(format string, args ...any) {
	_ = Log("answer", fmt.Sprintf(format, args...))
}

// Alertf logs an alertf message.
func Alertf(format string, args ...any) {
	_ = Log("alert", fmt.Sprintf(format, args...))
}

// Bugf logs a bugf message.
func Bugf(format string, args ...any) {
	_ = Log("bug", fmt.Sprintf(format, args...))
}

// Sdebugf logs a sdebugf message.
func Sdebugf(format string, args ...any) string {
	// Aqui há uma opção interna de habilitar a saída do log
	// Ela é lida internamente por esta função e não deve ser usada externamente.
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("debug", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("debugf", m)
	return m
}

// Sinfof logs a sinfof message.
func Sinfof(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("info", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("infof", m)
	return m
}

// Snoticef logs a snoticef message.
func Snoticef(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("notice", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("noticef", m)
	return m
}

// Ssuccessf logs a ssuccessf message.
func Ssuccessf(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("success", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("successf", m)
	return m
}

// Serrorf logs a serrorf message.
func Serrorf(format string, args ...any) error {
	m := fmt.Errorf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("error", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("errorf", m)
	return m
}

// Swarnf logs a swarnf message.
func Swarnf(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("warn", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("warnf", m)
	return m
}

// Sanswerf logs a sanswerf message.
func Sanswerf(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("answer", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("answerf", m)
	return m
}

// Salertf logs a salertf message.
func Salertf(format string, args ...any) string {
	m := fmt.Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("alert", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("alertf", m)
	return m
}

// Sbugf logs a sbugf message.
func Sbugf(format string, args ...any) string {
	m := Sprintf(format, args...)
	if len(args) > 1 && args[len(args)-1] == "%log=true%" {
		_ = Log("bug", m)
	}
	// Os "pseudo" níveis atuam somente como um mecanismo para
	// evitar a saída do log no output diretamente, e garantir o
	// retorno da string formatada pelo Logz (fmt.Sprintf).
	_ = Log("bugf", m)
	return m
}

// Panicf logs a panicf message.
func Panicf(format string, args ...any) {
	// O Log("panic", ...) sempre retorna o error, então passamos o resultado de Log
	// para o panic(any) que é o tipo esperado pelo Panicf.
	panic(Log("panic", Sprintf(format, args...)))
}

// SetGlobalLogger allows setting a custom global logger instance.
func SetGlobalLogger(logger *LogzLogger) {
	Logger = logger
}

// SetGlobalLoggerZ allows setting a custom global LoggerZ instance.
func SetGlobalLoggerZ(logger *LogzLoggerZ) {
	LoggerLogz = logger
}

// init initializes the logger.
func init() {
	if InitArgs == nil || kbx.LoggerArgs == nil {
		kbx.ParseLoggerArgs(
			kbx.DefaultLogLevel,
			kbx.DefaultLogMinLevel,
			kbx.DefaultLogMaxLevel,
			kbx.DefaultLogOutput,
		)
		InitArgs = kbx.LoggerArgs
	}
	if Logger == nil {
		Logger = defaultLogger()
	}
	if LoggerLogz == nil {
		LoggerLogz = defaultLoggerZ()
	}
}

// NewLogzFormatter creates a new LogzFormatter.
func NewLogzFormatter(args *LogzFormatOptions, format string, pretty bool) LogzFormatter {
	switch format {
	case "json":
		return formatter.NewJSONFormatter(pretty)
	case "pretty":
		return formatter.NewPrettyFormatter(pretty)
	default:
		return formatter.NewTextFormatter(pretty)
	}
}

// NewLogzWriter creates a new LogzWriter.
func NewLogzWriter(output string, w io.Writer) LogzWriter {
	if w == nil {
		w = writer.ParseWriter(output)
	}
	return writer.NewLogzWriter(w)
}

// NewLogzMultiWriter creates a new LogzMultiWriter.
func NewLogzMultiWriter(outputs ...writer.Writer) LogzWriter {
	return writer.NewMultiWriter(outputs...)
}

// NewLogzIOWriter creates a new LogzIOWriter.
func NewLogzIOWriter(w io.Writer) LogzWriter {
	if w == nil {
		w = os.Stdout
	}
	if wrt, ok := w.(writer.LogzWriter); ok {
		return writer.NewDynamicWriter(wrt)
	}
	return writer.NewDynamicWriter(writer.NewLogzWriter(w))
}
