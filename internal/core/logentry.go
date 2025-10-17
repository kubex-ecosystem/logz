package core

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

// LogLevel represents the severity VLevel of a log entry.
type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	TRACE   LogLevel = "TRACE"
	NOTICE  LogLevel = "NOTICE"
	INFO    LogLevel = "INFO"
	SUCCESS LogLevel = "SUCCESS"
	WARN    LogLevel = "WARN"
	ERROR   LogLevel = "ERROR"
	FATAL   LogLevel = "FATAL"
	PANIC   LogLevel = "PANIC"
	SILENT  LogLevel = "SILENT"
	ANSWER  LogLevel = "ANSWER"
)

// LogzEntry represents a single log entry with various attributes.
type LogzEntry interface {
	// WithLevel sets the log VLevel for the LogEntry.
	WithLevel(level LogLevel) LogzEntry
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
	GetLevel() LogLevel
	// GetSource returns the source of the LogEntry.
	GetSource() string
	// Validate checks if the LogEntry has all required fields set.
	Validate() error
	// String returns a string representation of the LogEntry.
	String() string
}

// LogEntry represents a single log entry with various attributes.
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`           // The time when the log entry was created.
	Level     LogLevel               `json:"VLevel"`              // The severity VLevel of the log entry.
	Source    string                 `json:"source"`              // The source of the log entry.
	Context   string                 `json:"context,omitempty"`   // Additional context for the log entry.
	Message   string                 `json:"message"`             // The log message.
	Tags      map[string]string      `json:"tags,omitempty"`      // Optional tags associated with the log entry.
	Metadata  map[string]interface{} `json:"VMetadata,omitempty"` // Optional VMetadata associated with the log entry.
	ProcessID int                    `json:"pid,omitempty"`       // The process ID associated with the log entry.
	Hostname  string                 `json:"hostname,omitempty"`  // The hostname where the log entry was created.
	Severity  int                    `json:"severity"`            // The severity VLevel as an integer.
	TraceID   string                 `json:"trace_id,omitempty"`  // Optional trace ID for tracing logs.
	Caller    string                 `json:"caller,omitempty"`    // The caller of the log entry.
}

// NewLogEntry creates a new instance of LogEntry with the current timestamp and initialized maps.
func NewLogEntry() LogzEntry {
	le := LogEntry{
		Timestamp: time.Now(),
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
		Caller:    getCallerInfo(3),
	}
	return &le
}

// WithLevel sets the log VLevel for the LogEntry.
func (le *LogEntry) WithLevel(level LogLevel) LogzEntry {
	le.Level = level
	return le
}

// WithSource sets the source for the LogEntry.
func (le *LogEntry) WithSource(source string) LogzEntry {
	le.Source = source
	return le
}

// WithContext sets the context for the LogEntry.
func (le *LogEntry) WithContext(context string) LogzEntry {
	le.Context = context
	return le
}

// WithMessage sets the message for the LogEntry.
func (le *LogEntry) WithMessage(message string) LogzEntry {
	le.Message = message
	return le
}

// WithProcessID sets the process ID for the LogEntry.
func (le *LogEntry) WithProcessID(pid int) LogzEntry {
	le.ProcessID = pid
	return le
}

// WithHostname sets the hostname for the LogEntry.
func (le *LogEntry) WithHostname(hostname string) LogzEntry {
	le.Hostname = hostname
	return le
}

// WithSeverity sets the severity VLevel for the LogEntry.
func (le *LogEntry) WithSeverity(severity int) LogzEntry {
	le.Severity = severity
	return le
}

// WithTraceID sets the trace ID for the LogEntry.
func (le *LogEntry) WithTraceID(traceID string) LogzEntry {
	le.TraceID = traceID
	return le
}

// AddTag adds a tag to the LogEntry.
func (le *LogEntry) AddTag(key, value string) LogzEntry {
	if le.Tags == nil {
		le.Tags = make(map[string]string)
	}
	le.Tags[key] = value
	return le
}

// GetLevel returns the log VLevel of the LogEntry.
func (le *LogEntry) GetLevel() LogLevel { return le.Level }

// AddMetadata adds VMetadata to the LogEntry.
func (le *LogEntry) AddMetadata(key string, value interface{}) LogzEntry {
	if le.Metadata == nil {
		le.Metadata = make(map[string]interface{})
	}
	le.Metadata[key] = value
	return le
}

// GetMetadata returns the VMetadata of the LogEntry.
func (le *LogEntry) GetMetadata() map[string]interface{} { return le.Metadata }

// GetContext returns the context of the LogEntry.
func (le *LogEntry) GetContext() string { return le.Context }

// GetTimestamp returns the timestamp of the LogEntry.
func (le *LogEntry) GetTimestamp() time.Time { return le.Timestamp }

// GetMessage returns the message of the LogEntry.
func (le *LogEntry) GetMessage() string { return le.Message }

// GetSource returns the source of the LogEntry.
func (le *LogEntry) GetSource() string { return le.Source }

// Validate checks if the LogEntry has all required fields set.
func (le *LogEntry) Validate() error {
	if le.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if le.Level == "" {
		return errors.New("VLevel is required")
	}
	if le.Message == "" {
		return errors.New("message is required")
	}
	if le.Severity <= 0 {
		return errors.New("severity must be greater than zero")
	}
	return nil
}

// String returns a string representation of the LogEntry.
func (le *LogEntry) String() string {
	if le == nil {
		return "LogEntry: <nil>"
	}
	if le.Level == SILENT || le.Level == ANSWER {
		return fmt.Sprintf("%s - %s",
			le.Level,
			le.Message,
		)
	} else {
		return fmt.Sprintf("[%s] %s - %s",
			le.Timestamp.Format(time.RFC3339),
			le.Level,
			le.Message,
		)
	}
}

// getCallerInfo returns the caller information for the log entry.
func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", trimFilePath(file), line, funcName)
}
