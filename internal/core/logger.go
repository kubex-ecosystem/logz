package core

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"

	//"io"
	"os"
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
	ModeService    LogMode = "service"    // Indicates that the core is being used by a detached process
	ModeStandalone LogMode = "standalone" // Indicates that the core is being used locally (e.g., CLI)
)

var logLevels = map[LogLevel]int{
	DEBUG:   1,
	TRACE:   2,
	INFO:    3,
	NOTICE:  4,
	SUCCESS: 5,
	WARN:    6,
	ERROR:   7,
	FATAL:   8,
	SILENT:  9,
	ANSWER:  10,
}

// LogzCoreImpl represents a core with configuration and VMetadata.
type LogzCoreImpl struct {
	// LogzLogger is a constraint to implement this interface
	LogzLogger

	// Logger is a promoted global Go Logger
	*log.Logger

	out       io.Writer                   // destination for output
	prefix    atomic.Pointer[string]      // prefix on each line to identify the logger (but see Lmsgprefix)
	prefixX   atomic.Pointer[*LogzLogger] // prefix on each line to identify the logger (but see Lmsgprefix)
	flag      atomic.Int32                // properties
	isDiscard atomic.Bool

	VLevel    LogLevel
	VWriter   LogWriter[any]
	VConfig   Config
	VMetadata map[string]interface{}
	VMode     LogMode // Mode control: service or standalone
	Mu        sync.RWMutex
}

// NewLogger creates a new instance of LogzCoreImpl with the provided configuration.
func NewLogger(prefix string) LogzLogger {
	return NewLoggerImpl(prefix)
}

func NewLoggerImpl(prefix string) *LogzCoreImpl {
	level := INFO // Default log VLevel

	writer := NewDefaultWriter[any](os.Stdout, &TextFormatter{}) //out, formatter)

	// Read the VMode from Config
	//VMode := VConfig.Mode()
	//if VMode != ModeService && VMode != ModeStandalone {
	mode := ModeStandalone // Default to standalone if not specified
	//}

	logg := log.New(
		writer.out,
		prefix,
		log.LstdFlags,
	)

	lgr := &LogzCoreImpl{
		Logger:    logg,
		VLevel:    level,
		VWriter:   writer,
		VMetadata: make(map[string]interface{}),
		VMode:     mode,
	}

	lgr.prefix.Store(&prefix)

	return lgr
}

// SetMetadata sets a VMetadata key-value pair for the LogzCoreImpl.
func (l *LogzCoreImpl) SetMetadata(key string, value interface{}) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.VMetadata[key] = value
}

// shouldLog checks if the log VLevel should be logged.
func (l *LogzCoreImpl) shouldLog(level LogLevel) bool {
	return logLevels[level] >= logLevels[l.VLevel]
}

// log logs a message with the specified VLevel and context.
func (l *LogzCoreImpl) log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	l.Mu.RLock()
	defer l.Mu.RUnlock()

	entry := NewLogEntry().
		WithLevel(level).
		WithMessage(msg).
		WithSeverity(logLevels[level])

	// Merge global and local VMetadata
	finalContext := mergeContext(l.VMetadata, ctx)
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	// Merge global and local VMetadata
	finalMetadata := mergeMetadata(l.VMetadata, ctx)
	for k, v := range finalMetadata {
		entry.AddMetadata(k, v)
	}

	if level != SILENT {
		// Write the log using the configured VWriter
		if err := l.VWriter.Write(entry); err != nil {
			log.Printf("ErrorCtx writing log: %v", err)
		}
	}

	// Update metrics in PrometheusManager, if enabled
	if l.VMode == ModeService {
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

// TraceCtx logs a trace message with context.
func (l *LogzCoreImpl) TraceCtx(msg string, ctx map[string]interface{}) { l.log(TRACE, msg, ctx) }

// NoticeCtx logs a notice message with context.
func (l *LogzCoreImpl) NoticeCtx(msg string, ctx map[string]interface{}) { l.log(NOTICE, msg, ctx) }

// SuccessCtx logs a success message with context.
func (l *LogzCoreImpl) SuccessCtx(msg string, ctx map[string]interface{}) { l.log(SUCCESS, msg, ctx) }

// DebugCtx logs a debug message with context.
func (l *LogzCoreImpl) DebugCtx(msg string, ctx map[string]interface{}) { l.log(DEBUG, msg, ctx) }

// InfoCtx logs an info message with context.
func (l *LogzCoreImpl) InfoCtx(msg string, ctx map[string]interface{}) { l.log(INFO, msg, ctx) }

// WarnCtx logs a warning message with context.
func (l *LogzCoreImpl) WarnCtx(msg string, ctx map[string]interface{}) { l.log(WARN, msg, ctx) }

// ErrorCtx logs an error message with context.
func (l *LogzCoreImpl) ErrorCtx(msg string, ctx map[string]interface{}) { l.log(ERROR, msg, ctx) }

// FatalCtx logs a fatal message with context and terminates the process.
func (l *LogzCoreImpl) FatalCtx(msg string, ctx map[string]interface{}) { l.log(FATAL, msg, ctx) }

// SilentCtx logs a message with context without any output.
func (l *LogzCoreImpl) SilentCtx(msg string, ctx map[string]interface{}) { l.log(SILENT, msg, ctx) }

// AnswerCtx logs an answer message with context.
func (l *LogzCoreImpl) AnswerCtx(msg string, ctx map[string]interface{}) { l.log(ANSWER, msg, ctx) }

// Silent logs a message without any output.
func (l *LogzCoreImpl) Silent(msg ...any) {
	if l.shouldLog(SILENT) {
		l.log(SILENT, fmt.Sprint(msg...), nil)
	}
}

// Answer logs a message without any output.
func (l *LogzCoreImpl) Answer(msg ...any) {
	if l.shouldLog(ANSWER) {
		l.log(ANSWER, fmt.Sprint(msg...), nil)
	}
}
func (l *LogzCoreImpl) SetLevel(level interface{}) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if lvl, ok := level.(LogLevel); ok {
		l.VLevel = lvl
	} else if lvlStr, ok := level.(string); ok {
		l.VLevel = LogLevel(lvlStr)
	} else {
		log.Println("Invalid log level type")
	}
}
func (l *LogzCoreImpl) GetLevel() interface{} {
	l.Mu.RLock()
	defer l.Mu.RUnlock()

	if l.VLevel == "" {
		l.VLevel = INFO
	}
	return l.VLevel
}
func (l *LogzCoreImpl) SetWriter(writer any) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if osFile, ok := writer.(*os.File); ok {
		l.VWriter = NewDefaultWriter[any](osFile, &TextFormatter{})
	} else if logWriter, ok := writer.(LogWriter[any]); ok {
		l.VWriter = logWriter
	} else {
		log.Println("Invalid writer type")
	}
}
func (l *LogzCoreImpl) GetWriter() interface{} {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VWriter == nil {
		l.VWriter = NewDefaultWriter[any](os.Stdout, &TextFormatter{})
	}
	return l.VWriter
}
func (l *LogzCoreImpl) GetMode() interface{} {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VMode == "" {
		l.VMode = ModeStandalone
	}
	return l.VMode
}
func (l *LogzCoreImpl) SetConfig(config interface{}) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if cfg, ok := config.(Config); ok {
		l.VConfig = cfg
	} else {
		log.Println("Invalid config type")
	}
}
func (l *LogzCoreImpl) GetConfig() interface{} {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VConfig == nil {
		c := NewConfigManager()
		c2 := *c
		c3 := c2.GetConfig()
		l.VConfig = c3
	}
	return l.VConfig
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
