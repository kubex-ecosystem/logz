// Package core contém funcionalidades centrais do logz.
package core

import (
	"fmt"
	"strings"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

func toEntry(level kbx.Level, args ...any) kbx.LogzEntry {
	if len(args) == 0 {
		en, err := NewEntry(level)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}

		return en.
			WithMessage("<empty>")
	}

	// Se já for Entry → retorna direto
	if e, ok := args[0].(*Entry); ok {
		return e
	}

	// Se for error
	if _, ok := args[0].(error); ok {
		en, err := NewEntry(kbx.LevelError)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}

		return en.
			WithError(err) /* .
			WithMessage(err.Error())
			WithLevel(kbx.LevelError) */
	}

	// Se for string
	if s, ok := args[0].(string); ok {
		en, err := NewEntry(level)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}

		return en.WithMessage(s)
	}

	// Se for []byte
	if b, ok := args[0].([]byte); ok {
		en, err := NewEntry(level)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}
		return en.WithMessage(string(b))
	}

	// Se for map
	if m, ok := args[0].(map[string]any); ok {
		en, err := NewEntry(level)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}

		return en.
			WithFields(m).
			WithMessage("map")
	}

	// Se for struct (fallback leve SEM panic)
	val := args[0]
	if kbx.IsObjSafe(val, false) {
		en, err := NewEntry(level)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}

		return en.
			WithMessage(fmt.Sprintf("%T", val)).
			WithData(val)
	}

	en, err := NewEntry(level)
	if err != nil {
		// fallback bruto
		return &Entry{
			Message: fmt.Sprintf("failed to create new entry: %v", err),
			Level:   kbx.LevelError,
		}
	}

	// fallback TOTAL
	return en.
		WithMessage(fmt.Sprintf("%v", args[0]))
}

func ToEntry(level kbx.Level, args ...any) *Entry {
	e, err := NewEntry(level)
	if err != nil {
		// fallback bruto
		return &Entry{
			Message: fmt.Sprintf("failed to create new entry: %v", err),
			Level:   kbx.LevelError,
		}
	}

	if len(args) == 0 {
		return e.WithMessage("<empty>").(*Entry)
	}

	first := args[0]

	// 1) se já é Entry
	if rec, ok := first.(*Entry); ok {
		return rec
	}

	// 2) se é erro
	if err, ok := first.(error); ok {
		e = e.WithError(err).(*Entry)
		if err != nil {
			// fallback bruto
			return &Entry{
				Message: fmt.Sprintf("failed to create new entry: %v", err),
				Level:   kbx.LevelError,
			}
		}
		if len(args) > 1 {
			e = e.WithField("args", args[1:]).(*Entry)
		}
		return e
	}

	// 3) string → mensagem
	if msg, ok := first.(string); ok {
		e = e.WithMessage(msg).(*Entry)
		if len(args) > 1 {
			// segundo arg error?
			if len(args) == 2 {
				if err, ok := args[1].(error); ok {
					return e.WithError(err).(*Entry)
				}
			}
			e = e.WithField("args", args[1:]).(*Entry)
		}
		return e
	}

	// 4) map → fields
	if m, ok := first.(map[string]any); ok {
		e = e.WithFields(m).(*Entry)
		if len(args) > 1 {
			e = e.WithField("args", args[1:]).(*Entry)
		}
		return e
	}

	// 5) []byte → mensagem
	if b, ok := first.([]byte); ok {
		e = e.WithMessage(string(b)).(*Entry)
		if len(args) > 1 {
			e = e.WithField("args", args[1:]).(*Entry)
		}
		return e
	}

	// 6) struct / qualquer coisa segura
	if kbx.IsObjSafe(first, false) {
		e = e.WithMessage(fmt.Sprintf("%T", first)).WithData(first).(*Entry)
		if len(args) > 1 {
			e = e.WithField("args", args[1:]).(*Entry)
		}
		return e
	}

	// 7) fallback
	return e.WithMessage(fmt.Sprintf("%v", first)).(*Entry)
}

func normalizeLevel(v any) kbx.Level {
	switch x := v.(type) {
	case kbx.Level:
		return x

	case string:
		l := kbx.Level(strings.ToLower(strings.TrimSpace(x)))
		if kbx.IsLevel(l.String()) {
			return l
		}
	}

	return kbx.LevelInfo
}
