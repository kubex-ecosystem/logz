package interfaces

import (
	"io"
	"time"
)

type Logger interface {
	SetFormatter(f Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	Enabled(level Level) bool
	GetMinLevel() Level
	GetLevel() Level
	SetRotate(rotate bool)
	SetRotateMaxSize(size int64)
	SetRotateMaxBack(back int64)
	SetRotateMaxAge(age int64)
	SetCompress(compress bool)
	SetBufferSize(size int)
	SetFlushInterval(interval time.Duration)
	SetHooks(hooks []Hook)
	SetLHooks(hooks LHook[any])
	SetMetadata(metadata map[string]any)
	Log(lvl string, rec Entry) error
	LogAny(args ...any) error

	AddHook(h Hook)
}

type LoggerZ[T Hook | *Entry] interface {
	SetFormatter(f Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	Enabled(level Level) bool
	Log(rec T) error

	AddHook(h HookG[T])
}

type LoggerFunc interface {
	SetFormatter(f FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	Enabled(level Level) bool
	Log(rec Entry) error

	AddHook(h HookFunc)
}

type LoggerFuncG[T Hook | *Entry] interface {
	SetFormatter(f FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	Enabled(level Level) bool
	Log(rec T) error

	AddHook(h HookFuncG[T])
}
