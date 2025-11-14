package loggerz

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"

	//"io"
	"os"
	"strings"
	"sync"

	li "github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	si "github.com/kubex-ecosystem/logz/internal/services"

	"github.com/kubex-ecosystem/logz/internal/interfaces/formatters"
)

// func init() {
// 	if LoggerG == nil {

// 		LoggerG = NewLoggerZ("")

// 		if logger, ok := LoggerG.(*LoggerImpl); ok {
// 			g = logger
// 			logLevel = getEnvOrDefault("GOBE_LOG_LEVEL", "error")
// 			debug = getEnvOrDefault("GOBE_DEBUG", false)
// 			showTrace = getEnvOrDefault("GOBE_SHOW_TRACE", false)
// 			g.gLogLevel = il.INFO
// 			g.gShowTrace = showTrace
// 			g.gDebug = debug
// 		}
// 	}
// }

// LogzCoreImpl represents a core with configuration and VMetadata.
type LogzCoreImpl struct {
	// LogzLogger is a constraint to implement this interface
	li.LogzLogger

	// Logger is a promoted global Go Logger
	*log.Logger

	out       io.Writer                      // destination for output
	prefix    atomic.Pointer[string]         // prefix on each line to identify the logger (but see Lmsgprefix)
	prefixX   atomic.Pointer[*li.LogzLogger] // prefix on each line to identify the logger (but see Lmsgprefix)
	flag      atomic.Int32                   // properties
	isDiscard atomic.Bool

	VLevel    LogLevel
	VWriter   io.Writer
	VConfig   li.Config
	VMetadata map[string]any
	VMode     kbx.LogMode // Mode control: service or standalone
	Mu        sync.RWMutex
}

// NewLogger creates a new instance of LogzCoreImpl with the provided configuration.
func NewLogger(prefix string) li.LogzLogger {
	return NewLoggerImpl(prefix)
}

func NewLoggerImpl(prefix string) *LogzCoreImpl {
	level := kbx.INFO // Default log VLevel

	writer := NewDefaultWriter[[]byte](os.Stdout, &formatters.TextFormatterImpl{}) //out, formatter)

	// Read the VMode from Config
	//VMode := VConfig.Mode()
	//if VMode != ModeService && VMode != ModeStandalone {
	mode := kbx.ModeStandalone // Default to standalone if not specified
	//}

	logg := log.New(
		writer.Out(),
		prefix,
		log.LstdFlags,
	)

	lgr := &LogzCoreImpl{
		Logger:    logg,
		VLevel:    level,
		VWriter:   writer,
		VMetadata: make(map[string]any),
		VMode:     mode,
	}

	lgr.prefix.Store(&prefix)

	return lgr
}

// SetMetadata sets a VMetadata key-value pair for the LogzCoreImpl.
func (l *LogzCoreImpl) SetMetadata(key string, value any) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.VMetadata[key] = value
}

// shouldLog checks if the log VLevel should be logged.
func (l *LogzCoreImpl) shouldLog(level string) bool {
	return kbx.LogLevels[kbx.LogLevel(strings.ToUpper(level))] >= kbx.LogLevels[l.VLevel]
}

// log logs a message with the specified VLevel and context.
func (l *LogzCoreImpl) log(level string, msg string, ctx map[string]any) {
	if !l.shouldLog(level) {
		return
	}

	l.Mu.RLock()
	defer l.Mu.RUnlock()

	entry := NewLogEntry().
		WithLevel(level).
		WithMessage(msg).
		WithSeverity(kbx.LogLevels[kbx.LogLevel(strings.ToUpper(level))])

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

	if strings.ToUpper(level) != string(kbx.SILENT) {
		// Write the log using the configured VWriter
		if _, err := l.VWriter.Write([]byte(entry.String())); err != nil {
			log.Printf("ErrorCtx writing log: %v", err)
		}
	}

	// Update metrics in PrometheusManager, if enabled
	if l.VMode == kbx.ModeService {
		pm := si.GetPrometheusManager()
		if pm.IsEnabled() {
			pm.IncrementMetric("logs_total", 1)
			pm.IncrementMetric("logs_total_"+string(level), 1)
		}
	}

	// Terminate the process in case of FATAL log
	if strings.ToUpper(level) == string(kbx.FATAL) {
		os.Exit(1)
	}
}

