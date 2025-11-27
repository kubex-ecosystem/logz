package interfaces

type Formatter interface {
	Format(e Entry) ([]byte, error)
}

// FormatterFunc é uma função que implementa a interface Formatter.
type FormatterFunc func(e Entry) ([]byte, error)

// Format formata a entry.
func (f FormatterFunc) Format(e Entry) ([]byte, error) {
	return f(e)
}
