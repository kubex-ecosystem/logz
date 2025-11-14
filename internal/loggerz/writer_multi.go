package loggerz

import (
	"io"
)

type LogzWriter[E any, T WriteLogz[E]] interface {
	Write(entry []byte) (int, error)
	WriteLogz(entry E, ioWriter T) (int, error)
}

// MultiWriter is a writer that can write to multiple LogWriters.
// It implements the LogMultiWriter interface.
// It allows adding new writers dynamically and writing log entries to all of them.
// It can be used to aggregate logs from different sources or formats.
type MultiWriter[T any] struct {
	writers []LogzWriter[T, WriteLogz[T]]
	out     io.Writer
}

// NewMultiWriter creates a new instance of MultiWriter with the provided writers.
func NewMultiWriter[T any](writers ...LogzWriter[T, WriteLogz[T]]) *MultiWriter[T] {
	return &MultiWriter[T]{writers: writers}
}

func (mw *MultiWriter[T]) AddWriter(w LogzWriter[T, WriteLogz[T]]) {
	mw.writers = append(mw.writers, w)
}

func (mw *MultiWriter[T]) Write(entry []byte) (int, error) {
	var total int
	for _, w := range mw.writers {
		n, err := w.Write(entry)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (mw *MultiWriter[T]) WriteLog(entry T, ioWriter WriteLogz[T]) (int, error) {
	var total int
	for _, w := range mw.writers {
		n, err := w.WriteLogz(entry, ioWriter)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (mw *MultiWriter[T]) GetWriters() []LogzWriter[T, WriteLogz[T]] { return mw.writers }

func (mw *MultiWriter[T]) Out() io.Writer {
	return mw.out
}
