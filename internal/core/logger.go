package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"

	"log"
)

// Logger é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type Logger struct {
	mu      sync.RWMutex
	flushMu sync.Mutex
	hooksMu sync.Mutex
	opts    *LoggerOptionsImpl
	*log.Logger
}

// LoggerZ é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type LoggerZ[T kbx.Entry] struct {
	ID       uuid.UUID
	flushMuZ sync.Mutex
	hooksMuZ sync.Mutex
	muZ      sync.RWMutex
	optsZ    *LoggerOptionsImpl
	interfaces.Logger
}

func NewLogger(prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	opts.Prefix = prefix

	// Configura o stdlog do Go para usar o mesmo output e prefixo
	var out io.Writer
	if opts.Output == nil {
		opts.Output = io.Discard
	} else {
		out = opts.Output
	}
	if out == nil {
		out = io.Discard
	}
	opts.Output = out
	logr := log.New(
		opts.Output,
		opts.Prefix,
		0,
	)

	lgr := &Logger{
		flushMu: sync.Mutex{},
		hooksMu: sync.Mutex{},
		mu:      sync.RWMutex{},
		opts:    opts,
		Logger:  logr,
	}
	// Reafirma configurações do log padrão
	lgr.SetFlags(0) // desativa flags automáticas do log padrão
	if kbx.DefaultFalse(opts.OutputTTY) {
		// se for TTY, desativa escrita direta no output padrão
		lgr.SetOutput(io.Discard) // evita escrita direta no output padrão
	} else {
		// se não for TTY, mantém escrita direta no output padrão
		lgr.SetOutput(out)
	}
	lgr.SetPrefix(prefix)
	lgr.SetFormatter(kbx.GetValueOrDefaultSimple(
		formatter.ParseFormatter(opts.LogzFormatOptions.Format, true),
		formatter.NewMinimalFormatter(true)),
	)
	lgr.SetPrefix(lgr.opts.Prefix)
	lgr.SetMinLevel(lgr.opts.MinLevel)

	return lgr
}

// NewLoggerZ cria um logger genérico:
// - formatter: serializa T em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLoggerZ[T kbx.Entry](prefix string, opts *LoggerOptionsImpl, withDefaults bool) *LoggerZ[T] {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	return &LoggerZ[T]{
		ID: uuid.New(),

		muZ:      sync.RWMutex{},
		flushMuZ: sync.Mutex{},
		hooksMuZ: sync.Mutex{},

		optsZ:  opts,
		Logger: NewLogger(prefix, opts, false), // evita chamada recursiva
	}
}

// NewLoggerZI cria um logger genérico:
// - formatter: serializa Record em []byte
// - out: destino final (io.Writer global, arquivo, socket, etc)
// - min: nível mínimo
func NewLoggerZI(prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	opts.Prefix = prefix

	// Configura o stdlog do Go para usar o mesmo output e prefixo
	var out io.Writer
	if opts.Output == nil {
		opts.Output = io.Discard
	} else {
		out = opts.Output
	}
	if out == nil {
		out = io.Discard
	}
	opts.Output = out
	logr := log.New(
		opts.Output,
		opts.Prefix,
		0,
	)

	lgr := &Logger{
		flushMu: sync.Mutex{},
		hooksMu: sync.Mutex{},
		mu:      sync.RWMutex{},
		opts:    opts,
		Logger:  logr,
	}
	// Reafirma configurações do log padrão
	lgr.SetFlags(0) // desativa flags automáticas do log padrão
	if kbx.DefaultFalse(opts.OutputTTY) {
		// se for TTY, desativa escrita direta no output padrão
		lgr.SetOutput(io.Discard) // evita escrita direta no output padrão
	} else {
		// se não for TTY, mantém escrita direta no output padrão
		lgr.SetOutput(out)
	}
	lgr.SetPrefix(prefix)
	lgr.SetFormatter(lgr.opts.Formatter)
	lgr.SetPrefix(lgr.opts.Prefix)
	lgr.SetMinLevel(lgr.opts.MinLevel)

	return lgr
}

func (l *Logger) SetFormatter(f formatter.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.opts.LogzFormatOptions == nil {
		l.opts.LogzFormatOptions = &kbx.LogzFormatOptions{}
	}
	l.opts.LogzFormatOptions.Format = f.Name()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Output = w
}

func (l *Logger) SetMinLevel(min kbx.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.MinLevel = min
}

func (l *Logger) AddHook(h interfaces.Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Hooks = append(l.opts.Hooks, h)
}

