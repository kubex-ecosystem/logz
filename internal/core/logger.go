package core

import (
	"io"
	"sync"
	"time"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// LoggerZ é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type LoggerZ[T interfaces.Entry] struct {
	mu        sync.RWMutex
	formatter interfaces.FormatterG[T]
	out       io.Writer
	minLevel  interfaces.Level
	hooks     []interfaces.HookG[T]
	*Logger
}

// LoggerG é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type LoggerG[T interfaces.Entry] struct {
	mu        sync.RWMutex
	formatter interfaces.FormatterG[T]
	out       io.Writer
	minLevel  interfaces.Level
	hooks     []interfaces.HookG[T]
	*Logger
}

// NewLoggerG cria um logger genérico:
// - formatter: serializa T em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLoggerG[T interfaces.Entry](formatter interfaces.FormatterG[T], out io.Writer, min interfaces.Level) *LoggerG[T] {
	if min == "" {
		min = interfaces.LevelDebug
	}
	return &LoggerG[T]{
		formatter: formatter,
		out:       out,
		minLevel:  min,
	}
}

type Logger struct {
	mu        sync.RWMutex
	formatter interfaces.Formatter
	out       io.Writer
	minLevel  interfaces.Level
	hooks     []interfaces.Hook
	level     interfaces.Level
	buffer    []byte
	lastEntry *Entry
	ticker    *time.Ticker
	flushMu   sync.Mutex
	hooksMu   sync.Mutex
	lHooks    []interfaces.LHook[interfaces.Entry]
}

// NewLogger cria um logger genérico:
// - formatter: serializa Record em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLogger(formatter interfaces.Formatter, out io.Writer, min interfaces.Level) *Logger {
	if min == "" {
		min = interfaces.LevelDebug
	}
	return &Logger{
		formatter: formatter,
		out:       out,
		minLevel:  min,
		level:     min,
	}
}

func (l *Logger) SetFormatter(f interfaces.FormatterG[interfaces.Entry]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = f
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetMinLevel(min interfaces.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = min
}

func (l *Logger) AddHook(h interfaces.Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, h)
}

func (l *Logger) Enabled(level interfaces.Level) bool {
	l.mu.RLock()
	min := l.minLevel
	l.mu.RUnlock()
	return level.Severity() >= min.Severity()
}

func (l *Logger) GetMinLevel() interfaces.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.minLevel
}

func (l *Logger) GetLevel() interfaces.Level {
	return l.minLevel
}

// Log é o caminho principal: recebe um Record pronto (T),
// dispara hooks, formata e escreve em out.
func (l *Logger) Log(rec interfaces.Entry) error {
	if !kbx.IsObjSafe(rec, false) {
		// nada a fazer, mas não vamos quebrar ninguém
		return nil
	}

	// r := *rec // copia pra evitar alterações concorrentes

	if !l.Enabled(interfaces.Level(rec.GetLevel().String())) {
		return nil
	}

	if err := rec.Validate(); err != nil {
		return err
	}

	l.mu.RLock()
	f := l.formatter
	out := l.out

	hooks := append([]interfaces.LHook[interfaces.Entry](nil), l.lHooks...)

	l.mu.RUnlock()

	if f == nil || out == nil {
		// logger não inicializado corretamente; falha silenciosa
		return nil
	}

	// hooks antes da formatação
	for _, h := range hooks {
		h.Fire(rec)
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
