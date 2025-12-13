// Package core provides the fundamental logging abstractions and implementations.
package core

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// Entry é a unidade básica de log do sistema.
// Tudo no Kubex que for "log estruturado" deveria conseguir ser expresso nisso.
type Entry struct {
	Timestamp time.Time `json:"ts" yaml:"ts" xml:"ts" mapstructure:"ts"`             // Sempre UTC.
	Level     kbx.Level `json:"level" yaml:"level" xml:"level" mapstructure:"level"` // debug / info / warn / error / fatal / silent
	Message   string    `json:"msg" yaml:"msg" xml:"msg" mapstructure:"msg"`         // Mensagem humana.

	ShowColor   bool   `json:"show_color,omitempty" yaml:"show_color,omitempty" xml:"show_color,omitempty" mapstructure:"show_color,omitempty"`             // Habilita cores na saída
	ShowIcon    bool   `json:"show_icon,omitempty" yaml:"show_icon,omitempty" xml:"show_icon,omitempty" mapstructure:"show_icon,omitempty"`                 // Habilita ícones na saída
	ShowTraceID bool   `json:"show_trace_id,omitempty" yaml:"show_trace_id,omitempty" xml:"show_trace_id,omitempty" mapstructure:"show_trace_id,omitempty"` // Habilita o ID de rastreamento na saída
	ShowCaller  bool   `json:"show_caller,omitempty" yaml:"show_caller,omitempty" xml:"show_caller,omitempty" mapstructure:"show_caller,omitempty"`         // Habilita informações do chamador na saída
	ShowStack   bool   `json:"show_stack,omitempty" yaml:"show_stack,omitempty" xml:"show_stack,omitempty" mapstructure:"show_stack,omitempty"`             // Habilita informações da pilha de chamadas na saída
	ShowFields  bool   `json:"show_fields,omitempty" yaml:"show_fields,omitempty" xml:"show_fields,omitempty" mapstructure:"show_fields,omitempty"`         // Habilita campos adicionais na saída
	Format      string `json:"format,omitempty" yaml:"format,omitempty" xml:"format,omitempty" mapstructure:"format,omitempty"`                             // json / text / xml / etc.

	Context  string `json:"ctx,omitempty" yaml:"ctx,omitempty" xml:"ctx,omitempty" mapstructure:"ctx,omitempty"`             // ex: "auth", "db", "billing"
	Source   string `json:"src,omitempty" yaml:"src,omitempty" xml:"src,omitempty" mapstructure:"src,omitempty"`             // componente/módulo/serviço
	TraceID  string `json:"trace,omitempty" yaml:"trace,omitempty" xml:"trace,omitempty" mapstructure:"trace,omitempty"`     // correlação
	Caller   string `json:"caller,omitempty" yaml:"caller,omitempty" xml:"caller,omitempty" mapstructure:"caller,omitempty"` // arquivo:linha função
	Severity int    `json:"sev,omitempty" yaml:"sev,omitempty" xml:"sev,omitempty" mapstructure:"sev,omitempty"`             // cache do Level.Severity()

	Tags   map[string]string `json:"tags,omitempty" yaml:"tags,omitempty" xml:"-" mapstructure:"tags,omitempty"`       // metadados arbitrários
	Fields map[string]any    `json:"fields,omitempty" yaml:"fields,omitempty" xml:"-" mapstructure:"fields,omitempty"` // dados estruturados arbitrários

	Error error `json:"error,omitempty"` // erro associado (se houver)
}

func NewKbxEntry(level kbx.Level) (kbx.LogzEntry, error) {
	return NewEntryImpl(level)
}

func NewEntry(level kbx.Level) (*Entry, error) {
	return &Entry{
		ShowColor:   true,
		ShowIcon:    true,
		ShowTraceID: false,
		Timestamp:   time.Now().UTC(),
		Tags:        make(map[string]string),
		Fields:      make(map[string]any),
		Caller:      captureCaller(3),
		Level:       level,
		Severity:    level.Severity(),
	}, nil
}