func (l *Logger) Enabled(level kbx.Level) bool {
	l.mu.RLock()
	min := l.opts.MinLevel
	l.mu.RUnlock()
	return level.Severity() >= min.Severity()
}

func (l *Logger) GetMinLevel() kbx.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.opts.MinLevel
}

func (l *Logger) GetLevel() kbx.Level {
	return l.opts.MinLevel
}

// SetRotate is the setter for setRotate
func (l *Logger) SetRotate(rotate bool) {
	// implementação fictícia
}

// SetRotateMaxSize is the setter for setRotateMaxSize
func (l *Logger) SetRotateMaxSize(size int64) {
	// implementação fictícia
}

// SetRotateMaxBack is the setter for setRotateMaxBack
func (l *Logger) SetRotateMaxBack(back int64) {
	// implementação fictícia
}

// SetRotateMaxAge is the setter for setRotateMaxAge
func (l *Logger) SetRotateMaxAge(age int64) {
	// implementação fictícia
}

// SetCompress is the setter for setCompress
func (l *Logger) SetCompress(compress bool) {
	// implementação fictícia
}

// SetBufferSize is the setter for setBufferSize
func (l *Logger) SetBufferSize(size int) {
	// implementação fictícia
}

// SetFlushInterval is the setter for setFlushInterval
func (l *Logger) SetFlushInterval(interval time.Duration) {
	// implementação fictícia
}

// SetHooks is the setter for setHooks
func (l *Logger) SetHooks(hooks []interfaces.Hook) {
	// implementação fictícia
}

// SetLHooks is the setter for setLHooks
func (l *Logger) SetLHooks(hooks interfaces.LHook[any]) {
	// implementação fictícia
}

// SetMetadata is the setter for setMetadata
func (l *Logger) SetMetadata(metadata map[string]any) {
	// implementação fictícia
}

func (l *Logger) GetConfig() *LoggerOptionsImpl {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.opts
}

type logParts struct {
	entries   []kbx.Entry
	others    []any
	jobLevel  kbx.Level
	timestamp time.Time
}

func (l *Logger) logEntryError(entry *Entry) error {
	// Ao ser criado o objeto já armazena o timestamp, que inclusive não
	// pode ser alterado depois.
	// Então se o timestamp estiver zerado, significa que o objeto
	// foi criado de forma incorreta. e não será possível corrigir isso aqui.
	// portanto será logado como erro de implementação.E NÃO SEGUIRÁ O FLUXO COM O RESTO!
	entryInstanceErrorLog, err := NewEntryImpl(kbx.LevelError)
	if err != nil {
		return err
	}
	pc, file, _, ok := runtime.Caller(2)

	if ok {
		fn := runtime.FuncForPC(pc)

		entryInstanceErrorLog = entryInstanceErrorLog.
			WithField("caller_function", fn.Name()).
			WithField("caller_file", file).
			WithField("caller_ok", ok)
	}

	entryInstanceErrorLog = entryInstanceErrorLog.
		WithMessage("logz: entry created with zero timestamp; this is an implementation error").
		WithField("entry_type", "Entry").
		WithField("entry_value", entry).
		WithError(fmt.Errorf("entry has zero timestamp"))

	if entryInstanceErrorLog.GetTimestamp().IsZero() {
		return nil
	}
	// Formata a Entry para ser logada
	b, err := l.opts.Formatter.Format(entryInstanceErrorLog)
	if err != nil {
		return err
	}
	// Garante newline pra saída de console / arquivos de texto.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	_, err = l.Writer().Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) getFormatter() (formatter.Formatter, error) {
	l.mu.RLock()
	f := kbx.GetValueOrDefaultSimple(
		formatter.ParseFormatter(kbx.GetValueOrDefaultSimple(l.opts, &LoggerOptionsImpl{
			LogzAdvancedOptions: &LogzAdvancedOptions{
				Formatter: formatter.ParseFormatter("minimal", true),
			},
		}).Format, true),
		formatter.ParseFormatter(l.opts.Format, true),
	)
	out := l.opts.Output
	l.mu.RUnlock()
	if f == nil || out == nil {
		// logger não inicializado corretamente; falha silenciosa
		return nil, fmt.Errorf("logger not properly initialized: formatter or output is nil")
	}
	return f, nil
}

