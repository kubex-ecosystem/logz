package logz

// Writer recebe bytes já formatados e empurra pra algum destino.
// NÃO sabe nada sobre Entry.
type Writer interface {
	Write([]byte) error
	Close() error
}
