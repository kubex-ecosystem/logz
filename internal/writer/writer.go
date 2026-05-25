// Package writer implementa Writers para diferentes destinos.
package writer

import (
	"io"
	"os"
)

// Writer recebe bytes já formatados e empurra pra algum destino.
// NÃO sabe nada sobre Entry. (Clone da interface io.Writer)
type Writer interface {
	Write([]byte) (int, error)
	Close() error
}

type LogzWriter interface {
	Writer

	GetIOWriter() io.Writer
	SetOutput(io.Writer)
	GetOutput() io.Writer
	WriteLogz([]byte) error
	Sync() error
	String() string
}

// LogzWriterImpl é a implementação padrão de LogzWriter.
// Promove (herda) todos os métodos e propriedades de io.Writer.
type LogzWriterImpl struct {
	output io.Writer
}

// NewLogzWriter cria um novo LogzWriter.
func NewLogzWriter(w io.Writer) LogzWriter {
	return &LogzWriterImpl{
		output: w,
	}
}

// GetIOWriter retorna o io.Writer.
func (l *LogzWriterImpl) GetIOWriter() io.Writer {
	return l.output
}

// SetOutput define um novo io.Writer.
func (l *LogzWriterImpl) SetOutput(w io.Writer) {
	l.output = w
}

// GetOutput retorna o io.Writer.
func (l *LogzWriterImpl) GetOutput() io.Writer { return l.output }

// Write implementa a interface io.Writer.
func (l *LogzWriterImpl) Write(p []byte) (n int, err error) {
	return l.output.Write(p)
}

// WriteLogz implementa a interface LogzWriter.
func (l *LogzWriterImpl) WriteLogz(p []byte) error {
	_, err := l.output.Write(p)
	return err
}

// Sync implementa a interface LogzWriter.
func (l *LogzWriterImpl) Sync() error {
	if syncer, ok := l.output.(interface{ Sync() error }); ok {
		return syncer.Sync()
	}
	return nil
}

// Close implementa a interface LogzWriter.
func (l *LogzWriterImpl) Close() error {
	if closer, ok := l.output.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
func (l *LogzWriterImpl) String() string {
	return "LogzWriter"
}

// ParseWriter retorna um Writer baseado na string de output.
func ParseWriter(output string) LogzWriter {
	switch output {
	case "stdout":
		return NewLogzWriter(os.Stdout)
	case "stderr":
		return NewLogzWriter(os.Stderr)
	default:
		if file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			return NewLogzWriter(file)
		}
		return NewLogzWriter(io.Discard)
	}
}
