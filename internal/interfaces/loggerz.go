package interfaces

import (
	"io"
	"time"
)

// LogzEntry represents a single log entry with various attributes.
type LogzEntry interface {
	// WithLevel sets the log VLevel for the LogEntry.
	WithLevel(level string) LogzEntry
	// WithSource sets the source for the LogEntry.
	WithSource(source string) LogzEntry
	// WithContext sets the context for the LogEntry.
	WithContext(context string) LogzEntry
	// WithMessage sets the message for the LogEntry.
	WithMessage(message string) LogzEntry
	// WithProcessID sets the process ID for the LogEntry.
	WithProcessID(pid int) LogzEntry
	// WithHostname sets the hostname for the LogEntry.
	WithHostname(hostname string) LogzEntry
	// WithSeverity sets the severity VLevel for the LogEntry.
	WithSeverity(severity int) LogzEntry
	// WithTraceID sets the trace ID for the LogEntry.
	WithTraceID(traceID string) LogzEntry
	// AddTag adds a tag to the LogEntry.
	AddTag(key, value string) LogzEntry
	// AddMetadata adds VMetadata to the LogEntry.
	AddMetadata(key string, value interface{}) LogzEntry
	// GetMetadata returns the VMetadata of the LogEntry.
	GetMetadata() map[string]interface{}
	// GetContext returns the context of the LogEntry.
	GetContext() string
	// GetTimestamp returns the timestamp of the LogEntry.
	GetTimestamp() time.Time
	// GetMessage returns the message of the LogEntry.
	GetMessage() string
	// GetLevel returns the log VLevel of the LogEntry.
	GetLevel() string
	// GetSource returns the source of the LogEntry.
	GetSource() string
	// Validate checks if the LogEntry has all required fields set.
	Validate() error
	// String returns a string representation of the LogEntry.
	String() string
}

// LogzLogger combines the existing core with the standard Go log methods.
type LogzLogger interface {
	LogzCore
	// GetLevel returns the current log VLevel.
	// Method signature:
	// GetLevel() interface{}
	GetLevel() string
	// SetLevel sets the log VLevel.
	// Method signature:
	// SetLevel(VLevel interface{})
	// The VLevel is an LogLevel type or string.
	SetLevel(string)
	// Silent logs a message without any output.
	// Method signature:
	// Silent(message string, args ...any)
	// The message is a string.
	// The args are optional arguments.
	Silent(...any)
	// Answer logs an answer message.
	// Method signature:
	// Answer(message string, args ...any)
	// The message is a string.
	// The args are optional arguments.
	Answer(...any)
}

// LogzCore is the interface with the basic methods of the existing il.
type LogzCore interface {

	// SetMetadata sets a VMetadata key-value pair.
	// If the key is empty, it returns all VMetadata.
	// Returns the value and a boolean indicating if the key exists.
	SetMetadata(string, any)
	// TraceCtx logs a trace message with context.
	// Method signature:
	// TraceCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	TraceCtx(string, map[string]any)
	// NoticeCtx logs a notice message with context.
	// Method signature:
	// NoticeCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	NoticeCtx(string, map[string]any)
	// SuccessCtx logs a success message with context.
	// Method signature:
	// SuccessCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	SuccessCtx(string, map[string]any)
	// DebugCtx logs a debug message with context.
	// Method signature:
	// DebugCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	DebugCtx(string, map[string]any)
	// InfoCtx logs an informational message with context.
	// Method signature:
	// InfoCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	InfoCtx(string, map[string]any)
	// WarnCtx logs a warning message with context.
	// Method signature:
	// WarnCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	WarnCtx(string, map[string]any)
	// ErrorCtx logs an error message with context.
	// Method signature:
	// ErrorCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	ErrorCtx(string, map[string]any)
	// FatalCtx logs a fatal message with context and exits the application.
	// Method signature:
	// FatalCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	FatalCtx(string, map[string]any)
	// SilentCtx logs a message with context without any output.
	// Method signature:
	// SilentCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	SilentCtx(string, map[string]any)
	// AnswerCtx logs an answer message with context.
	// Method signature:
	// AnswerCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	AnswerCtx(string, map[string]any)
	// GetWriter returns the current VWriter.
	// Method signature:
	// GetWriter() interface{}
	// The VWriter is an interface that implements the LogWriter interface.
	GetWriter() io.Writer
	// SetWriter sets the VWriter.
	// Method signature:
	// SetWriter(VWriter interface{})
	// The VWriter is an interface that implements the LogWriter interface or io.Writer.
	SetWriter(io.Writer)
	// GetConfig returns the current configuration.
	// Method signature:
	// GetConfig() interface{}
	// The configuration is an interface that implements the Config interface.
	GetConfig() Config
	// SetConfig sets the configuration.
	SetConfig(Config)
	// SetFormat sets the format for the log entries.
	SetFormat(string)
	//// GetLevel returns the current log VLevel.
	//// Method signature:
	//// GetLevel() interface{}
	//GetLevel() interface{}
	//// SetLevel sets the log VLevel.
	//// Method signature:
	//// SetLevel(VLevel interface{})
	//// The VLevel is an LogLevel type or string.
	// SetLevel sets the log VLevel.
	//SetLevel(interface{})

	// GetLogLevel returns the current log level.
	GetLogLevel() string
	// GetShowTrace returns the current show trace flag.
	GetShowTrace() bool
	// GetDebug returns the current debug flag.
	GetDebug() bool
	// SetLogLevel sets the log level.
	SetLogLevel(string)
	// SetDebug sets the debug flag.
	SetDebug(bool)
	// SetShowTrace sets the show trace flag.
	SetShowTrace(bool)
	// ShouldLog checks if a message with the given log level should be logged.
	ShouldLog(string) bool
	// Log logs a message with the given log level.
	Log(string, ...any)
}
