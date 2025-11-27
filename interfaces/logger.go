package interfaces

import (
	"io"
	"time"

	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type Logger interface {
	SetFormatter(f formatter.Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool
	GetMinLevel() kbx.Level
	GetLevel() kbx.Level
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
	Log(lvl string, rec kbx.Entry) error
	LogAny(args ...any) error

	AddHook(h Hook)
}

type LoggerZ[T Hook | *kbx.Entry] interface {
	SetFormatter(f formatter.Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool
	Log(rec T) error

	AddHook(h HookG[T])
}

type LoggerFunc interface {
	SetFormatter(f formatter.FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool
	Log(rec kbx.Entry) error

	AddHook(h HookFunc)
}

type LoggerFuncG[T Hook | *kbx.Entry] interface {
	SetFormatter(f formatter.FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool
	Log(rec T) error

	AddHook(h HookFuncG[T])
}
