package core

import (
	"fmt"
	"io"
)

// LogWriter defines the contract for writing logs.
type LogWriter[T any] interface {
	Write(entry T) error
}

// DefaultWriter is a generic writer that implements the LogWriter interface.
// It can write log entries of any type T to an io.Writer.
// It uses a LogFormatter to format the log entries before writing them.
type DefaultWriter[T any] struct {
	out       io.Writer
	formatter LogFormatter
}

// NewDefaultWriter cria um novo VWriter usando generics.
func NewDefaultWriter[T any](out io.Writer, formatter LogFormatter) *DefaultWriter[T] {
	return &DefaultWriter[T]{
		out:       out,
		formatter: formatter,
	}
}

// Write aceita qualquer tipo de entrada T e a processa.
func (w *DefaultWriter[T]) Write(entry T) error {
	var formatted string
	var err error

	// Verifique se a entrada é do tipo LogzEntry
	switch v := any(entry).(type) {
	case LogzEntry:
		formatted, err = w.formatter.Format(v)
	case []byte:
		// Converta o []byte em LogzEntry antes de formatar (exemplo simplificado)
		entry := NewLogEntry().WithMessage(string(v))
		formatted, err = w.formatter.Format(entry)
	default:
		return fmt.Errorf("unsupported log entry type: %T", entry)
	}

	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w.out, formatted)
	return err
}

// formatMetadata converts VMetadata to a JSON string.
// Returns the JSON string or an empty string if marshalling fails.
func formatMetadata(entry LogzEntry) string {
	metadata := entry.GetMetadata()
	if len(metadata) == 0 {
		return ""
	}
	prefix := "Context:\n"
	for k, v := range metadata {
		if k == "showContext" {
			continue
		}
		prefix += fmt.Sprintf("  - %s: %v\n", k, v)
	}
	return prefix
}
