package manager

import (
	"sync/atomic"

	control "github.com/kubex-ecosystem/logz/internal/manager/control"
	ctl "github.com/kubex-ecosystem/logz/internal/module/control"
)

func (m *Manager) advance(flag ctl.JobFlag) {
	m.stage.Store(uint32(flag))
}

func (m *Manager) in(flag ctl.JobFlag) bool {
	vl := m.state.Load()
	return ctl.JobFlag(atomic.LoadUint32(
		(*uint32)(&vl),
	))&flag != 0
}

func (m *Manager) fail(err error) error {
	m.advance(control.StepFailed)
	return err
}

func (m *Manager) done() error {
	m.advance(control.StepDone)
	return nil
}
