// Package control - job state machine baseada em bitflags para uso interno
// no pipeline do logger / manager.
package control

import (
	kbx "github.com/kubex-ecosystem/logz/internal/module/kbx"
)

/* ========= JOB STATE FLAGS ========= */

type JobFlag uint32

const (
	JobPendingA JobFlag = 1 << iota
	JobRunningA
	JobCancelRequestedA
	JobRetryingA
	JobCompletedA
	JobFailedA
	JobTimedOutA
)

const terminalMask JobFlag = JobCompletedA | JobFailedA | JobTimedOutA

func (j JobFlag) Has(flag JobFlag) bool {
	return j&flag != 0
}

/* ========= JOB STATE (FSM) ========= */

// JobState encapsula a evolução do job via bitflags, usando FlagReg32[JobFlag].
type JobState struct {
	r *FlagReg32[JobFlag]
}

// NewJobState cria um JobState novo, iniciando em zero (sem flags).
func NewJobState() *JobState {
	return &JobState{r: NewFlagReg32[JobFlag]()}
}

// Load retorna os flags atuais.
func (s *JobState) Load() JobFlag {
	return s.r.Load()
}

// IsTerminal retorna true se o job já atingiu um estado terminal
// (Completed, Failed ou TimedOut).
func (s *JobState) IsTerminal() bool {
	return s.r.Any(terminalMask)
}

/* ========= TRANSIÇÕES ========= */

// Start: permite transicionar para Running se não estiver em estado terminal.
// Se já estiver terminal, retorna ErrTerminal.
func (s *JobState) Start() error {
	for {
		old := s.r.Load()
		// não inicia se já terminal
		if old&terminalMask != 0 {
			return kbx.ErrTerminal
		}
		// já está rodando? nada a fazer
		if old&JobRunningA != 0 {
			return nil
		}
		// Running liga, e garantimos que flags terminais não coexistam
		newV := (old | JobRunningA) &^ (JobCompletedA | JobFailedA | JobTimedOutA)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

// RequestCancel seta o flag de cancelamento solicitado.
func (s *JobState) RequestCancel() {
	s.r.Set(JobCancelRequestedA)
}

// Retry: marca Retry e limpa Running, se ainda não terminal.
func (s *JobState) Retry() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return kbx.ErrTerminal
		}
		newV := (old | JobRetryingA) &^ JobRunningA
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

// Complete: marca Completed e limpa Running/Retrying/CancelRequested.
// Se já terminal, retorna ErrTerminal.
func (s *JobState) Complete() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return kbx.ErrTerminal
		}
		newV := (old | JobCompletedA) &^ (JobRunningA | JobRetryingA | JobCancelRequestedA)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

// Fail: marca Failed e limpa Running/Retrying, se ainda não terminal.
func (s *JobState) Fail() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return kbx.ErrTerminal
		}
		newV := (old | JobFailedA) &^ (JobRunningA | JobRetryingA)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

// Timeout: marca TimedOut e limpa Running/Retrying, se ainda não terminal.
func (s *JobState) Timeout() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return kbx.ErrTerminal
		}
		newV := (old | JobTimedOutA) &^ (JobRunningA | JobRetryingA)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}
