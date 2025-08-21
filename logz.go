package logz

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rafa-mori/logz/internal/core"
	logz "github.com/rafa-mori/logz/logger"
	vs "github.com/rafa-mori/logz/version"
)

var (
	pfx    = "Logz" // Default prefix
	logger Logger   // Global logger instance
	//mu             sync.RWMutex // Mutex for concurrency control
	once           sync.Once // Ensure single initialization
	versionService vs.Service
)

type LogLevel = core.LogLevel
type LogFormat = core.LogFormat

type Config interface{ core.Config }

// type ConfigManager interface{ core.LogzConfigManager }
type NotifierManager interface{ core.NotifierManager }
type Notifier interface{ core.Notifier }
type Logger interface{ logz.LogzLogger }

type JSONFormatter = core.JSONFormatter
type TextFormatter = core.TextFormatter

type Writer struct{ core.LogWriter[any] }

func (w Writer) Write(p []byte) (n int, err error) {
	var decodedMessage map[string]interface{}
	if jsonErr := json.Unmarshal(p, &decodedMessage); jsonErr == nil {
		entry := core.NewLogEntry().
			WithMessage(decodedMessage["message"].(string)).
			WithLevel(GetLogLevel()).
			AddMetadata("original", decodedMessage)

		writeErr := w.LogWriter.Write(entry)
		if writeErr != nil {
			return 0, writeErr
		}
	} else {
		entry := core.NewLogEntry().
			WithMessage(string(p)).
			WithLevel(GetLogLevel())
		writeErr := w.LogWriter.Write(entry)
		if writeErr != nil {
			return 0, writeErr
		}
	}

	return len(p), nil
}

// SetLogWriter sets the log writer for the global core.
func SetLogWriter(writer interface{}) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		nWriter := core.NewDefaultWriter[any](writer.(Writer), &TextFormatter{})
		logger.SetWriter(nWriter)
	}
}

// GetLogWriter returns the log writer of the global core.
func GetLogWriter() *Writer {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	writer := logger.GetWriter().(logz.LogWriter[any])
	return &Writer{LogWriter: writer}
}

func NewWriter(out *os.File, formatter core.LogFormatter) Writer {
	if out == nil {
		out = os.Stdout
	}
	return Writer{LogWriter: core.NewDefaultWriter[any](out, formatter)}
}

// initializeLogger initializes the global logger with the given prefix.
func initializeLogger(prefix string) {
	//	once.Do(func() {
	if prefix == "" {
		prefix = pfx
	}
	if logger != nil {
		return
	}
	logger = logz.NewLogger(prefix).(Logger)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		logger.SetLevel(core.LogLevel(logLevel))
	} else {
		logger.SetLevel(core.INFO)
	}
	logFormat := os.Getenv("LOG_FORMAT")
	//config := logger.GetConfig().(*core.Config)
	if logFormat != "" {
		logger.SetFormat(core.LogFormat(logFormat))
	} else {
		//logger.GetConfig().SetFormat(core.TEXT)
	}
	logOutput := os.Getenv("LOG_OUTPUT")
	if logOutput != "" {
		//logger.GetConfig().SetOutput(logOutput)
	} else {
		//logger.GetConfig().SetOutput(os.Stdout.Name())
	}
	//	})
}

// GetLogger returns the global core instance, initializing it if necessary.
func GetLogger(prefix string) Logger {
	initializeLogger(prefix)

	////mu.RLock()
	//defer mu.RUnlock()
	return logger
}

// NewLogger creates a new core instance with the given prefix.
func NewLogger(prefix string) Logger {
	return logz.NewLogger(prefix)
}

// SetLogger sets the global core instance to the provided core.
func SetLogger(newLogger Logger) {
	//mu.Lock()
	//defer mu.Unlock()
	logger = newLogger
}

// SetPrefix sets the global prefix for the core.
func SetPrefix(prefix string) {
	//mu.Lock()
	//defer mu.Unlock()
	pfx = prefix
}

// GetPrefix returns the global prefix for the core.
func GetPrefix() string {
	////mu.RLock()
	//defer mu.RUnlock()
	return pfx
}

// SetLogLevel sets the log level for the global core.
func SetLogLevel(level LogLevel) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetLevel(level)
	}
}

// GetLogLevel returns the log level of the global core.
func GetLogLevel() LogLevel {
	if logger == nil {
		return core.DEBUG
	}
	return LogLevel(logger.GetLevel().(string))
}

// SetLogConfig sets the configuration for the global core.
func SetLogConfig(config Config) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetConfig(config)
	}
}

