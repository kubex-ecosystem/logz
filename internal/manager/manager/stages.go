package manager

import (
	"errors"

	"github.com/kubex-ecosystem/logz"
)

func (m *Manager) stageValidate(entry logz.Entry) error {
	if m.levelEnabled == nil {
		return nil
	}
	lvl := logz.Level(entry.GetLevel().String())
	if !m.levelEnabled(lvl) {
		return errors.New("level not enabled")
	}
	return nil
}

func (m *Manager) stageFormat(entry logz.Entry) ([]byte, error) {
	m.mu.RLock()
	f := m.formatter
	m.mu.RUnlock()

	if f == nil {
		return nil, errors.New("logz: no formatter configured in Manager")
	}

	b, err := f.Format(&entry)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return b, nil
}

func (m *Manager) stageWrite(b []byte) error {
	m.mu.RLock()
	out := m.writer
	m.mu.RUnlock()

	if out == nil {
		return errors.New("logz: no writer configured in Manager")
	}

	_, err := out.Write(b)
	return err
}
