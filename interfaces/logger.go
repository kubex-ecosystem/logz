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

	Log(lvl kbx.Level, args ...any) error
	LogAny(level kbx.Level, args ...any) error

	AddHook(h Hook)
}


type LoggerZ[T Hook | *kbx.Entry] interface {
	SetFormatter(f formatter.Formatter)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool

	Log(level kbx.Level, rec T) error

	AddHook(h HookG[T])

	SetDebugMode(debug bool)
	Debug(msg ...any)
	Notice(msg ...any)
	Info(msg ...any)
	Success(msg ...any)
	Warn(msg ...any)
	Error(msg ...any) error
	Fatal(msg ...any)
	Trace(msg ...any)
	Critical(msg ...any)
	Answer(msg ...any)
	Alert(msg ...any)
	Bug(msg ...any)
	Panic(msg ...any)
	Println(msg ...any)
	Printf(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Noticef(format string, args ...any)
	Successf(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any) error
	Fatalf(format string, args ...any)
	Tracef(format string, args ...any)
	Criticalf(format string, args ...any)
	Answerf(format string, args ...any)
	Alertf(format string, args ...any)
	Bugf(format string, args ...any)
	Panicf(format string, args ...any)
}

type LoggerFunc interface {
	SetFormatter(f formatter.FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool

	Log(level kbx.Level, rec kbx.Entry) error

	AddHook(h HookFunc)
}

type LoggerFuncG[T Hook | *kbx.Entry] interface {
	SetFormatter(f formatter.FormatterFunc)
	SetOutput(w io.Writer)
	SetMinLevel(min kbx.Level)
	Enabled(level kbx.Level) bool

	Log(level kbx.Level, rec T) error

	AddHook(h HookFuncG[T])
}
