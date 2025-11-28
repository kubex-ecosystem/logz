// Package control provides atomic flag registries and helpers for
// job state and security flags.
package control

import "sync/atomic"

// FlagReg32 é um registrador genérico baseado em uint32 com operações
// atômicas para trabalhar com bitflags (T deve ser ~uint32).
type FlagReg32[T ~uint32] struct {
	v atomic.Uint32
}

// NewFlagReg32 cria um registrador novo, com valor inicial 0.
func NewFlagReg32[T ~uint32]() *FlagReg32[T] {
	return &FlagReg32[T]{}
}

// Load lê o valor atual do registrador.
func (r *FlagReg32[T]) Load() T {
	return T(r.v.Load())
}

// Store grava o valor bruto no registrador (cuidado: sobrescreve tudo).
func (r *FlagReg32[T]) Store(val T) {
	r.v.Store(uint32(val))
}

// CompareAndSwap executa um CAS no valor inteiro do registrador.
func (r *FlagReg32[T]) CompareAndSwap(old, new T) bool {
	return r.v.CompareAndSwap(uint32(old), uint32(new))
}

// Set faz OR bitwise com a máscara, de forma atômica.
func (r *FlagReg32[T]) Set(mask T) {
	for {
		old := r.v.Load()
		newV := old | uint32(mask)
		if r.v.CompareAndSwap(old, newV) {
			return
		}
	}
}

// Clear faz AND-NOT bitwise com a máscara, de forma atômica.
func (r *FlagReg32[T]) Clear(mask T) {
	for {
		old := r.v.Load()
		newV := old &^ uint32(mask)
		if r.v.CompareAndSwap(old, newV) {
			return
		}
	}
}

// Has retorna true se TODOS os bits de mask estiverem presentes.
func (r *FlagReg32[T]) Has(mask T) bool {
	v := r.v.Load()
	return v&uint32(mask) == uint32(mask)
}

// Any retorna true se QUALQUER bit de mask estiver presente.
func (r *FlagReg32[T]) Any(mask T) bool {
	v := r.v.Load()
	return v&uint32(mask) != 0
}

// All é alias semântico pra Has.
func (r *FlagReg32[T]) All(mask T) bool {
	return r.Has(mask)
}

// Alias de compatibilidade: se em algum lugar você usou FlagReg32A/ NewFlagReg32A,
// continua funcionando sem quebrar nada.
type FlagReg32A[T ~uint32] = FlagReg32[T]

func NewFlagReg32A[T ~uint32]() *FlagReg32A[T] {
	return NewFlagReg32[T]()
}
