package interfaces

type LHook[T any] interface {
	Fire(record T)
	String() string
	Clone() LHook[T]
	Type() T
}

type FHook[T any] func(record T)

func (f FHook[T]) Fire(record T) {
	f(record)
}

func (f FHook[T]) String() string {
	return "FHook"
}

func (f FHook[T]) Clone() FHook[T] {
	return f
}

type Hook func(record Entry)

type Hooks []Hook

// Add adiciona um hook à coleção.
func (h Hooks) Add(hook Hook) {
	h = append(h, hook)
}

// Fire executa todos os hooks da coleção.
func (h Hooks) Fire(record Entry) {
	for _, hook := range h {
		hook(record)
	}
}

// HookG é executado antes da formatação/escrita.
// Pode enriquecer o record, coletar métricas, enviar para outro sistema, etc.
type HookG[T any] func(T)

// HooksG é uma coleção de hooks.
type HooksG[T any] []HookG[T]

// Add adiciona um hook à coleção.
func (h *HooksG[T]) Add(hook HookG[T]) {
	*h = append(*h, hook)
}

// Fire executa todos os hooks da coleção.
func (h HooksG[T]) Fire(record T) {
	for _, hook := range h {
		hook(record)
	}
}

// HookFunc é uma função que implementa a interface Hook.
type HookFunc func(record Entry)

// Run executa o hook.
func (f HookFunc) Run(record Entry) {
	f(record)
}

// HookFuncG é uma função genérica que implementa a interface HookG.
type HookFuncG[T Entry] func(record T)

// Run executa o hook genérico.
func (f HookFuncG[T]) Run(record T) {
	f(record)
}
