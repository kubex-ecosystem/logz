package core

import (
	"encoding/json"
	"os"
)

// LogMultiWriter defines the contract for writing logs to multiple writers.
type LogMultiWriter[T any] interface {
	Write(entry T) error
	AddWriter(w LogWriter[T])
	GetWriters() []LogWriter[T]
}

// MultiWriter is a writer that can write to multiple LogWriters.
// It implements the LogMultiWriter interface.
// It allows adding new writers dynamically and writing log entries to all of them.
// It can be used to aggregate logs from different sources or formats.
type MultiWriter[T any] struct {
	writers []LogWriter[T]
}

// NewMultiWriter creates a new instance of MultiWriter with the provided writers.
func NewMultiWriter[T any](writers ...LogWriter[T]) *MultiWriter[T] {
	return &MultiWriter[T]{writers: writers}
}

func (mw *MultiWriter[T]) AddWriter(w LogWriter[T]) {
	mw.writers = append(mw.writers, w)
}

func (mw *MultiWriter[T]) Write(entry T) error {
	for _, w := range mw.writers {
		if err := w.Write(entry); err != nil {
			return err
		}
	}

	/* structTest := json.Marshal(entry)
	if structTest != nil {
		if strEntry, ok := entry.(string); ok {
			for _, w := range mw.writers {
				if err := w.Write(strEntry); err != nil {
					return err
				}
			}
		}
	} */

	json.NewEncoder(os.Stdout).Encode(entry)

	return nil
}

func (mw *MultiWriter[T]) GetWriters() []LogWriter[T] { return mw.writers }
