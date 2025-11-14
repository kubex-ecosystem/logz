package interfaces

import "io"

// Writer defines the contract for writing logs.
type Writer interface {
	Write(entry []byte) (int, error)
	Out() io.Writer
}
