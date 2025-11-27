// Package interfaces defines the Entry interface for log entries.
package interfaces

import "time"

type Entry interface {
	// --- Chainable ---

	WithLevel(l Level) Entry
	WithMessage(msg string) Entry
	WithContext(ctx string) Entry
	WithSource(src string) Entry
	WithTraceID(id string) Entry
	WithError(err error) Entry
	WithField(key string, value any) Entry
	WithFields(fields map[string]any) Entry
	WithData(data any) Entry

	// ---- Tags and Fields ---

	Tag(k, v string) Entry
	Field(k string, v any) Entry
	WithCaller(c string) Entry
	CaptureCaller(skip int) Entry
	Clone() Entry
	GetLevel() Level
	Validate() error
	String() string

	GetTimestamp() time.Time
	GetContext() string
	GetMessage() string
	GetCaller() string
	GetTags() map[string]string
	GetFields() map[string]any
}
