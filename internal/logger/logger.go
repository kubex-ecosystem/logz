package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LogMode string
type LogFormat string

const (
	JSON LogFormat = "json"
	TEXT LogFormat = "text"
	YAML LogFormat = "yaml"
	XML  LogFormat = "xml"
	RAW  LogFormat = "raw"
)

const (
	ModeService    LogMode = "service"    // Indicates that the logger is being used by a detached process
	ModeStandalone LogMode = "standalone" // Indicates that the logger is being used locally (e.g., CLI)
)

var logLevels = map[LogLevel]int{
	DEBUG: 1,
	INFO:  2,
	WARN:  3,
	ERROR: 4,
	FATAL: 5,
}

// LogzCoreImpl represents a logger with configuration and metadata.
type LogzCoreImpl struct {
	level    LogLevel
	writer   LogWriter
	config   Config
	metadata map[string]interface{}
	mode     LogMode // Mode control: service or standalone
	mu       sync.RWMutex
}

// NewLogger creates a new instance of LogzCoreImpl with the provided configuration.
func NewLogger(config Config) *LogzCoreImpl {
	level := INFO // Default log level

	if config == nil {
		// If no config is provided, create a new default config
		// Create a new ConfigManager
		if cfgMgr := NewConfigManager(); cfgMgr != nil {
			configMgr := *cfgMgr
			config = configMgr.GetConfig()
		} else {
			log.Println("Error creating ConfigManager")
			return nil
		}
	}

	// Set the log level from the Config
	level = LogLevel(config.Level()) // Method config.Level() returns the log level as a string

	var out *os.File
	var err error
	// Set the output to stdout if not specified or if the output is invalid
	if strings.ToLower(config.Output()) == "stdout" || config.Output() == "" || config.Output() == os.Stdout.Name() {
		out = os.Stdout
	} else {
		fmt.Println("Output: ", config.Output())
		// Ensure the log file exists and has the correct permissions
		if _, err = os.Stat(config.Output()); os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(config.Output()), 0755); err != nil {
				log.Printf("Error creating log directory: %v\nRedirecting to stdout...\n", err)
				out = os.Stdout
			} else {
				out, err = os.OpenFile(config.Output(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Printf("Error opening log file: %v\nRedirecting to stdout...\n", err)
					out = os.Stdout
				}
			}
		} else {
			out, err = os.OpenFile(config.Output(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening log file: %v\nRedirecting to stdout...\n", err)
				out = os.Stdout
			}
		}
	}

	// Initialize the formatter (JSON or text)
	var formatter LogFormatter
	if config.Format() == "json" {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{}
	}
	writer := NewDefaultWriter(out, formatter)

	// Read the mode from Config
	mode := config.Mode()
	if mode != ModeService && mode != ModeStandalone {
		mode = ModeStandalone // Default to standalone if not specified
	}

	return &LogzCoreImpl{
		level:    level,
		writer:   writer,
		config:   config,
		metadata: make(map[string]interface{}),
		mode:     mode,
	}
}

// SetMetadata sets a metadata key-value pair for the LogzCoreImpl.
func (l *LogzCoreImpl) SetMetadata(key string, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.metadata[key] = value
}

// shouldLog checks if the log level should be logged.
func (l *LogzCoreImpl) shouldLog(level LogLevel) bool {
	return logLevels[level] >= logLevels[l.level]
}

// log logs a message with the specified level and context.
func (l *LogzCoreImpl) log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	entry := NewLogEntry().
		WithLevel(level).
		WithMessage(msg).
		WithSeverity(logLevels[level])

	// Merge global and local metadata
	finalContext := mergeContext(l.metadata, ctx)
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	// Merge global and local metadata
	finalMetadata := mergeMetadata(l.metadata, ctx)
	for k, v := range finalMetadata {
		entry.AddMetadata(k, v)
	}

	// Write the log using the configured writer
	if err := l.writer.Write(entry); err != nil {
		log.Printf("Error writing log: %v", err)
	}

	// Only in service mode, notify via Notifiers
	if l.mode == ModeService && l.config != nil {
		for _, name := range l.config.NotifierManager().ListNotifiers() {
			if notifier, ok := l.config.NotifierManager().GetNotifier(name); ok {
				if notifier != nil {
					ntf := notifier
					if ntfErr := ntf.Notify(entry); ntfErr != nil {
						log.Printf("Error notifying %s: %v", name, ntfErr)
					}
				}
			}
		}
	}

	// Update metrics in PrometheusManager, if enabled
	if l.mode == ModeService {
		pm := GetPrometheusManager()
		if pm.IsEnabled() {
			pm.IncrementMetric("logs_total", 1)
			pm.IncrementMetric("logs_total_"+string(level), 1)
		}
	}

	// Terminate the process in case of FATAL log
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message with context.
func (l *LogzCoreImpl) Debug(msg string, ctx map[string]interface{}) { l.log(DEBUG, msg, ctx) }

// Info logs an info message with context.
func (l *LogzCoreImpl) Info(msg string, ctx map[string]interface{}) { l.log(INFO, msg, ctx) }

// Warn logs a warning message with context.
func (l *LogzCoreImpl) Warn(msg string, ctx map[string]interface{}) { l.log(WARN, msg, ctx) }

// Error logs an error message with context.
func (l *LogzCoreImpl) Error(msg string, ctx map[string]interface{}) { l.log(ERROR, msg, ctx) }

// FatalC logs a fatal message with context and terminates the process.
func (l *LogzCoreImpl) FatalC(msg string, ctx map[string]interface{}) { l.log(FATAL, msg, ctx) }

func (l *LogzCoreImpl) SetLevel(level LogLevel) { l.level = level }
func (l *LogzCoreImpl) GetLevel() LogLevel      { return l.level }

func (l *LogzCoreImpl) SetWriter(writer LogWriter) { l.writer = writer }
func (l *LogzCoreImpl) GetWriter() LogWriter       { return l.writer }

func (l *LogzCoreImpl) GetMode() LogMode { return l.mode }

func (l *LogzCoreImpl) SetConfig(config Config) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config = config
}
func (l *LogzCoreImpl) GetConfig() Config {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.config
}

// trimFilePath trims the file path to show only the last two segments.
func trimFilePath(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return filePath
}

// mergeContext merges global and local context maps.
func mergeContext(global, local map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}
	return merged
}

// mergeMetadata merges global and local context maps.
func mergeMetadata(global, local map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}
	return merged
}
