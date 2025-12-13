package writer

import "io"

type MultiWriter struct {
	writers []LogzWriter
}

func NewMultiWriter(writers ...Writer) LogzWriter {
	return NewMultiWriterType(writers...)
}
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

func (m *MultiWriter) Close() error {
	var lastErr error
	for _, w := range m.writers {
		if err := w.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (m *MultiWriter) String() string {
	return "MultiWriter"
}

func (m *MultiWriter) GetIOWriter() io.Writer {
	return m
}

func (m *MultiWriter) SetOutput(w io.Writer) {
	// Not implemented for MultiWriter
}

func (m *MultiWriter) GetOutput() io.Writer {
	return m
}

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
func (m *MultiWriter) WriteLogz(b []byte) error {
	return m.LogzWrite(b)
}
func (m *MultiWriter) AddWriter(w LogzWriter) {
	m.writers = append(m.writers, w)
}

func (m *MultiWriter) RemoveWriter(w LogzWriter) {
	for i, writer := range m.writers {
		if writer == w {
			m.writers = append(m.writers[:i], m.writers[i+1:]...)
			break
		}
	}
}
func (m *MultiWriter) Writers() []LogzWriter {
	return m.writers
}
func (m *MultiWriter) Count() int {
	return len(m.writers)
}
func (m *MultiWriter) IsEmpty() bool {
	return len(m.writers) == 0
}
func (m *MultiWriter) Clear() {
	m.writers = []LogzWriter{}
}
func (m *MultiWriter) GetWriters() []LogzWriter {
	return m.writers
}
func (m *MultiWriter) SetWriters(writers []LogzWriter) {
	m.writers = writers
}
func (m *MultiWriter) GetWriterAt(index int) LogzWriter {
	if index < 0 || index >= len(m.writers) {
		return nil
	}
	return m.writers[index]
}
func (m *MultiWriter) SetWriterAt(index int, w LogzWriter) {
	if index < 0 || index >= len(m.writers) {
		return
	}
	m.writers[index] = w
}
