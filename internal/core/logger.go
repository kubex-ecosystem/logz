package core

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"

	"log"
)

// Logger é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type Logger struct {
	mu      sync.RWMutex
	flushMu sync.Mutex
	hooksMu sync.Mutex
	opts    *LoggerOptionsImpl
	*log.Logger
}

// LoggerZ é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type LoggerZ[T kbx.Entry] struct {
	ID       uuid.UUID
	flushMuZ sync.Mutex
	hooksMuZ sync.Mutex
	muZ      sync.RWMutex
	optsZ    *LoggerOptionsImpl
	interfaces.Logger
}

func NewLogger(prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	opts.Prefix = prefix

	// Configura o stdlog do Go para usar o mesmo output e prefixo
	var out io.Writer
	if opts.Output == nil {
		opts.Output = io.Discard
	} else {
		out = opts.Output
	}
	if out == nil {
		out = io.Discard
	}
	opts.Output = out
	logr := log.New(
		opts.Output,
		opts.Prefix,
		0,
	)

	lgr := &Logger{
		flushMu: sync.Mutex{},
		hooksMu: sync.Mutex{},
		mu:      sync.RWMutex{},
		opts:    opts,
		Logger:  logr,
	}
	// Reafirma configurações do log padrão
	lgr.SetFlags(0) // desativa flags automáticas do log padrão
	if kbx.DefaultFalse(opts.OutputTTY) {
		// se for TTY, desativa escrita direta no output padrão
		lgr.SetOutput(io.Discard) // evita escrita direta no output padrão
	} else {
		// se não for TTY, mantém escrita direta no output padrão
		lgr.SetOutput(out)
	}
	lgr.SetPrefix(prefix)
	lgr.SetFormatter(lgr.opts.Formatter)
	lgr.SetPrefix(lgr.opts.Prefix)
	lgr.SetMinLevel(lgr.opts.MinLevel)

	return lgr
}

// NewLoggerZ cria um logger genérico:
// - formatter: serializa T em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLoggerZ[T kbx.Entry](prefix string, opts *LoggerOptionsImpl, withDefaults bool) *LoggerZ[T] {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	return &LoggerZ[T]{
		ID: uuid.New(),

		muZ:      sync.RWMutex{},
		flushMuZ: sync.Mutex{},
		hooksMuZ: sync.Mutex{},

		optsZ:  opts,
		Logger: NewLogger(prefix, opts, false), // evita chamada recursiva
	}
}

// NewLogger cria um logger genérico:
// - formatter: serializa Record em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLoggerZI(prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	opts.Prefix = prefix

	// Configura o stdlog do Go para usar o mesmo output e prefixo
	var out io.Writer
	if opts.Output == nil {
		opts.Output = io.Discard
	} else {
		out = opts.Output
	}
	if out == nil {
		out = io.Discard
	}
	opts.Output = out
	logr := log.New(
		opts.Output,
		opts.Prefix,
		0,
	)

	lgr := &Logger{
		flushMu: sync.Mutex{},
		hooksMu: sync.Mutex{},
		mu:      sync.RWMutex{},
		opts:    opts,
		Logger:  logr,
	}
	// Reafirma configurações do log padrão
	lgr.SetFlags(0) // desativa flags automáticas do log padrão
	if kbx.DefaultFalse(opts.OutputTTY) {
		// se for TTY, desativa escrita direta no output padrão
		lgr.SetOutput(io.Discard) // evita escrita direta no output padrão
	} else {
		// se não for TTY, mantém escrita direta no output padrão
		lgr.SetOutput(out)
	}
	lgr.SetPrefix(prefix)
	lgr.SetFormatter(lgr.opts.Formatter)
	lgr.SetPrefix(lgr.opts.Prefix)
	lgr.SetMinLevel(lgr.opts.MinLevel)

	return lgr
}

func (l *Logger) SetFormatter(f formatter.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Formatter = f
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Output = w
}

func (l *Logger) SetMinLevel(min kbx.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.MinLevel = min
}

