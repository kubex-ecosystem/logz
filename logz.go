package logz

import (
	"fmt"
	core "github.com/faelmori/logz/internal/logger"
	logz "github.com/faelmori/logz/logger"
	vs "github.com/faelmori/logz/version"
	"os"
	"sync"
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
type Config = core.Config
type ConfigManager = core.LogzConfigManager
type NotifierManager = core.NotifierManager
type Notifier = core.Notifier
type Logger = logz.LogzLogger
type Writer = core.LogWriter

// initializeLogger initializes the global logger with the given prefix.
func initializeLogger(prefix string) {
	//	once.Do(func() {
	if prefix == "" {
		prefix = pfx
	}
	if logger != nil {
		return
	}
	logger = logz.NewLogger(prefix)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		logger.SetLevel(core.LogLevel(logLevel))
	} else {
		logger.SetLevel(core.INFO)
	}
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat != "" {
		logger.GetConfig().SetFormat(core.LogFormat(logFormat))
	} else {
		logger.GetConfig().SetFormat(core.TEXT)
	}
	logOutput := os.Getenv("LOG_OUTPUT")
	if logOutput != "" {
		logger.GetConfig().SetOutput(logOutput)
	} else {
		logger.GetConfig().SetOutput(os.Stdout.Name())
	}
	//	})
}

// GetLogger returns the global logger instance, initializing it if necessary.
func GetLogger(prefix string) Logger {
	initializeLogger(prefix)

	////mu.RLock()
	//defer mu.RUnlock()
	return logger
}

// NewLogger creates a new logger instance with the given prefix.
func NewLogger(prefix string) Logger {
	return logz.NewLogger(prefix)
}

// SetLogger sets the global logger instance to the provided logger.
func SetLogger(newLogger Logger) {
	//mu.Lock()
	//defer mu.Unlock()
	logger = newLogger
}

// SetPrefix sets the global prefix for the logger.
func SetPrefix(prefix string) {
	//mu.Lock()
	//defer mu.Unlock()
	pfx = prefix
}

// GetPrefix returns the global prefix for the logger.
func GetPrefix() string {
	////mu.RLock()
	//defer mu.RUnlock()
	return pfx
}

// SetLogLevel sets the log level for the global logger.
func SetLogLevel(level LogLevel) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetLevel(level)
	}
}

// GetLogLevel returns the log level of the global logger.
func GetLogLevel() LogLevel {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return core.DEBUG
	}
	return logger.GetLevel()
}

// SetLogWriter sets the log writer for the global logger.
func SetLogWriter(writer Writer) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetWriter(writer)
	}
}

// GetLogWriter returns the log writer of the global logger.
func GetLogWriter() Writer {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.GetWriter()
}

// SetLogConfig sets the configuration for the global logger.
func SetLogConfig(config Config) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetConfig(config)
	}
}

// GetLogConfig returns the configuration of the global logger.
func GetLogConfig() Config {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.GetConfig()
}

// SetMetadata sets a metadata key-value pair for the global logger.
func SetMetadata(key string, value interface{}) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.SetMetadata(key, value)
	}
}

// Trace logs a trace message with the given context.
func Trace(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Trace(msg, ctx)
	}
}

// Notice logs a notice message with the given context.
func Notice(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Notice(msg, ctx)
	}
}

// Success logs a success message with the given context.
func Success(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Success(msg, ctx)
	}
}

// Debug logs a debug message with the given context.
func Debug(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Debug(msg, ctx)
	}
}

// Info logs an info message with the given context.
func Info(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Info(msg, ctx)
	}
}

// Warn logs a warning message with the given context.
func Warn(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Warn(msg, ctx)
	}
}

// Error logs an error message with the given context.
func Error(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.Error(msg, ctx)
	}
}

// FatalC logs a fatal message with the given context and exits the application.
func FatalC(msg string, ctx map[string]interface{}) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger != nil {
		logger.FatalC(msg, ctx)
	}
}

// AddNotifier adds a notifier to the global logger's configuration.
func AddNotifier(name string, notifier Notifier) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.
			GetConfig().
			NotifierManager().
			AddNotifier(name, notifier)
	}
}

// GetNotifier returns the notifier with the given name from the global logger's configuration.
func GetNotifier(name string) (Notifier, bool) {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil, false
	}
	return logger.
		GetConfig().
		NotifierManager().
		GetNotifier(name)
}

// ListNotifiers returns a list of all notifier names in the global logger's configuration.
func ListNotifiers() []string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.
		GetConfig().
		NotifierManager().
		ListNotifiers()
}

// SetLogFormat sets the log format for the global logger.
func SetLogFormat(format LogFormat) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.
			GetConfig().
			SetFormat(format)
	}
}

// GetLogFormat returns the log format of the global logger.
func GetLogFormat() string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return "text"
	}
	return logger.GetConfig().Format()
}

// SetLogOutput sets the log output for the global logger.
func SetLogOutput(output string) {
	//mu.Lock()
	//defer mu.Unlock()
	if logger != nil {
		logger.GetConfig().SetOutput(output)
	}
}

// GetLogOutput returns the log output of the global logger.
func GetLogOutput() string {
	//mu.RLock()
	//defer mu.RUnlock()
	if logger == nil {
		return os.Stdout.Name()
	}
	return logger.GetConfig().Output()
}

// CheckVersion checks the version of the logger.
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

// Version returns the current version of the logger.
func Version() string {
	if versionService == nil {
		versionService = vs.NewVersionService()
	}
	return versionService.GetCurrentVersion()
}
