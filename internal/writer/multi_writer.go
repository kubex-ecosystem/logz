package writer

type MultiWriter struct {
	writers []Writer
}

func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

func (m *MultiWriter) Write(b []byte) error {
	var lastErr error
	for _, w := range m.writers {
		if err := w.Write(b); err != nil {
			lastErr = err
		}
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
