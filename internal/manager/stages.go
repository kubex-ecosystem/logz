package manager

import (
	"context"
	"errors"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

func (m *Manager) stageValidate(entry *core.Entry) error {
	if m.levelEnabled == nil {
		// se ninguém configurou levelEnabled, deixa passar.
		return nil
	}
	lvl := kbx.Level(entry.GetLevel().String())
	if !m.levelEnabled(lvl) {
		// log ignorado silenciosamente, igual seu comportamento atual.
		return kbx.ErrTerminal // ou nil, se quiser só "não logar"
	}
	return nil
}

func (m *Manager) stagePreHooks(ctx context.Context, entry *core.Entry) error {
	m.mu.RLock()
	hooks := append([]interfaces.LHook[any](nil), m.hooks...)
	m.mu.RUnlock()

	for _, h := range hooks {
		if h == nil {
			continue
		}
		// Assumindo que sua interface é algo como Fire(any) error;
		// se quiser contexto, pode adaptar a interface pra aceitar ctx.
		if err := h.Fire(entry); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) stageFormat(entry *core.Entry) ([]byte, error) {
	m.mu.RLock()
	f := m.formatter
	m.mu.RUnlock()

	if f == nil {
		return nil, errors.New("logz: no formatter configured in Manager")
	}

	b, err := f.Format(entry)
	if err != nil {
		return nil, err
	}

	// garante newline pra saída de console / arquivos de texto.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return b, nil
}

func (m *Manager) stagePostHooks(ctx context.Context, entry kbx.Entry) error {
	// hoje você não diferencia pre/post de verdade;
	// aqui você tem liberdade de, no futuro, separar hooks por fase.
	m.mu.RLock()
	hooks := append([]interfaces.LHook[any](nil), m.hooks...)
	m.mu.RUnlock()

	for _, h := range hooks {
		if h == nil {
			continue
		}
		if err := h.Fire(entry); err != nil {
			return err
		}
	}
	return nil
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
