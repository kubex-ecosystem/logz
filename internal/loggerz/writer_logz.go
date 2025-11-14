package loggerz

import (
	"fmt"
	"io"

	li "github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/interfaces/formatters"
)

type WriteLogz[E any] = func(entry E, ioWriter func([]byte) error) error

// DefaultWriter is a generic writer that implements the LogWriter interface.
// It can write log entries of any type T to an io.Writer.
// It uses a LogFormatter to format the log entries before writing them.
type LogzWriterImpl[E any, T WriteLogz[E]] struct {
	// E type of log entry and it is already fixed by go compiler.
	// We don't need to declare it again here, because it is already declared in the interface.

	// The function responsible for writing log entries.
	writeLog  T
	output    io.Writer
	formatter formatters.LogFormatter
	writer    io.Writer
}

// NewLogzWriter cria um novo VWriter usando generics.
func NewLogzWriter[E any, T WriteLogz[E]](typ E, output T, format string, writer io.Writer) LogzWriter[E, T] {
	var formatter formatters.LogFormatter
	switch format {
	case "json":
		formatter = &formatters.JSONFormatterImpl{}
	case "text":
		fallthrough
	default:
		formatter = &formatters.TextFormatterImpl{}
	}
	return &LogzWriterImpl[E, T]{
		output:    writer,
		writeLog:  output,
		formatter: formatter,
		writer:    writer,
	}
}

// Write aceita qualquer tipo de entrada T e a processa.
func (w *LogzWriterImpl[E, T]) Write(entry []byte) (int, error) {
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

	_, err = fmt.Fprintln(w.output, formatted)
	return 0, err
}

func (w *LogzWriterImpl[E, T]) WriteLogz(entry E, ioWriter T) (int, error) {
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

	_, err = fmt.Fprintln(w.output, formatted)
	return 0, err
}
