// Package core provides the fundamental logging abstractions and implementations.
package core

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// Entry é a unidade básica de log do sistema. E a maior diferença estrutural
// entre o logz e o stdlog. Ele encaminha tudo no fim para o stdlog? Não.
// Ele é capaz de lidar e trabalhar com outputs customizados, que são interfaces de io.Writer.
// Porém, ele conta com fmt para parsear algumas strings e formatar algumas saídas (não são todas: vide package writer)
// Tudo no Kubex que for "log estruturado" deveria conseguir ser expresso nisso. Se não vier nesse formato,
// nós fazemos um parse para transformar nessa estrutura. Ela é a "vida" do sistema.
// Tudo aqui deve ser serializável, para que possamos enviar para um output customizado, se não for, fazemos
// de forma que seja o mais compatível possível.
type Entry struct {
	ctx context.Context `json:"-" yaml:"-" xml:"-"` // Atribuído na criação, no TraceID, e usado no Writer (output)

	Timestamp time.Time `json:"ts" yaml:"ts" xml:"ts" mapstructure:"ts"`
	Level     kbx.Level `json:"level" yaml:"level" xml:"level" mapstructure:"level"`
	Message   string    `json:"msg" yaml:"msg" xml:"msg" mapstructure:"msg"`

	ShowColor   bool `json:"show_color,omitempty" yaml:"show_color,omitempty" xml:"show_color,omitempty" mapstructure:"show_color,omitempty"`             // Habilita cores na saída
	ShowIcon    bool `json:"show_icon,omitempty" yaml:"show_icon,omitempty" xml:"show_icon,omitempty" mapstructure:"show_icon,omitempty"`                 // Habilita ícones na saída
	ShowTraceID bool `json:"show_trace_id,omitempty" yaml:"show_trace_id,omitempty" xml:"show_trace_id,omitempty" mapstructure:"show_trace_id,omitempty"` // Habilita o ID de rastreamento na saída
	ShowCaller  bool `json:"show_caller,omitempty" yaml:"show_caller,omitempty" xml:"show_caller,omitempty" mapstructure:"show_caller,omitempty"`         // Habilita informações do chamador na saída
	ShowStack   bool `json:"show_stack,omitempty" yaml:"show_stack,omitempty" xml:"show_stack,omitempty" mapstructure:"show_stack,omitempty"`             // Habilita informações da pilha de chamadas na saída
	ShowFields  bool `json:"show_fields,omitempty" yaml:"show_fields,omitempty" xml:"show_fields,omitempty" mapstructure:"show_fields,omitempty"`         // Habilita campos adicionais na saída

	// O formatter é armazenado somente como texto de referência ao tipo, e instanciado no momento da formatação.
	Format string `json:"format,omitempty" yaml:"format,omitempty" xml:"format,omitempty" mapstructure:"format,omitempty"` // json / text / xml / etc.

	Context string `json:"ctx,omitempty" yaml:"ctx,omitempty" xml:"ctx,omitempty" mapstructure:"ctx,omitempty"` // ex: "auth", "db", "billing"

	// Os campos Source, TraceID e Caller são campos com valores adquiridos por libs de sistema.
	// Mesmo se tentarem ser alterados via builder/chainable, o valor original será preservado e exibido que houve alteração na saída.
	// Garantindo que o valor verdadeiro seja sempre preservado.
	Source   string `json:"src,omitempty" yaml:"src,omitempty" xml:"src,omitempty" mapstructure:"src,omitempty"`             // componente/módulo/serviço
	TraceID  string `json:"trace,omitempty" yaml:"trace,omitempty" xml:"trace,omitempty" mapstructure:"trace,omitempty"`     // correlação
	Caller   string `json:"caller,omitempty" yaml:"caller,omitempty" xml:"caller,omitempty" mapstructure:"caller,omitempty"` // arquivo:linha função
	Severity int    `json:"sev,omitempty" yaml:"sev,omitempty" xml:"sev,omitempty" mapstructure:"sev,omitempty"`             // cache do Level.Severity()

	Tags   map[string]string `json:"tags,omitempty" yaml:"tags,omitempty" xml:"-" mapstructure:"tags,omitempty"`       // metadados arbitrários
	Fields map[string]any    `json:"fields,omitempty" yaml:"fields,omitempty" xml:"-" mapstructure:"fields,omitempty"` // dados estruturados arbitrários

	err error `json:"-" yaml:"-" xml:"-" mapstructure:"-"` // erro associado (se houver)
}

func NewKbxEntry(level kbx.Level) (kbx.LogzEntry, error) {
	return NewEntryImpl(string(level))
}

