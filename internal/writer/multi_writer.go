package writer

import "io"

// MultiWriter é um array de writers que recebem bytes já formatados e empurram para múltiplos destinos.
type MultiWriter struct {
	writers []LogzWriter
}

// NewMultiWriter cria um novo MultiWriter.
func NewMultiWriter(writers ...Writer) LogzWriter {
	return NewMultiWriterType(writers...)
}

// NewMultiWriterType cria um novo MultiWriter.
// Promove (herda) todos os métodos e propriedades de LogzWriter.
// Se o Writer não implementar LogzWriter, ele será promovido para LogzWriter.
func NewMultiWriterType(writers ...Writer) *MultiWriter {
	logzWriters := make([]LogzWriter, 0, len(writers))
	for _, w := range writers {
		if lw, ok := w.(LogzWriter); ok {
			logzWriters = append(logzWriters, lw)
		} else {
			logzWriters = append(logzWriters, NewLogzWriter(w))
		}
	}
	return &MultiWriter{writers: logzWriters}
}

// Write implementa a interface io.Writer.
func (m *MultiWriter) Write(b []byte) (n int, err error) {
	total := 0
	var lastErr error
	for _, w := range m.writers {
		n, err := w.Write(b)
		if err != nil {
			lastErr = err
		} else {
			total += n
		}
	}
	return total, lastErr
}

// LogzWrite implementa a interface LogzWriter.
func (m *MultiWriter) LogzWrite(b []byte) error {
	var lastErr error
	for _, w := range m.writers {
		// if lw, ok := w; ok {
		if err := w.WriteLogz(b); err != nil {
			lastErr = err
		}
		// }
	}
	return lastErr
}

// Close implementa a interface LogzWriter.
func (m *MultiWriter) Close() error {
	var lastErr error
	for _, w := range m.writers {
		if err := w.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// String retorna o nome do writer.
func (m *MultiWriter) String() string {
	return "MultiWriter"
}

// GetIOWriter retorna a instância de io.Writer do MultiWriter.
func (m *MultiWriter) GetIOWriter() io.Writer {
	return m
}

// SetOutput não é implementado para MultiWriter.
func (m *MultiWriter) SetOutput(w io.Writer) {
	// Not implemented for MultiWriter
}

// GetOutput retorna a instância de io.Writer do MultiWriter.
func (m *MultiWriter) GetOutput() io.Writer {
	return m
}

// Sync implementa a interface LogzWriter.
func (m *MultiWriter) Sync() error {
	var lastErr error
	for _, w := range m.writers {
		// if lw, ok := w.(LogzWriter); ok {
		if err := w.Sync(); err != nil {
			lastErr = err
		}
		// }
	}
	return lastErr
}

// WriteLogz implementa a interface LogzWriter.
func (m *MultiWriter) WriteLogz(b []byte) error {
	return m.LogzWrite(b)
}

// AddWriter adiciona um writer ao MultiWriter.
func (m *MultiWriter) AddWriter(w LogzWriter) {
	m.writers = append(m.writers, w)
}

// RemoveWriter remove um writer do MultiWriter.
func (m *MultiWriter) RemoveWriter(w LogzWriter) {
	for i, writer := range m.writers {
		if writer == w {
			m.writers = append(m.writers[:i], m.writers[i+1:]...)
			break
		}
	}
}

// Writers retorna todos os writers do MultiWriter.
func (m *MultiWriter) Writers() []LogzWriter {
	return m.writers
}

// Count retorna o número de writers do MultiWriter.
func (m *MultiWriter) Count() int {
	return len(m.writers)
}

// IsEmpty retorna true se o MultiWriter estiver vazio.
func (m *MultiWriter) IsEmpty() bool {
	return len(m.writers) == 0
}

// Clear limpa todos os writers do MultiWriter.
func (m *MultiWriter) Clear() {
	m.writers = []LogzWriter{}
}

// GetWriters retorna todos os writers do MultiWriter.
func (m *MultiWriter) GetWriters() []LogzWriter {
	return m.writers
}

// SetWriters define todos os writers do MultiWriter.
func (m *MultiWriter) SetWriters(writers []LogzWriter) {
	m.writers = writers
}

// GetWriterAt retorna um writer no índice especificado.
func (m *MultiWriter) GetWriterAt(index int) LogzWriter {
	if index < 0 || index >= len(m.writers) {
		return nil
	}
	return m.writers[index]
}

// SetWriterAt define um writer no índice especificado.
func (m *MultiWriter) SetWriterAt(index int, w LogzWriter) {
	if index < 0 || index >= len(m.writers) {
		return
	}
	m.writers[index] = w
}
