package interfaces

import (
	"io"
	"time"
)

type Logger interface {
	SetFormatter(f Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h Hook)
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
}

type LoggerZ[T Entry] interface {
	SetFormatter(f Formatter)
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
	SetFormatter(f FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min Level)
	AddHook(h HookFuncG[T])
	Enabled(level Level) bool
	Log(rec T) error
}
