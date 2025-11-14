package loggerz

import (
	"fmt"
	"io"

	li "github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/interfaces/formatters"
)

// DefaultWriter is a generic writer that implements the Writer interface.
// It can write log entries of any type T to an io.Writer.
// It uses a LogFormatter to format the log entries before writing them.
type DefaultWriter[T any] struct {
	out       io.Writer
	formatter formatters.LogFormatter
}

// NewDefaultWriter cria um novo VWriter usando generics.
func NewDefaultWriter[T any](out io.Writer, formatter formatters.LogFormatter) *DefaultWriter[T] {
	return &DefaultWriter[T]{
		out:       out,
		formatter: formatter,
	}
}

// Write aceita qualquer tipo de entrada T e a processa.
func (w *DefaultWriter[T]) Write(entry T) (int, error) {
	var formatted string
	var err error

	// Verifique se a entrada é do tipo LogzEntry
	switch v := any(entry).(type) {
	case li.LogzEntry:
		formatted, err = w.formatter.Format(v)
	case []byte:
		// Converta o []byte em LogzEntry antes de formatar (exemplo simplificado)
		entry := NewLogEntry().WithMessage(string(v))
		formatted, err = w.formatter.Format(entry)
	default:
		return 0, fmt.Errorf("unsupported log entry type: %T", entry)
	}

	if err != nil {
		return 0, err
	}

	_, err = fmt.Fprintln(w.out, formatted)
	return 0, err
}

// Out returns the underlying io.Writer.
func (w *DefaultWriter[T]) Out() io.Writer {
	return w.out
}