func NewEntry(level string) (*Entry, error) {
	l := kbx.ParseLevel(level)
	return &Entry{
		ctx:         context.Background(),
		ShowColor:   true,
		ShowIcon:    true,
		ShowTraceID: false,
		Timestamp:   time.Now().UTC(),
		Tags:        make(map[string]string),
		Fields:      make(map[string]any),
		Caller:      captureCaller(3),
		Level:       l,
		Severity:    l.Severity(),
	}, nil
}

// NewEntryImpl cria uma entry com:
// - timestamp UTC
// - maps inicializados
// - caller capturado
// - context
func NewEntryImpl(level string) (*Entry, error) {
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

func (e *Entry) WithLevel(level string) kbx.LogzEntry {
	e.Level = kbx.Level(level)
	e.Severity = e.Level.Severity()
	return e
}

// TODO: Fazer de forma que seja possível atribuir um context.Context.
// Vai ficar pra depois...
func (e *Entry) WithTraceID(id string) kbx.LogzEntry {
	e.TraceID = id
	return e
}

// TODO: Fazer de forma que seja possível atribuir um contexto.
func (e *Entry) WithContext(ctx context.Context) kbx.LogzEntry {
	e.ctx = ctx
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
	// Aqui não é feita nenhum processamento, somente a atribuição e retorno do ponteiro.
	// o uso de variadic arguments na criação dos métodos do logger já fazem a concatenação
	// e formatação dos valores em uma única string, o que facilita a vida do usuário.
	// TODO: A montagem da mensagem ocorre no formatter, porém aqui a propriedade deveria empilhar para
	// permitir que fosse feita somente lá.
	// Fazer essa implementação NÃO vai exigir uma pequena mudança na interface do io.Writer. Mas sim nos formatters.
	// Hoje a interface do io.Writer é simplesmente um escritor (io.Writer), ela só lida com []byte.
	// A montagem ocorre no Formatter, no método de dispatch. Porque assim permitimos que durante o ciclo de vida
	// do Entry (com o log original), só seja alterada no próprio formatter, preservando a que está aqui.
	// Nós recebemos o valor do formatter e o direcionamos para o writer, que irá escrever no output (io.Writer), independete
	// da mensagem ou do destino (console, arquivo, socket, etc.). Ele só escreve. O formatter só formata. É o que permite
	// ter uma lógica mais limpa e separada.
	e.Message = msg
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
	e.err = err
	return e
}

func (e *Entry) Tag(k, v string) kbx.LogzEntry {
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}
	e.Tags[k] = v
	return any(e).(kbx.LogzEntry)
}

func (e *Entry) Field(k string, v any) kbx.LogzEntry {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[k] = v
	return e
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

func (e *Entry) GetLevel() string {
	if e == nil {
		return string(kbx.LevelSilent)
	}
	return string(e.Level)
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

func (e *Entry) Error() error {
	if e == nil {
		return nil
	}
	return e.err
}

//
// ---------- Debug-friendly String() ----------
//

func (e *Entry) String() string {
	// Montagem da mensagem sem formatação especial/LogzFormatter.
	// Isso serve para que o usuário possa fazer um "fmt.Println(entry)" e
	// obter uma saída legível sem precisar chamar um formatter específico.
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
// ---------- Métodos de Destruição/Descarte ----------
//

// Reset reseta os campos da Entry
//
// Importante: O Reset não zera o TraceID, pois ele é gerado no momento da criação
// do Logger ou pelo método WithTraceID, e pode ser reutilizado em diversos Entries.
// Isso é feito propositalmente para que o TraceID seja preservado entre as chamadas
// do Logger, permitindo o monitoramento contínuo do fluxo de requisições.
// Isso também evita que a chave TraceID apareça em cada entrada individualmente.
//

func (e *Entry) Reset() {
	if e == nil {
		return
	}

	e.Message = ""
	e.Timestamp = time.Time{}
	e.Level = kbx.LevelSilent
	e.Severity = 0
	e.Caller = ""
	e.Context = ""
	e.Source = ""
	e.Format = ""
	e.err = nil
	e.Tags = nil
	e.Fields = nil
	e.ShowColor = true
	e.ShowStack = false
	e.ShowCaller = false
	e.ShowFields = false
	e.ShowIcon = false
	e.ShowTraceID = false
	e.TraceID = ""
}

// Propagate propaga a os eventos de relativos à entrada.
// O entry não vê o logger, ele só acessa os fluxos e estruturas que
// são utilizadas através dos dados alocados na própria entry.
// Isso faz com que o entry tenha um acoplamento baixo com o logger.
// E possa ser usado de forma independente, além de permitir a definição
// de diferentes fluxos de saída (IOBridge) para cenários mais complexos,
// flexíveis e com o máximo de customizações possíveis considerando a lógica
// de um logger moderno e robusto.
func (e *Entry) Propagate() error {
	if e == nil {
		return errors.New("entry is nil")
	}
	// lgr := logz.GetLoggerZ("")

	// lgr.

	// return e.Logger.Output(e.Formatter.Format(e))
	return nil
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
