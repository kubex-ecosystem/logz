package kbx

import (
	"context"
	"time"
)

type Record interface {
	GetTimestamp() time.Time
	GetContext() string
	GetLevel() string
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

	WithLevel(l string) LogzEntry
	WithMessage(msg any) LogzEntry
	WithMessages(msg ...any) LogzEntry
	WithSource(src string) LogzEntry
	WithTraceID(id string) LogzEntry
	WithError(err error) LogzEntry
	WithField(key string, value any) LogzEntry
	WithFields(fields map[string]any) LogzEntry
	WithData(data any) LogzEntry

	WithContext(ctx context.Context) LogzEntry

	WithFormat(format string) LogzEntry
	WithColor(color bool) LogzEntry
	WithStack(stack bool) LogzEntry
	WithIcon(icon bool) LogzEntry
	WithShowTraceID(show bool) LogzEntry
	WithShowCaller(show bool) LogzEntry
	WithShowFields(show bool) LogzEntry

	Tag(k, v string) LogzEntry
	WithCaller(c string) LogzEntry

	Field(k string, v any) LogzEntry
	Error() error

	Reset()
}
