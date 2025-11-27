package interfaces

// FormatterG é genérico em T, mas T precisa ser um Record.
// Para o logz "normal", T será *Entry.
// Para outros módulos do Kubex, T pode ser outro tipo que implemente Record.
type FormatterG[T Entry] interface {
	Format(T) ([]byte, error)
}

// FormatterFuncG é uma função genérica que implementa a interface FormatterG.
type FormatterFuncG[T Entry] func(T) ([]byte, error)

// Format formata o record genérico.
func (f FormatterFuncG[T]) Format(rec T) ([]byte, error) {
	return f(rec)
}
