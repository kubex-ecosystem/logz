package kbx

import "time"

type Record interface {
	GetTimestamp() time.Time
	GetContext() string
	GetMessage() string

	// GetData() any

	Validate() error
	String() string
}

type Entry interface {
	Record

	// ---- Tags and Fields ---
	Clone() Entry
	GetLevel() Level
	GetCaller() string
	GetTags() map[string]string
	GetFields() map[string]any
}

type LogzEntry[T Entry] interface {
	Entry

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

	Tag(k, v string) T
	Field(k string, v any) T
	WithCaller(c string) T
	CaptureCaller(skip int) T
}
