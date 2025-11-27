// Package writer implementa Writers para diferentes destinos.
package writer

import (
	"io"
	"os"
)

// Writer recebe bytes já formatados e empurra pra algum destino.
// NÃO sabe nada sobre Entry.
type Writer interface {
	Write([]byte) error
	Close() error
}

func ParseWriter(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		if file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			return file
		}
		return io.Discard
	}
}
