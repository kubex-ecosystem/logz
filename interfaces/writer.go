package interfaces

type Writer interface {
	Write([]byte) (int, error)
}