func (l *Logger) AddHook(h interfaces.Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Hooks = append(l.opts.Hooks, h)
}

func (l *Logger) Enabled(level kbx.Level) bool {
	l.mu.RLock()
	min := l.opts.MinLevel
	l.mu.RUnlock()
	return level.Severity() >= min.Severity()
}

func (l *Logger) GetMinLevel() kbx.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.opts.MinLevel
}

func (l *Logger) GetLevel() kbx.Level {
	return l.opts.MinLevel
}

// SetRotate is the setter for setRotate
func (l *Logger) SetRotate(rotate bool) {
	// implementação fictícia
}

// SetRotateMaxSize is the setter for setRotateMaxSize
func (l *Logger) SetRotateMaxSize(size int64) {
	// implementação fictícia
}

// SetRotateMaxBack is the setter for setRotateMaxBack
func (l *Logger) SetRotateMaxBack(back int64) {
	// implementação fictícia
}

// SetRotateMaxAge is the setter for setRotateMaxAge
func (l *Logger) SetRotateMaxAge(age int64) {
	// implementação fictícia
}

// SetCompress is the setter for setCompress
func (l *Logger) SetCompress(compress bool) {
	// implementação fictícia
}

// SetBufferSize is the setter for setBufferSize
func (l *Logger) SetBufferSize(size int) {
	// implementação fictícia
}

// SetFlushInterval is the setter for setFlushInterval
func (l *Logger) SetFlushInterval(interval time.Duration) {
	// implementação fictícia
}

// SetHooks is the setter for setHooks
func (l *Logger) SetHooks(hooks []interfaces.Hook) {
	// implementação fictícia
}

// SetLHooks is the setter for setLHooks
func (l *Logger) SetLHooks(hooks interfaces.LHook[any]) {
	// implementação fictícia
}

// SetMetadata is the setter for setMetadata
func (l *Logger) SetMetadata(metadata map[string]any) {
	// implementação fictícia
}

// Log é o caminho principal: recebe um Record pronto (T),
// dispara hooks, formata e escreve em out.
func (l *Logger) Log(lvl string, rec kbx.Entry) error {
	if !kbx.IsObjSafe(rec, false) {
		// nada a fazer, mas não vamos quebrar ninguém
		return nil
	}

	// r := *rec // copia pra evitar alterações concorrentes

	if !l.Enabled(kbx.Level(rec.GetLevel().String())) {
		return nil
	}

	if len(lvl) == 0 {
		if err := rec.Validate(); err != nil {
			return err
		}
		lvl = rec.GetLevel().String()
	}

	l.mu.RLock()
	f := l.opts.Formatter
	out := l.opts.Output
	hooks := append([]interfaces.LHook[any](nil), l.opts.LHooks)
	l.mu.RUnlock()

	if f == nil || out == nil {
		// logger não inicializado corretamente; falha silenciosa
		return nil
	}

	// hooks antes da formatação
	for _, h := range hooks {
		if h != nil {
			err := h.Fire(rec)
			if err != nil {
				return err
			}
		}
	}

	b, err := f.Format(rec)
	if err != nil {
		return err
	}

	// garante newline pra saída de console / arquivos de texto.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	_, err = out.Write(b)
	return err
}

func (l *Logger) LogAny(args ...any) error {
	if l == nil {
		return nil
	}
	if len(args) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			if l.Logger != nil {
				l.Logger.Printf("logz: panic in LogAny: %v (args=%#v)", r, args)
			}
		}
	}()

	// se args[0] é level → old API compat
	if len(args) > 1 && (kbx.IsLevel(fmt.Sprintf("%v", args[0]))) {
		level := normalizeLevel(args[0])
		entry := toEntry(args[1:]...)
		entry = entry.WithLevel(level)
		return l.Log(level.String(), entry)
	}

	// modo moderno: nada garante level → assume Info
	entry := toEntry(args...)
	if entry.GetLevel() == "" {
		entry = entry.WithLevel(kbx.LevelInfo)
	}

	return l.Log(entry.GetLevel().String(), entry)
}
