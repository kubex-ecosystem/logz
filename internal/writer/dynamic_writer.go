// Package writer implementa Writers para diferentes destinos.
package writer

import (
	"io"
	"sync"
)

// DynamicWriter é um wrapper de io.Writer que permite trocar o destino em runtime.
// Ele implementa todos os tipos de interfaces do Logz, permitindo que ele seja usado em qualquer lugar
// onde um LogzWriter é esperado. Apesar dele ser um wrapper/abstract, ele implementa realmente
// todas as interfaces do Logz, permitindo que ele possa atuar de fato como um LogzWriter genéricp.
// Por ele não possuir as propriedades que alguns writers possuem, ele é capaz de atuar, porém
// executando de uma forma simplificada o que alguns writers fazem nativamente (como o MultiWriter).
type DynamicWriter struct {
	mu     sync.RWMutex
	target LogzWriter
}

// NewDynamicWriter cria um novo DynamicWriter.
func NewDynamicWriter(initial Writer) LogzWriter {
	return NewDynamicWriterType(initial)
}

// NewDynamicWriterType cria um novo DynamicWriter.
func NewDynamicWriterType(initial Writer) *DynamicWriter {
	return &DynamicWriter{target: initial.(LogzWriter)}
}

// Set define um novo LogzWriter para o DynamicWriter.
func (d *DynamicWriter) Set(w LogzWriter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.target = w
}

// Write implementa a interface io.Writer.
func (d *DynamicWriter) Write(b []byte) (int, error) {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return 0, nil
	}
	return t.Write(b)
}

// WriteLogz implementa a interface LogzWriter.
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

// Close implementa a interface LogzWriter.
func (d *DynamicWriter) Close() error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.Close()
}

// GetIOWriter retorna a instância de io.Writer do DynamicWriter.
func (d *DynamicWriter) GetIOWriter() io.Writer {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.GetIOWriter()
}

// SetOutput define um novo io.Writer para o DynamicWriter.
func (d *DynamicWriter) SetOutput(w io.Writer) {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return
	}
	t.SetOutput(w)
}

// GetOutput retorna a instância de io.Writer do DynamicWriter.
func (d *DynamicWriter) GetOutput() io.Writer {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.GetOutput()
}

// Sync implementa a interface LogzWriter.
func (d *DynamicWriter) Sync() error {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return nil
	}
	return t.Sync()
}

// String retorna o nome do writer.
func (d *DynamicWriter) String() string {
	d.mu.RLock()
	t := d.target
	d.mu.RUnlock()
	if t == nil {
		return "DynamicWriter(nil)"
	}
	return "DynamicWriter(" + t.String() + ")"
}
