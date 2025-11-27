package logz

import "sync"

// DynamicWriter permite trocar o destino em runtime.
type DynamicWriter struct {
	mu     sync.RWMutex
	target Writer
}

func NewDynamicWriter(initial Writer) *DynamicWriter {
	return &DynamicWriter{target: initial}
}

func (d *DynamicWriter) Set(w Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.target = w
}

func (d *DynamicWriter) Write(b []byte) error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.Write(b)
}

func (d *DynamicWriter) Close() error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.Close()
}