// TraceCtx logs a trace message with context.
func (l *LogzCoreImpl) TraceCtx(msg string, ctx map[string]any) { l.log(string(kbx.TRACE), msg, ctx) }

// NoticeCtx logs a notice message with context.
func (l *LogzCoreImpl) NoticeCtx(msg string, ctx map[string]any) { l.log(string(kbx.NOTICE), msg, ctx) }

// SuccessCtx logs a success message with context.
func (l *LogzCoreImpl) SuccessCtx(msg string, ctx map[string]any) {
	l.log(string(kbx.SUCCESS), msg, ctx)
}

// DebugCtx logs a debug message with context.
func (l *LogzCoreImpl) DebugCtx(msg string, ctx map[string]any) { l.log(string(kbx.DEBUG), msg, ctx) }

// InfoCtx logs an info message with context.
func (l *LogzCoreImpl) InfoCtx(msg string, ctx map[string]any) { l.log(string(kbx.INFO), msg, ctx) }

// WarnCtx logs a warning message with context.
func (l *LogzCoreImpl) WarnCtx(msg string, ctx map[string]any) { l.log(string(kbx.WARN), msg, ctx) }

// ErrorCtx logs an error message with context.
func (l *LogzCoreImpl) ErrorCtx(msg string, ctx map[string]any) { l.log(string(kbx.ERROR), msg, ctx) }

// FatalCtx logs a fatal message with context and terminates the process.
func (l *LogzCoreImpl) FatalCtx(msg string, ctx map[string]any) { l.log(string(kbx.FATAL), msg, ctx) }

// SilentCtx logs a message with context without any output.
func (l *LogzCoreImpl) SilentCtx(msg string, ctx map[string]any) { l.log(string(kbx.SILENT), msg, ctx) }

// AnswerCtx logs an answer message with context.
func (l *LogzCoreImpl) AnswerCtx(msg string, ctx map[string]any) { l.log(string(kbx.ANSWER), msg, ctx) }

// Silent logs a message without any output.
func (l *LogzCoreImpl) Silent(msg ...any) {
	if l.shouldLog(string(kbx.SILENT)) {
		l.log(string(kbx.SILENT), fmt.Sprint(msg...), nil)
	}
}

func (l *LogzCoreImpl) GetConfig() li.Config {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VConfig == nil {
		// c := NewConfigManager()
		// c2 := *c
		// c3 := c2.GetConfig()
		// l.VConfig = c3
	}
	return l.VConfig
}
func (l *LogzCoreImpl) SetConfig(config li.Config) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.VConfig = config
}

// Answer logs a message without any output.
func (l *LogzCoreImpl) Answer(msg ...any) {
	if l.shouldLog(string(kbx.ANSWER)) {
		l.log(string(kbx.ANSWER), fmt.Sprint(msg...), nil)
	}
}
func (l *LogzCoreImpl) SetLevel(level string) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, exists := kbx.LogLevels[kbx.LogLevel(strings.ToUpper(level))]; exists {
		l.VLevel = kbx.LogLevel(strings.ToUpper(level))
	} else {
		log.Printf("Invalid %s log level type", level)
	}
}

func (l *LogzCoreImpl) GetLevel() string {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VLevel == "" {
		l.VLevel = kbx.INFO
	}
	return strings.ToLower(string(l.VLevel))
}
func (l *LogzCoreImpl) SetWriter(writer io.Writer) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if osFile, ok := writer.(*os.File); ok {
		l.VWriter = NewDefaultWriter[[]byte](osFile, &formatters.TextFormatterImpl{})
	} else if logWriter, ok := writer.(io.Writer); ok {
		l.VWriter = logWriter
	} else {
		log.Println("Invalid writer type")
	}
}
func (l *LogzCoreImpl) GetWriter() io.Writer {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VWriter == nil {
		l.VWriter = NewDefaultWriter[[]byte](os.Stdout, &formatters.TextFormatterImpl{})
	}
	return l.VWriter
}
func (l *LogzCoreImpl) GetMode() interface{} {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	if l.VMode == "" {
		l.VMode = kbx.ModeStandalone
	}
	return l.VMode
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
func mergeContext(global, local map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}
	return merged
}

// mergeMetadata merges global and local context maps.
func mergeMetadata(global, local map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}
	return merged
}
