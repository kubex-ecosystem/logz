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

	Error error `json:"error,omitempty"` // erro associado (se houver)
}

func NewEntryImpl() *Entry {
	return &Entry{
		Timestamp: time.Now().UTC(),
		Tags:      make(map[string]string),
		Fields:    make(map[string]any),
		Caller:    captureCaller(3),
	}
}

// NewEntry cria uma entry com:
// - timestamp UTC
// - maps inicializados
// - caller capturado
func NewEntry() interfaces.Entry {
	return NewEntryImpl()
}

//
// ---------- Chainable builders ----------
//

func (e *Entry) WithLevel(l interfaces.Level) interfaces.Entry {
	e.Level = l
	e.Severity = l.Severity()
	return e
}

func (e *Entry) WithMessage(msg string) interfaces.Entry {
	e.Message = msg
	return e
}

func (e *Entry) WithContext(ctx string) interfaces.Entry {
	e.Context = ctx
	return e
}

func (e *Entry) WithSource(src string) interfaces.Entry {
	e.Source = src
	return e
}

func (e *Entry) WithTraceID(id string) interfaces.Entry {
	e.TraceID = id
	return e
}

func (e *Entry) WithField(key string, value any) interfaces.Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[key] = value
	return e
}

func (e *Entry) WithFields(fields map[string]any) interfaces.Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func (e *Entry) WithData(data any) interfaces.Entry {
	e.Fields["data"] = data
	return e
}

func (e *Entry) WithError(err error) interfaces.Entry {
	e.Error = err
	return e
}

func (e *Entry) Tag(k, v string) interfaces.Entry {
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}
	e.Tags[k] = v
	return e
}

func (e *Entry) Field(k string, v any) interfaces.Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[k] = v
	return e
}

func (e *Entry) WithCaller(c string) interfaces.Entry {
	e.Caller = c
	return e
}

func (e *Entry) CaptureCaller(skip int) interfaces.Entry {
	e.Caller = captureCaller(skip + 1)
	return e
}

func (e *Entry) GetTimestamp() time.Time {
	if e == nil {
		return time.Time{}
	}
	return e.Timestamp
}

func (e *Entry) GetContext() string {
	if e == nil {
		return ""
	}
	return e.Context
}

func (e *Entry) GetCaller() string {
	if e == nil {
		return ""
	}
	return e.Caller
}

func (e *Entry) GetTags() map[string]string {
	if e == nil {
		return nil
	}
	return e.Tags
}

func (e *Entry) GetFields() map[string]any {
	if e == nil {
		return nil
	}
	return e.Fields
}

//
// ---------- Clone sem aliasing ----------
//

func (e *Entry) Clone() interfaces.Entry {
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

func (e *Entry) GetMessage() string {
	if e == nil {
		return ""
	}
	return e.Message
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
