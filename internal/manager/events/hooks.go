// Package events provides atomic event hooks for the Kubex ecosystem.
package events

import (
	"fmt"
	"sync"
)

// Hook defines the contract for a generic event consumer.
type Hook[T any] interface {
	Fire(T) error
}

// HookFunc allows plain functions to act as Hooks.
type HookFunc[T any] func(T) error

func (f HookFunc[T]) Fire(record T) error {
	return f(record)
}

// Collection is a thread-safe set of hooks.
type Collection[T any] struct {
	mu    sync.RWMutex
	hooks []Hook[T]
}

// Add adds a new hook to the collection.
func (c *Collection[T]) Add(h Hook[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hooks = append(c.hooks, h)
}

// Fire executes all hooks in the collection.
func (c *Collection[T]) Fire(record T) error {
	c.mu.RLock()
	// Copy slice to avoid holding lock during execution
	hooks := append([]Hook[T](nil), c.hooks...)
	c.mu.RUnlock()

	for _, h := range hooks {
		if err := h.Fire(record); err != nil {
			return fmt.Errorf("hook execution failed: %w", err)
		}
	}
	return nil
}

// Global registry example or shortcut
type LogHook = Hook[any]
type LogHooks = Collection[any]
