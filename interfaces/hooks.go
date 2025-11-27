package interfaces

import "fmt"

type LHook[T any] interface {
	Fire(record T) error
	String() string
	Clone() LHook[T]
	Type() T
}

type FHook[T any] func(record T) error

func (f FHook[T]) Fire(record T) error {
	return f(record)
}

func (f FHook[T]) String() string {
	return "FHook"
}

func (f FHook[T]) Clone() FHook[T] {
	return f
}

type Hook func(record Entry) error

type Hooks []Hook

// Add adiciona um hook à coleção.
func (h Hooks) Add(hook Hook) (Hooks, error) {
	if h == nil {
		h = make(Hooks, 0)
	}
	if hook == nil {
		return nil, fmt.Errorf("hook is nil")
	} else {
		newHook := true
		if len(h) > 0 {
			for _, hkk := range h {
				if hkk != nil {
					if &hkk == &hook {
						newHook = false
						return h, fmt.Errorf("hook already exists in collection")
					}
				}
			}
		}
		if newHook {
			h = append(h, hook)
		} else {
			return h, fmt.Errorf("hook already exists in collection")
		}
	}
	return h, nil
}

// Fire executa todos os hooks da coleção.
func (h Hooks) Fire(record Entry) error {
	for _, hook := range h {
		err := hook(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// HookG é executado antes da formatação/escrita.
// Pode enriquecer o record, coletar métricas, enviar para outro sistema, etc.
type HookG[T any] func(T) error

// HooksG é uma coleção de hooks.
type HooksG[T any] []HookG[T]

// Add adiciona um hook à coleção.
func (h *HooksG[T]) Add(hook HookG[T]) error {
	if hook == nil {
		return fmt.Errorf("hook is nil")
	}
	*h = append(*h, hook)
	return nil
}

// Fire executa todos os hooks da coleção.
func (h HooksG[T]) Fire(record T) error {
	for _, hook := range h {
		err := hook(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// HookFunc é uma função que implementa a interface Hook.
type HookFunc func(record Entry) error

// Run executa o hook.
func (f HookFunc) Run(record Entry) error {
	return f(record)
}

// HookFuncG é uma função genérica que implementa a interface HookG.
type HookFuncG[T Entry] func(record T) error

// Run executa o hook genérico.
func (f HookFuncG[T]) Run(record T) error {
	return f(record)
}
