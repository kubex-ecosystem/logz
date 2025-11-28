// Package manager implements the processing manager for log entries.
package manager

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	control "github.com/kubex-ecosystem/logz/internal/manager/control"
	// . "github.com/kubex-ecosystem/logz/internal/manager/control"

	// . "github.com/kubex-ecosystem/logz/internal/module/kbx"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// Entry aqui é o seu tipo concreto.
// No seu código, você tem *Entry, então sigo essa linha.
type Manager struct {
	mu sync.RWMutex

	formatter    Formatter               // normalmente l.opts.Formatter
	writer       io.Writer               // normalmente l.Writer()
	hooks        []interfaces.LHook[any] // normalmente l.opts.LHooks
	levelEnabled func(kbx.Level) bool    // ponte pro Enabled do logger (por enquanto)

	stage atomic.Uint32
	state atomic.Uint32
	entry *core.Entry

	options *kbx.InitArgs
	ctl     control.ManagerControl
}

// Formatter é o contrato já existente no seu logger.
type Formatter interface {
	Format(e *core.Entry) ([]byte, error)
}

// IsTerminal verifica se o manager está em estado terminal (done ou failed).
func (m *Manager) IsTerminal() bool {
	return m.ctl.State.Any(control.StepTerminal)
}

// Process é o pipeline principal para ENTRIES saudáveis/sóbrios.
// Logger NÃO faz mais formatação, hooks, write — só chama isso aqui.
func (m *Manager) Process(ctx context.Context, entry *core.Entry) error {
	if entry == nil {
		return nil
	}
	if m.IsTerminal() {
		return kbx.ErrTerminal
	}

	// ---- Stage 1: validate/level gate --------------------------------------
	if err := m.stageValidate(entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepValidate)

	// ---- Stage 2: pre-hooks -------------------------------------------------
	if err := m.stagePreHooks(ctx, entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepPreHooks)

	// ---- Stage 3: format ----------------------------------------------------
	b, err := m.stageFormat(entry)
	if err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepFormat)

	// ---- Stage 4: post-hooks ------------------------------------------------
	if err := m.stagePostHooks(ctx, entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepPostHooks)

	// ---- Stage 5: write -----------------------------------------------------
	if err := m.stageWrite(b); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepWrite)

	m.advance(control.StepDone)
	return nil
}
