package core

// LogzLogger combines the existing core with the standard Go log methods.
type LogzLogger interface {
	LogzCore
	// GetLevel returns the current log VLevel.
	// Method signature:
	// GetLevel() interface{}
	GetLevel() interface{}
	// SetLevel sets the log VLevel.
	// Method signature:
	// SetLevel(VLevel interface{})
	// The VLevel is an LogLevel type or string.
	SetLevel(interface{})
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
	SetMetadata(string, interface{})
	// TraceCtx logs a trace message with context.
	// Method signature:
	// TraceCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	TraceCtx(string, map[string]interface{})
	// NoticeCtx logs a notice message with context.
	// Method signature:
	// NoticeCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	NoticeCtx(string, map[string]interface{})
	// SuccessCtx logs a success message with context.
	// Method signature:
	// SuccessCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	SuccessCtx(string, map[string]interface{})
	// DebugCtx logs a debug message with context.
	// Method signature:
	// DebugCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	DebugCtx(string, map[string]interface{})
	// InfoCtx logs an informational message with context.
	// Method signature:
	// InfoCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	InfoCtx(string, map[string]interface{})
	// WarnCtx logs a warning message with context.
	// Method signature:
	// WarnCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	WarnCtx(string, map[string]interface{})
	// ErrorCtx logs an error message with context.
	// Method signature:
	// ErrorCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	ErrorCtx(string, map[string]interface{})
	// FatalCtx logs a fatal message with context and exits the application.
	// Method signature:
	// FatalCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	FatalCtx(string, map[string]interface{})
	// SilentCtx logs a message with context without any output.
	// Method signature:
	// SilentCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	SilentCtx(string, map[string]interface{})
	// AnswerCtx logs an answer message with context.
	// Method signature:
	// AnswerCtx(message string, context map[string]interface{})
	// The message is a string.
	// The context is a map of key-value pairs.
	AnswerCtx(string, map[string]interface{})
	// GetWriter returns the current VWriter.
	// Method signature:
	// GetWriter() interface{}
	// The VWriter is an interface that implements the LogWriter interface.
	GetWriter() interface{}
	// SetWriter sets the VWriter.
	// Method signature:
	// SetWriter(VWriter interface{})
	// The VWriter is an interface that implements the LogWriter interface or io.Writer.
	SetWriter(interface{})
	// GetConfig returns the current configuration.
	// Method signature:
	// GetConfig() interface{}
	// The configuration is an interface that implements the Config interface.
	GetConfig() interface{}
	// SetConfig sets the configuration.
	SetConfig(interface{})
	// SetFormat sets the format for the log entries.
	SetFormat(interface{})
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
	GetLogLevel() LogLevel
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
	// Log logs a message with the given log level.
	Log(string, ...any)
}
