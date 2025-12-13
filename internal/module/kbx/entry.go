package kbx

import "time"

type Record interface {
	GetTimestamp() time.Time
	GetContext() string
	GetMessage() string
	GetPrefix() string

	// GetData() any

	Validate() error
	String() string
}

type Entry interface {
	Record

	// ---- Tags and Fields ---
	Clone() Entry
	GetPrefix() string
	GetLevel() Level
	GetMessage() string
	GetCaller() string
	GetTags() map[string]string
	GetFields() map[string]any
	GetTraceID() string
	GetShowColor() bool
	GetShowStack() bool
	GetShowCaller() bool
	GetShowFields() bool
	GetShowIcon() bool
	GetFormat() string
	GetShowTraceID() bool

	// --- Chainable ---
	CaptureCaller(skip int) Entry
}

type LogzEntry interface {
	// --- Base ---

	Entry

	// ---- Tags and Fields ---

	GetLevel() Level
	GetCaller() string
	GetTags() map[string]string
	GetFields() map[string]any
	GetPrefix() string
	GetMessage() string
	GetTraceID() string
	GetShowColor() bool
	GetShowStack() bool
	GetShowCaller() bool
	GetShowFields() bool
	GetShowIcon() bool
	GetShowTraceID() bool
	GetFormat() string

	// --- Chainable ---

	WithLevel(l Level) LogzEntry
	WithMessage(msg string) LogzEntry
	WithContext(ctx string) LogzEntry
	WithSource(src string) LogzEntry
	WithTraceID(id string) LogzEntry
	WithError(err error) LogzEntry
	WithField(key string, value any) LogzEntry
	WithFields(fields map[string]any) LogzEntry
	WithData(data any) LogzEntry

	WithFormat(format string) LogzEntry
	WithColor(color bool) LogzEntry
	WithStack(stack bool) LogzEntry
	WithIcon(icon bool) LogzEntry
	WithShowTraceID(show bool) LogzEntry
	WithShowCaller(show bool) LogzEntry
	WithShowFields(show bool) LogzEntry

	Tag(k, v string) LogzEntry
	WithCaller(c string) LogzEntry

	Field(k string, v any) Entry
}

type logzEntry[E Entry, T LogzEntry] interface {
	// --- Base ---

	Entry

	// ---- Tags and Fields ---

	GetLevel() Level
	GetCaller() string
	GetTags() map[string]string
	GetFields() map[string]any
	GetTraceID() string
	GetPrefix() string
	GetMessage() string
	GetShowColor() bool
	GetShowStack() bool
	GetShowCaller() bool
	GetShowFields() bool
	GetShowIcon() bool
	GetFormat() string
	GetShowTraceID() bool

	// --- Chainable ---

	WithLevel(l Level) T
	WithMessage(msg string) T
	WithContext(ctx string) T
	WithSource(src string) T
	WithTraceID(id string) T
	WithError(err error) T
	WithField(key string, value any) T
	WithFields(fields map[string]any) T
	WithData(data any) T

	WithFormat(format string) T
	WithColor(color bool) T
	WithStack(stack bool) T
	WithIcon(icon bool) T
	WithShowTraceID(show bool) T
	WithShowCaller(show bool) T
	WithShowFields(show bool) T

	Tag(k, v string) T
	Field(k string, v any) T
	WithCaller(c string) T
}
