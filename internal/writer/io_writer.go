package logz

import (
	"io"
	"sync"
)

type IOWriter struct {
	w  io.Writer
	mu sync.Mutex
}

func NewIOWriter(w io.Writer) *IOWriter {
	return &IOWriter{w: w}
}

func (w *IOWriter) Write(b []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := w.w.Write(b)
	return err
}

func (w *IOWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
