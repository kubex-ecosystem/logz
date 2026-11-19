package writer

import (
	"io"
	"sync"
)

// IOWriter é um wrapper de io.Writer que implementa LogzWriter.
type IOWriter struct {
	w  io.Writer
	mu sync.Mutex
}

// NewIOWriter cria um novo IOWriter.
func NewIOWriter(w io.Writer) *IOWriter {
	return &IOWriter{w: w}
}

// Write implementa a interface LogzWriter.
func (w *IOWriter) Write(b []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := w.w.Write(b)
	return err
}

// Close implementa a interface LogzWriter.
func (w *IOWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// String retorna o nome do writer.
func (w *IOWriter) String() string {
	return "IOWriter"
}
