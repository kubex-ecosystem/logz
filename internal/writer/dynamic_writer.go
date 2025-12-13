package writer

import (
	"io"
	"sync"
)

// DynamicWriter permite trocar o destino em runtime.
type DynamicWriter struct {
	mu     sync.RWMutex
	target LogzWriter
}

func NewDynamicWriter(initial Writer) LogzWriter {
	return NewDynamicWriterType(initial)
}

func NewDynamicWriterType(initial Writer) *DynamicWriter {
	return &DynamicWriter{target: initial.(LogzWriter)}
}

func (d *DynamicWriter) Set(w LogzWriter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.target = w
}

func (d *DynamicWriter) Write(b []byte) (int, error) {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return 0, nil
	}
	return t.Write(b)
}

func (d *DynamicWriter) WriteLogz(b []byte) error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	_, err := t.Write(b)
	return err
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

func (d *DynamicWriter) GetIOWriter() io.Writer {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.GetIOWriter()
}
func (d *DynamicWriter) SetOutput(w io.Writer) {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return
	}
	t.SetOutput(w)
}

func (d *DynamicWriter) GetOutput() io.Writer {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.GetOutput()
}
func (d *DynamicWriter) Sync() error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.Sync()
}
func (d *DynamicWriter) String() string {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return "DynamicWriter(nil)"
	}
	return "DynamicWriter(" + t.String() + ")"
}
