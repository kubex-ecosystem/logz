// Package interfaces defines the Entry interface for log entries.
package interfaces

type Entry interface {
	WithLevel(l Level) *Entry
	WithMessage(msg string) *Entry
	WithContext(ctx string) *Entry
	WithSource(src string) *Entry
	WithTraceID(id string) *Entry
	Tag(k, v string) *Entry
	Field(k string, v any) *Entry
	WithCaller(c string) *Entry
	CaptureCaller(skip int) *Entry
	Clone() *Entry
	GetLevel() Level
	Validate() error
	String() string
}