func (l *Logger) dispatchLogEntry(entry *Entry) error {
	if l == nil || entry == nil {
		return nil
	}
	if !kbx.IsObjSafe(l, false) {
		return nil
	}
	if !kbx.IsObjSafe(entry, false) {
		return nil
	}

	// Isso já foi inferido antes de entrar nesse método. Ele é
	// privado, portanto é para uso interno apenas.
	// Nós somente iremos reafirmar o que é passível de ser reafirmado.
	// Como o nível do log. O level DEVE SER PASSADO por argumento.
	// Porém, como o Entry também possui a informação, nós iremos
	// considerar o que está no entry SOMENTE QUANDO HOUVER VÁRIOS ENTRIES!!!
	// Isso porque, se houver vários entries, pode haver intençãoes
	// diferentes entre eles, podem compor um bloco de log enviado de uma vez.
	if !l.Enabled(entry.GetLevel()) {
		return nil
	}

	// obtém o formatter

	f, err := l.getFormatter()
	if err != nil {
		return err
	}

	// formata a Entry

	b, err := f.Format(entry)
	if err != nil {
		return err
	}

	// garante newline pra saída de console / arquivos de texto.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	// dispara hooks pré-formatação
	if l.GetConfig() != nil {
		if l.GetConfig().LogzAdvancedOptions != nil {
			if l.GetConfig().LogzAdvancedOptions.Hooks != nil {
				err := l.GetConfig().LHooks.Fire(entry)
				if err != nil {
					return err
				}
			}
		}
	}

	// escreve no destino final
	_, err = l.Writer().Write(b)
	if err != nil {
		return err
	}

	// tudo ok
	return nil
}

// Log é o caminho principal: recebe um Record pronto (T),
// dispara hooks, formata e escreve em out.
func (l *Logger) Log(lvl kbx.Level, rec ...any) error {
	if !kbx.IsObjSafe(rec, false) {
		// nada a fazer, mas não vamos quebrar ninguém
		return nil
	}

	var logParts = logParts{
		entries:   make([]kbx.Entry, 0),
		others:    make([]any, 0),
		jobLevel:  lvl,
		timestamp: time.Now(),
	}
	if len(rec) > 0 {
		for _, r := range rec {
			if !kbx.IsObjSafe(r, false) {
				continue
			}
			if e, ok := r.(kbx.Entry); ok {
				logParts.entries = append(logParts.entries, e)
			} else {
				logParts.others = append(logParts.others, r)
			}
		}
	}

	/////////////////////////////////////////////////////////////////////
	/// Agoa vamos separar o que é Entry do que não é. PRimeiro Entries
	/// Vamos pegar todas elas e simplesmente disparar um Log para cada
	/// visto que cada entry possui todo o contexto que compõe o log.
	/////////////////////////////////////////////////////////////////////
	for pos, entry := range logParts.entries {
		// garante que o nível do job seja respeitado
		if entry.GetLevel() == "" || kbx.Level(entry.GetLevel()).Severity() < logParts.jobLevel.Severity() {
			entry = entry.(*Entry).WithLevel(logParts.jobLevel)
		}
		// garante timestamp
		if err := entry.Validate(); err != nil {
			if err := l.logEntryError(entry.(*Entry)); err != nil {
				return err
			}
			logParts.entries[pos] = nil
			continue
		} else {
			if err := l.dispatchLogEntry(entry.(*Entry)); err != nil {
				return err
			}
		}
	}

	/////////////////////////////////////////////////////////////////////
	/// Agora, TODOS OS OUTROS objetos que estavam na lista de argumentos
	/////////////////////////////////////////////////////////////////////
	if len(logParts.others) > 0 {
		entryInstance, err := NewEntryImpl(lvl)
		if err != nil {
			return err
		}
		entry := entryInstance.
			WithLevel(lvl)

		var msgParts = make([]string, 0)
		for _, other := range logParts.others {
			if str, ok := other.(string); ok {
				if str != "" {
					msgParts = append(msgParts, str)
				}
			} else if errObj, ok := other.(error); ok {
				entry = entry.WithError(errObj)
			} else if m, ok := other.(map[string]any); ok {
				for k, v := range m {
					entry = entry.WithField(k, v)
				}
			} else {
				// tenta serializar como json
				jsonBytes, err := json.MarshalIndent(other, "", "  ")
				if err == nil && len(jsonBytes) > 0 {
					msgParts = append(msgParts, string(jsonBytes))
				} else {
					// fallback simples
					msgParts = append(msgParts, fmt.Sprintf("%v", other))
				}
			}
		}
		entry = entry.WithMessage(fmt.Sprintf("%s", msgParts))
		// dispara o log
		if err := l.dispatchLogEntry(entry); err != nil {
			return err
		}
	}
	/////////////////////////////////////////////////////////////////////
	/// Fim do pipeline de log para objetos não estruturados.
	///////////////////////////////////////////////////////////////////
	return nil
}

