// Package events provides events for Logz.
package events

import (
	"fmt"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// LHook defines the interface for a hook that can be fired with a record.
type LHook[T any] interface {
	Fire(record T) error
	String() string
	Clone() LHook[T]
	Type() T
}

// FHook is a function type that processes a record of type T.
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

// Hook is a function type that takes a kbx.Entry and returns an error.
type Hook func(record kbx.Entry) error

// Hooks is a collection of Hook functions.
type Hooks []Hook

// Add adiciona um hook à coleção.
func (h Hooks) Add(hook Hook) (Hooks, error) {
	if h == nil {
		h = make(Hooks, 0)
	}
	if hook != nil {
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
		return h, nil
	}
	return nil, fmt.Errorf("hook is nil")
}

// Fire executa todos os hooks da coleção.
func (h Hooks) Fire(record kbx.Entry) error {
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
func (f *HooksG[T]) Add(hook HookG[T]) error {
	if hook == nil {
		return fmt.Errorf("hook is nil")
	}
	*f = append(*f, hook)
	return nil
}

// Fire executa todos os hooks da coleção.
func (f HooksG[T]) Fire(record T) error {
	for _, hook := range f {
		err := hook(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// HookFunc é uma função que implementa a interface Hook.
type HookFunc func(record kbx.Entry) error

// Fire executa o hook.
func (f HookFunc) Fire(record kbx.Entry) error {
	return f(record)
}

// HookFuncG é uma função genérica que implementa a interface HookG.
type HookFuncG[T Hook | *kbx.Entry] func(record T) error

// Fire executa o hook genérico.
func (f HookFuncG[T]) Fire(record T) error {
	return f(record)
}
