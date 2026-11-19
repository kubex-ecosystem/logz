// Package control provides atomic flag registers and utilities for managing
// sets of bit flags on top of basic integer types.
package control

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

// FlagReg32 is an atomic register for bit flags with an underlying uint32.
type FlagReg32[T ~uint32] struct{ v atomic.Uint32 }

func (r *FlagReg32[T]) Load() T          { return T(r.v.Load()) }
func (r *FlagReg32[T]) Store(x T)        { r.v.Store(uint32(x)) }
func (r *FlagReg32[T]) Has(f T) bool     { return r.v.Load()&uint32(f) == uint32(f) }
func (r *FlagReg32[T]) Any(mask T) bool  { return r.v.Load()&uint32(mask) != 0 }
func (r *FlagReg32[T]) None(mask T) bool { return r.v.Load()&uint32(mask) == 0 }
func (r *FlagReg32[T]) Set(f T)          { r.cas(func(old uint32) uint32 { return old | uint32(f) }) }
func (r *FlagReg32[T]) Clear(f T)        { r.cas(func(old uint32) uint32 { return old &^ uint32(f) }) }
func (r *FlagReg32[T]) Toggle(f T)       { r.cas(func(old uint32) uint32 { return old ^ uint32(f) }) }
func (r *FlagReg32[T]) Mask(mask T) T    { return T(r.v.Load() & uint32(mask)) }

func (r *FlagReg32[T]) SetMask(mask, value T) {
	r.cas(func(old uint32) uint32 { return (old &^ uint32(mask)) | (uint32(value) & uint32(mask)) })
}

func (r *FlagReg32[T]) CompareAndSwap(old, new T) bool {
	return r.v.CompareAndSwap(uint32(old), uint32(new))
}

func (r *FlagReg32[T]) cas(f func(old uint32) uint32) {
	for {
		o := r.v.Load()
		n := f(o)
		if r.v.CompareAndSwap(o, n) {
			return
		}
	}
}

func (r *FlagReg32[T]) SetIf(mustBeClear, set T) bool {
	for {
		old := r.v.Load()
		if old&uint32(mustBeClear) != 0 {
			return false
		}
		newV := old | uint32(set)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}

func (r *FlagReg32[T]) ClearIf(mustBeSet, clr T) bool {
	for {
		old := r.v.Load()
		if old&uint32(mustBeSet) != uint32(mustBeSet) {
			return false
		}
		newV := old &^ uint32(clr)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}

// JobFlag representa os estados do job.
type JobFlag uint32

const (
	JobPending JobFlag = 1 << iota
	JobRunning
	JobCancelRequested
	JobRetrying
	JobCompleted
	JobFailed
	JobTimedOut
)

// StepFlag representa estágios de um pipeline (ex: logger).
type StepFlag uint32

const (
	StepValidate StepFlag = 1 << iota
	StepPreHooks
	StepFormat
	StepPostHooks
	StepWrite
	StepDone
	StepFailed
)

const (
	terminalMask     JobFlag  = JobCompleted | JobFailed | JobTimedOut
	StepTerminalMask StepFlag = StepDone | StepFailed
)

var ErrTerminal = errors.New("job is in a terminal state")

// JobState é o estado do job baseado em FlagReg32.
type JobState struct{ r FlagReg32[JobFlag] }

func (s *JobState) Load() JobFlag      { return s.r.Load() }
func (s *JobState) Has(f JobFlag) bool { return s.r.Has(f) }

func (s *JobState) Start() error {
	ok := s.r.SetIf(terminalMask|JobRunning, JobRunning)
	if !ok {
		return ErrTerminal
	}
	return nil
}

func (s *JobState) Complete() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobCompleted) &^ (JobRunning | JobRetrying | JobCancelRequested)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) Fail() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobFailed) &^ (JobRunning | JobRetrying)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) IsTerminal() bool { return s.r.Any(terminalMask) }

// ManagerControl gerencia estágios e estados.
type ManagerControl struct {
	Stage FlagReg32[StepFlag]
	State FlagReg32[StepFlag]
}

func (c *ManagerControl) IsTerminal() bool {
	return c.State.Any(StepTerminalMask)
}

// FlagReg64 mirrors FlagReg32 for 64-bit sets.
type FlagReg64[T ~uint64] struct{ v atomic.Uint64 }

func (r *FlagReg64[T]) Load() T      { return T(r.v.Load()) }
func (r *FlagReg64[T]) Store(x T)    { r.v.Store(uint64(x)) }
func (r *FlagReg64[T]) Has(f T) bool { return r.v.Load()&uint64(f) == uint64(f) }
func (r *FlagReg64[T]) Set(f T)      { r.cas(func(o uint64) uint64 { return o | uint64(f) }) }
func (r *FlagReg64[T]) Clear(f T)    { r.cas(func(o uint64) uint64 { return o &^ uint64(f) }) }
func (r *FlagReg64[T]) cas(f func(uint64) uint64) {
	for {
		o := r.v.Load()
		n := f(o)
		if r.v.CompareAndSwap(o, n) {
			return
		}
	}
}

func FlagString[T ~uint32](val T, pairs map[string]T) string {
	if len(pairs) == 0 {
		return fmt.Sprintf("0x%X", uint32(val))
	}
	names := make([]string, 0, len(pairs))
	for name, bit := range pairs {
		if uint32(val)&uint32(bit) != 0 {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return "<none>"
	}
	return strings.Join(names, "|")
}