// NewEntryImpl cria uma entry com:
// - timestamp UTC
// - maps inicializados
// - caller capturado
func NewEntryImpl(level kbx.Level) (*Entry, error) {
	e, err := NewEntry(level)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func NewLogzEntry(level kbx.Level) kbx.LogzEntry {
	entry, err := NewKbxEntry(level)
	if err != nil {
		// Handle error by returning a default entry with level Info
		defaultEntry, _ := NewKbxEntry(kbx.LevelInfo)
		return defaultEntry
	}
	return entry
}

//
// ---------- Chainable builders ----------
//

func (e *Entry) WithLevel(l kbx.Level) kbx.LogzEntry {
	e.Level = l
	e.Severity = l.Severity()
	return e
}

func (e *Entry) WithTraceID(id string) kbx.LogzEntry {
	e.TraceID = id
	return e
}

func (e *Entry) WithColor(color bool) kbx.LogzEntry {
	e.ShowColor = color
	return e
}

func (e *Entry) WithIcon(icon bool) kbx.LogzEntry {
	e.ShowIcon = icon
	return e
}

func (e *Entry) WithMessage(msg string) kbx.LogzEntry {
	e.Message = msg
	return e
}

func (e *Entry) WithContext(ctx string) kbx.LogzEntry {
	e.Context = ctx
	return e
}

func (e *Entry) WithSource(src string) kbx.LogzEntry {
	e.Source = src
	return e
}

func (e *Entry) WithFormat(format string) kbx.LogzEntry {
	e.Format = format
	return e
}

func (e *Entry) WithField(key string, value any) kbx.LogzEntry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[key] = value
	return e
}

func (e *Entry) WithFields(fields map[string]any) kbx.LogzEntry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func (e *Entry) WithData(data any) kbx.LogzEntry {
	e.Fields["data"] = data
	return e
}

func (e *Entry) WithError(err error) kbx.LogzEntry {
	e.Error = err
	return e
}

func (e *Entry) Tag(k, v string) kbx.LogzEntry {
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}
	e.Tags[k] = v
	return any(e).(kbx.LogzEntry)
}

func (e *Entry) Field(k string, v any) kbx.Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[k] = v
	return any(e).(kbx.Entry)
}

func (e *Entry) WithCaller(c string) kbx.LogzEntry {
	// e.Caller = c
	frames := runtime.CallersFrames([]uintptr{uintptr(0)})
	for {
		framesLen, more := frames.Next()
		if more {
			// Se houver mais frames, pular este e pegar o próximo
			continue
		} else {
			// Último frame disponível
			e.Caller = fmt.Sprintf("%s:%d %s (byArg: %s)", framesLen.File, framesLen.Line, framesLen.Function, c)
			// Source é a FuncForPC
			e.Source = runtime.FuncForPC(framesLen.PC).Name()
			break
		}
	}

	return e
}

func (e *Entry) WithStack(show bool) kbx.LogzEntry {
	e.ShowStack = show
	return e
}

func (e *Entry) WithShowTraceID(show bool) kbx.LogzEntry {
	e.ShowTraceID = show
	return e
}

func (e *Entry) WithShowCaller(show bool) kbx.LogzEntry {
	e.ShowCaller = show
	return e
}

func (e *Entry) WithShowFields(show bool) kbx.LogzEntry {
	e.ShowFields = show
	return e
}

//
// ---------- Getters ----------
//

func (e *Entry) CaptureCaller(skip int) kbx.Entry {
	e.Caller = captureCaller(skip + 1)
	return any(e).(kbx.Entry)
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

func (e *Entry) GetShowColor() bool {
	if e == nil {
		return true
	}
	return e.ShowColor
}

func (e *Entry) GetShowStack() bool {
	if e == nil {
		return false
	}
	return e.ShowStack
}

func (e *Entry) GetTraceID() string {
	if e == nil {
		return ""
	}
	return e.TraceID
}

func (e *Entry) GetShowCaller() bool {
	if e == nil {
		return false
	}
	return e.ShowCaller
}

func (e *Entry) GetShowFields() bool {
	if e == nil {
		return false
	}
	return e.ShowFields
}

func (e *Entry) GetShowIcon() bool {
	if e == nil {
		return false
	}
	return e.ShowIcon
}

func (e *Entry) GetShowTraceID() bool {
	if e == nil {
		return false
	}
	return e.ShowTraceID
}

func (e *Entry) GetFormat() string {
	if e == nil {
		return ""
	}
	return e.Format
}

func (e *Entry) GetPrefix() string {
	if e == nil {
		return ""
	}
	return kbx.GetValueOrDefaultSimple(
		kbx.LoggerArgs.Prefix,
		"Logz",
	)
}

//
// ---------- Clone sem aliasing ----------
//

func (e *Entry) Clone() kbx.Entry {
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

func (e *Entry) GetLevel() kbx.Level {
	if e == nil {
		return kbx.LevelSilent
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
	if e.Timestamp.IsZero() || e.Timestamp.Location() != time.UTC {
		e.Timestamp = time.Now().UTC()
	}
	if len(strings.TrimSpace(string(e.Level))) == 0 {
		return errors.New("level is required")
	}
	if len(strings.TrimSpace(e.Message)) == 0 {
		return errors.New("message is required")
	}
	// Silent pode ter severidade 0.
	if e.Level != kbx.LevelSilent && e.Severity <= 0 {
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