func (l *Logger) LogAny(level kbx.Level, args ...any) error {
	if l == nil {
		return nil
	}
	if len(args) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			if l.Logger != nil {
				l.Logger.Printf("logz: panic in LogAny: %v (args=%#v)", r, args)
			}
		}
	}()

	// modo moderno: nada garante level → assume Info
	entry := toEntry(level, args...)
	if entry.GetLevel() == "" {
		entry = entry.WithLevel(level)
	}

	return l.Log(level, entry.GetLevel().String(), entry)
}

func (l *LoggerZ[T]) Clone() *LoggerZ[T] {
	l.muZ.RLock()
	defer l.muZ.RUnlock()
	newOpts := l.optsZ.Clone()
	return NewLoggerZ[T](l.optsZ.Prefix, newOpts, false)
}

// SetDebugMode habilita ou desabilita o modo debug do logger global.
// Quando debug=true, mostra logs de todos os níveis (incluindo debug e trace).
// Quando debug=false, mostra apenas logs de nível info ou superior.
func (l *LoggerZ[T]) SetDebugMode(debug bool) {
	if l == nil {
		return
	}
	if debug {
		l.SetMinLevel(kbx.LevelDebug)
	} else {
		l.SetMinLevel(kbx.LevelInfo)
	}
}

// Debug loga uma mensagem de debug
func (l *LoggerZ[T]) Debug(msg ...any) {
	l.Log("debug", msg...)
}

// Notice loga uma mensagem de notice
func (l *LoggerZ[T]) Notice(msg ...any) {
	l.Log("notice", msg...)
}

// Info loga uma mensagem informativa
func (l *LoggerZ[T]) Info(msg ...any) {
	l.Log("info", msg...)
}

// Success loga uma mensagem de sucesso
func (l *LoggerZ[T]) Success(msg ...any) {
	l.Log("success", msg...)
}

// Warn loga um aviso
func (l *LoggerZ[T]) Warn(msg ...any) {
	l.Log("warn", msg...)
}

// Error loga um erro e retorna error
func (l *LoggerZ[T]) Error(msg ...any) error {
	return l.Log("error", msg...)
}

// Fatal loga uma mensagem fatal e encerra o programa com exit code 1
func (l *LoggerZ[T]) Fatal(msg ...any) {
	l.Log("fatal", msg...)
	os.Exit(1)
}

func (l *LoggerZ[T]) Trace(msg ...any) {
	l.Log("trace", msg...)
}

func (l *LoggerZ[T]) Critical(msg ...any) {
	l.Log("critical", msg...)
}

func (l *LoggerZ[T]) Answer(msg ...any) {
	l.Log("answer", msg...)
}

func (l *LoggerZ[T]) Alert(msg ...any) {
	l.Log("alert", msg...)
}

func (l *LoggerZ[T]) Bug(msg ...any) {
	l.Log("bug", msg...)
}

func (l *LoggerZ[T]) Panic(msg ...any) {
	l.Log("panic", msg...)
}
func (l *LoggerZ[T]) Println(msg ...any) {
	l.Log("println", msg...)
}
func (l *LoggerZ[T]) Printf(format string, args ...any) {
	l.Log("printf", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Debugf(format string, args ...any) {
	l.Log("debug", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Infof(format string, args ...any) {
	l.Log("info", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Noticef(format string, args ...any) {
	l.Log("notice", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Successf(format string, args ...any) {
	l.Log("success", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Warnf(format string, args ...any) {
	l.Log("warn", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Errorf(format string, args ...any) error {
	return l.Log("error", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Fatalf(format string, args ...any) {
	l.Log("fatal", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Tracef(format string, args ...any) {
	l.Log("trace", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Criticalf(format string, args ...any) {
	l.Log("critical", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Answerf(format string, args ...any) {
	l.Log("answer", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Alertf(format string, args ...any) {
	l.Log("alert", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Bugf(format string, args ...any) {
	l.Log("bug", fmt.Sprintf(format, args...))
}
func (l *LoggerZ[T]) Panicf(format string, args ...any) {
	l.Log("panic", fmt.Sprintf(format, args...))
}
