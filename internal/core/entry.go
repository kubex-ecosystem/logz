// Package core provides the fundamental logging abstractions and implementations.
package core

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kubex-ecosystem/logz/interfaces"
)

// Entry é a unidade básica de log do sistema.
// Tudo no Kubex que for "log estruturado" deveria conseguir ser expresso nisso.
type Entry struct {
	Timestamp time.Time        `json:"ts"`    // Sempre UTC.
	Level     interfaces.Level `json:"level"` // debug / info / warn / error / fatal / silent
	Message   string           `json:"msg"`   // Mensagem humana.

	Context  string `json:"ctx,omitempty"`    // ex: "auth", "db", "billing"
	Source   string `json:"src,omitempty"`    // componente/módulo/serviço
	TraceID  string `json:"trace,omitempty"`  // correlação
	Caller   string `json:"caller,omitempty"` // arquivo:linha função
	Severity int    `json:"sev,omitempty"`    // cache do Level.Severity()

	Tags   map[string]string `json:"tags,omitempty"`   // chave/valor indexáveis
	Fields map[string]any    `json:"fields,omitempty"` // payload arbitrário
}

// NewEntry cria uma entry com:
// - timestamp UTC
// - maps inicializados
// - caller capturado
func NewEntry() *Entry {
	return &Entry{
		Timestamp: time.Now().UTC(),
		Tags:      make(map[string]string),
		Fields:    make(map[string]any),
		Caller:    captureCaller(3),
	}
}

//
// ---------- Chainable builders ----------
//

func (e *Entry) WithLevel(l interfaces.Level) *Entry {
	e.Level = l
	e.Severity = l.Severity()
	return e
}

func (e *Entry) WithMessage(msg string) *Entry {
	e.Message = msg
	return e
}

func (e *Entry) WithContext(ctx string) *Entry {
	e.Context = ctx
	return e
}

func (e *Entry) WithSource(src string) *Entry {
	e.Source = src
	return e
}

func (e *Entry) WithTraceID(id string) *Entry {
	e.TraceID = id
	return e
}

func (e *Entry) Tag(k, v string) *Entry {
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}
	e.Tags[k] = v
	return e
}

func (e *Entry) Field(k string, v any) *Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[k] = v
	return e
}

func (e *Entry) WithCaller(c string) *Entry {
	e.Caller = c
	return e
}

func (e *Entry) CaptureCaller(skip int) *Entry {
	e.Caller = captureCaller(skip + 1)
	return e
}

//
// ---------- Clone sem aliasing ----------
//

func (e *Entry) Clone() *Entry {
	if e == nil {
		return nil
	}

	clone := *e

	if e.Tags != nil {
		clone.Tags = make(map[string]string, len(e.Tags))
		for k, v := range e.Tags {
			clone.Tags[k] = v
		}
	}

	if e.Fields != nil {
		clone.Fields = make(map[string]any, len(e.Fields))
		for k, v := range e.Fields {
			clone.Fields[k] = v
		}
	}

	return &clone
}

//
// ---------- Record interface ----------
//

func (e *Entry) GetLevel() interfaces.Level {
	if e == nil {
		return interfaces.LevelSilent
	}
	return e.Level
}

//
// ---------- Sanity check ----------
//

func (e *Entry) Validate() error {
	if e == nil {
		return errors.New("entry is nil")
	}
	if e.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if strings.TrimSpace(string(e.Level)) == "" {
		return errors.New("level is required")
	}
	if strings.TrimSpace(e.Message) == "" {
		return errors.New("message is required")
	}
	// Silent pode ter severidade 0.
	if e.Level != interfaces.LevelSilent && e.Severity <= 0 {
		return errors.New("invalid severity (did you forget WithLevel?)")
	}
	return nil
}

//
// ---------- Debug-friendly String() ----------
//

func (e *Entry) String() string {
	if e == nil {
		return "<nil entry>"
	}
	return fmt.Sprintf(
		"%s [%s] %s",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Message,
	)
}

//
// ---------- Auxiliares internos ----------
//

func captureCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return fmt.Sprintf("%s:%d %s", file, line, fn.Name())
}