// GetLogConfig returns the configuration of the global core.
func GetLogConfig() Config {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.GetConfig().(core.Config)
}

// SetMetadata sets a metadata key-value pair for the global core.
func SetMetadata(key string, value interface{}) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetMetadata(key, value)
	}
}

// Trace logs a trace message with the given context.
func TraceCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.TraceCtx(msg, ctx)
	}
}

// Notice logs a notice message with the given context.
func NoticeCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.NoticeCtx(msg, ctx)
	}
}

// Success logs a success message with the given context.
func SuccessCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.SuccessCtx(msg, ctx)
	}
}

// Debug logs a debug message with the given context.
func DebugCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.DebugCtx(msg, ctx)
	}
}

// InfoCtx logs an info message with the given context.
func InfoCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.InfoCtx(msg, ctx)
	}
}

// Warn logs a warning message with the given context.
func WarnCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.WarnCtx(msg, ctx)
	}
}

// Error logs an error message with the given context.
func ErrorCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.ErrorCtx(msg, ctx)
	}
}

// FatalC logs a fatal message with the given context and exits the application.
func FatalCtx(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.FatalCtx(msg, ctx)
	}
}

// AddNotifier adds a notifier to the global core's configuration.
func AddNotifier(name string, notifier Notifier) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.
			GetConfig() /*.
			NotifierManager().
			AddNotifier(name, notifier)*/
	}
}

// GetNotifier returns the notifier with the given name from the global core's configuration.
func GetNotifier(name string) (interface{}, bool) { //(Notifier, bool) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil, false
	}
	return logger.GetConfig(), true
}

// ListNotifiers returns a list of all notifier names in the global core's configuration.
func ListNotifiers() []string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return nil /*logger.
	GetConfig().
	NotifierManager().
	ListNotifiers()*/
}

// SetLogFormat sets the log format for the global core.
func SetLogFormat(format LogFormat) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		/*logger.
		GetConfig().
		SetFormat(format)*/
	}
}

// GetLogFormat returns the log format of the global core.
func GetLogFormat() string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return "text"
	}
	cfg := logger.GetConfig()
	if cfg == nil {
		return "text"
	} else {
		return "json"
	}

	//return logger.GetConfig().Format()
}

// SetLogOutput sets the log output for the global core.
func SetLogOutput(output string) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		cfg := logger.GetConfig()
		ccc := cfg.(*core.Config)
		cc := *ccc
		cc.SetOutput(output)
		logger.SetConfig(ccc)

	}
}

// GetLogOutput returns the log output of the global core.
func GetLogOutput() string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return os.Stdout.Name()
	}
	cfg := logger.GetConfig()
	if cfg == nil {
		return os.Stdout.Name()
	} else {
		return cfg.(core.Config).Output()
	}
}

// CheckVersion checks the version of the core.
func CheckVersion() string {
	if versionService == nil {
		versionService = vs.NewVersionService()
	}
	if isLatest, err := versionService.IsLatestVersion(); err != nil {
		return "error checking version"
	} else {
		if isLatest {
			return "latest version"
		}
	}
	if latestVersion, err := versionService.GetLatestVersion(); err != nil {
		return "error getting latest version"
	} else {
		return fmt.Sprintf("latest version: %s\nYou are using version: %s", latestVersion, versionService.GetCurrentVersion())
	}
}

// Version returns the current version of the core.
func Version() string {
	if versionService == nil {
		versionService = vs.NewVersionService()
	}
	return versionService.GetCurrentVersion()
}

// Info returns the log output of the global core.
func Info(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.InfoCtx(fmt.Sprint(args...), nil)
}

// Debug returns the log output of the global core.
func Debug(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.DebugCtx(fmt.Sprint(args...), nil)
}

// Warn returns the log output of the global core.
func Warn(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.WarnCtx(fmt.Sprint(args...), nil)
}

// Error returns the log output of the global core.
func Error(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.ErrorCtx(fmt.Sprint(args...), nil)
}

// Fatal returns the log output of the global core.
func Fatal(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.FatalCtx(fmt.Sprint(args...), nil)
}

// Trace returns the log output of the global core.
func Trace(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.TraceCtx(fmt.Sprint(args...), nil)
}

// Notice returns the log output of the global core.
func Notice(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.NoticeCtx(fmt.Sprint(args...), nil)
}

// Success returns the log output of the global core.
func Success(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.SuccessCtx(fmt.Sprint(args...), nil)
}

// Panic returns the log output of the global core.
func Panic(args ...any) {
	if logger == nil {
		logger = logz.NewLogger(pfx)
	}
	logger.FatalCtx(fmt.Sprint(args...), nil)
}
