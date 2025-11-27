// Package core contém funcionalidades centrais do logz.
package core

import (
	"fmt"
	"strings"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

func toEntry(args ...any) interfaces.Entry {
	if len(args) == 0 {
		return NewEntry().WithMessage("<empty>")
	}

	// Se já for Entry → retorna direto
	if e, ok := args[0].(interfaces.Entry); ok {
		return e
	}

	// Se for error
	if err, ok := args[0].(error); ok {
		return NewEntry().
			WithError(err).
			WithMessage(err.Error()).
			WithLevel(interfaces.LevelError)
	}

	// Se for string
	if s, ok := args[0].(string); ok {
		return NewEntry().WithMessage(s)
	}

	// Se for []byte
	if b, ok := args[0].([]byte); ok {
		return NewEntry().WithMessage(string(b))
	}

	// Se for map
	if m, ok := args[0].(map[string]any); ok {
		return NewEntry().
			WithFields(m).
			WithMessage("map")
	}

	// Se for struct (fallback leve SEM panic)
	val := args[0]
	if kbx.IsObjSafe(val, false) {
		return NewEntry().
			WithMessage(fmt.Sprintf("%T", val)).
			WithData(val)
	}

	// fallback TOTAL
	return NewEntry().
		WithMessage(fmt.Sprintf("%v", args[0]))
}

func ToEntry(args ...any) interfaces.Entry {
	e := NewEntry()

	if len(args) == 0 {
		return e.WithMessage("<empty>")
	}

	first := args[0]

	// 1) se já é Entry
	if rec, ok := first.(interfaces.Entry); ok {
		return rec
	}

	// 2) se é erro
	if err, ok := first.(error); ok {
		e = e.WithError(err).WithMessage(err.Error())
		if len(args) > 1 {
			e = e.WithField("args", args[1:])
		}
		return e
	}

	// 3) string → mensagem
	if msg, ok := first.(string); ok {
		e = e.WithMessage(msg)
		if len(args) > 1 {
			// segundo arg error?
			if len(args) == 2 {
				if err, ok := args[1].(error); ok {
					return e.WithError(err)
				}
			}
			e = e.WithField("args", args[1:])
		}
		return e
	}

	// 4) map → fields
	if m, ok := first.(map[string]any); ok {
		e = e.WithFields(m)
		if len(args) > 1 {
			e = e.WithField("args", args[1:])
		}
		return e
	}

	// 5) []byte → mensagem
	if b, ok := first.([]byte); ok {
		e = e.WithMessage(string(b))
		if len(args) > 1 {
			e = e.WithField("args", args[1:])
		}
		return e
	}

	// 6) struct / qualquer coisa segura
	if kbx.IsObjSafe(first, false) {
		e = e.WithMessage(fmt.Sprintf("%T", first)).WithData(first)
		if len(args) > 1 {
			e = e.WithField("args", args[1:])
		}
		return e
	}

	// 7) fallback
	return e.WithMessage(fmt.Sprintf("%v", first))
}

func normalizeLevel(v any) interfaces.Level {
	switch x := v.(type) {
	case interfaces.Level:
		return x

	case string:
		l := interfaces.Level(strings.ToLower(strings.TrimSpace(x)))
		if interfaces.IsLevel(l.String()) {
			return l
		}
	}

	return interfaces.LevelInfo
}
