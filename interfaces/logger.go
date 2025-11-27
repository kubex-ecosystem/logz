package interfaces

import (
	"io"
)

type Logger interface {
	SetFormatter(f Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h Hook)
	Enabled(level Level) bool
	Log(rec Entry) error
}

type LoggerZ[T Entry] interface {
	SetFormatter(f FormatterG[T])
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h HookG[T])
	Enabled(level Level) bool
	Log(rec T) error
}

type LoggerFunc interface {
	SetFormatter(f FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h HookFunc)
	Enabled(level Level) bool
	Log(rec Entry) error
}

type LoggerFuncG[T Entry] interface {
	SetFormatter(f FormatterFuncG[T])
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h HookFuncG[T])
	Enabled(level Level) bool
	Log(rec T) error
}
