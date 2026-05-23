package manager

import (
	"context"
	"io"
	"sync"

	"github.com/kubex-ecosystem/logz"
	"github.com/kubex-ecosystem/logz/internal/manager/control"
	"github.com/kubex-ecosystem/logz/internal/manager/events"
)

// Manager handles the processing pipeline for log entries.
type Manager struct {
	mu sync.RWMutex

	formatter    Formatter
	writer       io.Writer
	hooks        events.Collection[logz.Entry]
	levelEnabled func(logz.Level) bool

	entry *logz.Entry
	state control.JobState
}

// Formatter defines the log formatting contract.
type Formatter interface {
	Format(e *logz.Entry) ([]byte, error)
}

// IsTerminal checks if the manager is in a terminal state.
func (m *Manager) IsTerminal() bool {
	return m.state.IsTerminal()
}

// Process runs the entry through the processing pipeline.
func (m *Manager) Process(ctx context.Context, entry logz.Entry) error {
	if entry == nil {
		return nil
	}
	if m.IsTerminal() {
		return control.ErrTerminal
	}

	// ---- Stage 1: Validate --------------------------------------
	if err := m.stageValidate(entry); err != nil {
		m.state.Fail()
		return err
	}

	// ---- Stage 2: Pre-Hooks -------------------------------------
	if err := m.hooks.Fire(entry); err != nil {
		m.state.Fail()
		return err
	}

	// ---- Stage 3: Format ----------------------------------------
	b, err := m.stageFormat(entry)
	if err != nil {
		m.state.Fail()
		return err
	}

	// ---- Stage 4: Write -----------------------------------------
	if err := m.stageWrite(b); err != nil {
		m.state.Fail()
		return err
	}

	m.state.Complete()
	return nil
}
