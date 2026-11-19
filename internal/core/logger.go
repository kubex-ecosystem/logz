package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/logz/internal/events"
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
	// Aqui a struct do stdlog está agregada ao logger, promovendo (herdando) todos os métodos do stdlog
	// Ela é um ponteiro, logo deve ser instanciada e atribuída sempre ao criar um Logger
	*log.Logger
	ID      uuid.UUID
	mu      sync.RWMutex
	flushMu sync.Mutex
	hooksMu sync.Mutex
	opts    *LoggerOptionsImpl
}

// LoggerZ é o núcleo do pipeline:
//
//	Record (T) -> hooks -> formatter -> io.Writer
//
// Não sabe nada de linha, arquivo, CLI, JSON, etc.
// Isso é responsabilidade do Formatter + destino (io.Writer).
type LoggerZ[T kbx.Entry] struct {
	ID      uuid.UUID
	muZ     sync.RWMutex
	optsZ   *LoggerOptionsImpl
	*Logger // Essa struct promove a struct global do stdlog
}

// NewLogger is the same as NewLoggerWithLogger but with nil logger and false withDefaults
// Compatibility constructor
// @deprecated
func NewLogger(prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	return NewLoggerWithLogger(nil, prefix, opts, withDefaults)
}

// NewLoggerWithLogger creates a new logger with the specified prefix.
// If l is not nil, it will use the same logger as l.
// If opts is not nil, it will use the same options as opts.
// If withDefaults is true, it will use the default options.
func NewLoggerWithLogger(l *log.Logger, prefix string, opts *LoggerOptionsImpl, withDefaults bool) *Logger {
	if opts == nil {
		opts = NewLoggerOptions(kbx.LoggerArgs)
	}
	if withDefaults {
		opts = opts.WithDefaults(opts)
	}
	opts.Prefix = prefix
	var out io.Writer
	if opts.Output == nil {
		opts.Output = io.Discard
	} else {
		out = opts.Output
	}
	if out == nil {
		out = io.Discard // Se for nulo, inicializa com o Discard
	}
	opts.Output = out
	// Configura o stdlog do Go para usar o mesmo output e prefixo
	if l == nil {
		// Aqui nós terminamos de checar todas as configurações e preenchemos com padrões as que são requeridas
		// Agora podemos criar o logger padrão
		l = log.New(opts.Output, opts.Prefix, 0)
	} else {
		// Configura o stdlog do Go para usar o mesmo output e prefixo
		l.SetOutput(out)
		l.SetPrefix(prefix)
		l.SetFlags(0)
	}
	// Cria o Logger do logz, com o stdlog (log.Logger) em composition
	// Nós estamos compondo com um ponteiro do log.Logger global
	// Isso permite que o Logger tenha em sua struct todos os métodos e propriedades
	// do stdlog. Com isso, além de conseguir extrair a instância original
	// do stdlog com o método .Logger(), nós conseguimos usar todos os métodos
	// do stdlog diretamente no Logger. Por isso, quando criarmos uma nova instância
	// Se não passamos a instância do logger global, nós criamos uma nova baseada nas configurações
	// do Logger (logz.LoggerOptions), com isso, o logz segue com as capacidades do stdlog.
	// Porém o logz pode ser usado standalone, sem depender do stdlog, caso seja setado um writer
	// e outras configurações estruturais, como formatador e etc.. ele não irá utilizar o stdlog.
	// Nós preservamos a instância do stdlog para permitir uma propagação de mensagens e outras coisas
	// que o stdlog não promove de formá passível de ser manipulada de qualquer lib externa.
	// Quando o usuário não fornece uma instância do stdlog, criamos uma nova baseada nas configurações
	// do Logger (logz.LoggerOptions) somente para propagar determinados comportamentos e outras propriedades
	// que demandam a instância dele. Porém, nesse caso nós usamos os nossos próprios mecanismos para logging,
	// o que garante que o logz funcione de forma independente, sem precisar do stdlog.
	// Por isso, quando usamos essa opção de criar uma nova instância do stdlog,
	// A nossa estrutura é dotada de métodos chainable, que nos permite configurar o logger de forma encadeada
	// independente da instância do stdlog. Por isso, quando criamos uma nova instância
	// do stdlog, nós não inicializamos as outras configurações, já que serão "clonadas/herdadas" através
	// de qualquer método chainable acionado
	lgr := &Logger{
		flushMu: sync.Mutex{},
		hooksMu: sync.Mutex{},
		mu:      sync.RWMutex{},
		opts:    opts,
		// Aqui é inserida a instância original do stdlog. A primeira e única que o logz usa para ocupar essa instância.
		// Ela possui esse nome Logger, porque da forma que o Go é estruturado, o composition que promove a estrutura
		// inteira é somente com um ponteiro na raiz da struct, para evitar que se crie outra instância.
		Logger: l,
	}

	// A cada chamada desse constructor será criada uma nova instância do logger global...
	// Mas isso só acontece se o dev/usuario NÃO fornecer uma instância do stdlog nas suas
	// chamadas desse constructor, ou se ele simplesmente não criar infinitos loggers.
	// De qualquer forma, o logz possui uma vida curta e será destruído em pouco tempo, a não ser
	// que a instância seja globalizada ao ser declarada na raiz de algum módulo/package do projeto que
	// esteja utilizando o logz. Nesse ultimo caso, será criada apenas uma instância do logger global,
	// e essa única instância será compartilhada por todos os módulos/packages do projeto.

	// Reafirma configurações do log padrão
	//
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
		formatter.ParseFormatter(opts.Format, true),
		formatter.NewMinimalFormatter(true)),
	)
	lgr.SetPrefix(lgr.opts.Prefix)
	lgr.SetMinLevel(lgr.opts.MinLevel)
	lgr.SetConfig(lgr.opts.LoggerConfig)
	metaData := make(map[string]any)
	for k, v := range lgr.opts.LoggerConfig.Metadata {
		metaData[k] = v
	}
	lgr.SetMetadata(metaData)

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

		muZ: sync.RWMutex{},

		optsZ:  opts,
		Logger: NewLoggerWithLogger(nil, prefix, opts, false), // evita chamada recursiva
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
	lgr.SetConfig(opts.LoggerConfig)
	metaData := make(map[string]any)
	for k, v := range lgr.opts.LoggerConfig.Metadata {
		metaData[k] = v
	}
	lgr.SetMetadata(metaData)

	return lgr
}

// SetFormatter sets the formatter for the logger.
func (l *Logger) SetFormatter(f formatter.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.opts.LogzFormatOptions == nil {
		if kbx.LoggerArgs == nil {
			kbx.ParseLoggerArgs(l.opts.Level.String(), l.opts.MinLevel.String(), l.opts.MaxLevel.String(), kbx.GetValueOrDefaultSimple(*l.opts.OutputFile, "stdout"))
		}
		l.opts.LogzFormatOptions = kbx.LoggerArgs.LogzFormatOptions
	}
	l.opts.Format = f.Name()
}

// SetOutput sets the output writer for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Output = w
}

// SetMinLevel sets the minimum level for the logger.
func (l *Logger) SetMinLevel(min kbx.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.MinLevel = min
}

// AddHook adds a hook to the logger.
func (l *Logger) AddHook(h events.Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Hooks = append(l.opts.Hooks, h)
}

// Enabled checks if the logger is enabled for the given level.
func (l *Logger) Enabled(level kbx.Level) bool {
	l.mu.RLock()
	min := l.opts.MinLevel
	l.mu.RUnlock()
	return level.Severity() >= min.Severity()
}

// GetMinLevel returns the minimum level for the logger.
func (l *Logger) GetMinLevel() kbx.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.opts.MinLevel
}

// GetLevel returns the level for the logger.
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
func (l *Logger) SetHooks(hooks []events.Hook) {
	// implementação fictícia
}

// SetLHooks is the setter for setLHooks
func (l *Logger) SetLHooks(hooks events.LHook[any]) {
	// implementação fictícia
}

// SetMetadata is the setter for setMetadata
func (l *Logger) SetMetadata(metadata map[string]any) {
	// implementação fictícia
}

// GetConfig returns the config for the logger.
func (l *Logger) GetConfig() *LoggerOptionsImpl {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.opts
}

// SetConfig sets the config for the logger.
func (l *Logger) SetConfig(opts *kbx.LogzConfig) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if opts == nil {
		return
	}
	if opts.LogzFormatOptions != nil {
		l.opts.LogzFormatOptions = opts.LogzFormatOptions
	}
	if opts.Output != nil {
		l.opts.Output = opts.Output
	}

	l.opts.Level = opts.Level
	l.opts.MinLevel = opts.MinLevel
	l.opts.MaxLevel = opts.MaxLevel

	if opts.Prefix != "" {
		l.opts.Prefix = opts.Prefix
	}

	l.opts.ShowFields = opts.ShowFields
	l.opts.ShowIcons = opts.ShowIcons
	l.opts.ShowColor = opts.ShowColor
	l.opts.ShowStack = opts.ShowStack
	l.opts.ShowTraceID = opts.ShowTraceID
	l.opts.Output = opts.Output
	l.opts.LoggerConfig.Metadata = opts.Metadata
	l.opts.StackTrace = opts.StackTrace
}

type logParts struct {
	entries   []kbx.LogzEntry
	others    []any
	jobLevel  kbx.Level
	timestamp time.Time
}

func (l *Logger) logEntryError(entry kbx.LogzEntry) error {
	// Ao ser criado o objeto já armazena o timestamp, que inclusive não
	// pode ser alterado depois.
	// Então se o timestamp estiver zerado, significa que o objeto
	// foi criado de forma incorreta. e não será possível corrigir isso aqui.
	// portanto será logado como erro de implementação.E NÃO SEGUIRÁ O FLUXO COM O RESTO!
	entryInstanceErrorLog := NewLogzEntry(kbx.LevelError)
	pc, file, _, ok := runtime.Caller(2)
	if ok {
		fn := runtime.FuncForPC(pc)

		entryInstanceErrorLog = entryInstanceErrorLog.
			WithField("caller_function", fn.Name()).
			WithField("caller_file", file).
			WithField("caller_ok", ok)
	}

	if entry.GetTimestamp().IsZero() {
		entryInstanceErrorLog = entryInstanceErrorLog.
			WithMessage("logz: entry created with zero timestamp; this is an implementation error").
			WithField("entry_type", "Entry").
			WithField("entry_value", entry).
			WithError(fmt.Errorf("entry has zero timestamp"))
	}

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

	// Escreve no destino final
	if l.Enabled(l.GetLevel()) {
		_, err = l.Writer().Write(b)
		if err != nil {
			return err
		}
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

func (l *Logger) dispatchLogEntry(entry kbx.LogzEntry) error {
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
	if !l.Enabled(kbx.Level(entry.GetLevel())) {
		return nil
	}

	// Obtém o formatter escolhido e alocado na Entry

	f, err := l.getFormatter()
	if err != nil {
		return err
	}

	// Formata a Entry utilizando o Formatter correto e transformando-a em []byte
	// O valor formatado é armazenado na variável b. Se b não for nil, significa que o log foi formatado corretamente.
	// Se b for nil, significa que o log não foi formatado corretamente e será escrito no destino final, porém em formato de erro.
	// O método Format tem essa responsabilidade de, caso ocorra um erro no processo, o erro é logado no formato de erro
	// e é retornando a Entry formatada com level de erro.
	b, err := f.Format(entry)
	if err != nil {
		return err
	}

	// Garante newline pra saída de console / arquivos de texto.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	// Dispara hooks pré-formatação (se houver algum anexado à Entry, ou global no próprio Logger)
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

	// Se a entry estiver desabilitada, o logger não escreve ela no output. Porém
	// todos os outros métodos como o disparo do hook, etc, vão ocorrer normalmente permitindo
	// que haja observabilidade, monitoramento, etc... Conforme o usuário definir.
	if l.Enabled(kbx.Level(entry.GetLevel())) {
		// escreve no destino final
		_, err = l.Writer().Write(b)
		if err != nil {
			return err
		}
	}

	// Após escrevermos a mensagem no destino final,
	// verificamos se o nível é Fatal, Panic, ou Critical.
	// Se for, o logger deve ser encerrado imediatamente e é enviado um SIGNAL para
	// informar o sistema operacional que a aplicação deve ser encerrada.
	if entry.GetLevel() == kbx.LevelFatal.String() ||
		entry.GetLevel() == kbx.LevelPanic.String() ||
		entry.GetLevel() == kbx.LevelCritical.String() {
		os.Exit(1)

		entry.Reset()
	}

	// Finalizamos o fluxo das Entries e Logging.
	return nil
}

// Log is the main entrypoint to log:
//   - dispatches hooks
//   - formats
//   - writes to out
func (l *Logger) Log(lvl kbx.Level, rec ...any) error {
	if !kbx.IsObjSafe(rec, false) {
		// nada a fazer, mas não vamos quebrar ninguém
		return nil
	}

	if lvl.Severity() < l.GetMinLevel().Severity() {
		// Quem determina se o log será impresso é a lógica
		// do dispatcher.
		// O minimum level é somente para evitar que
		// se grave níveis abaixo do mínimo no output, mesmo caso
		// o usuário não tenha passado nenhuma configuração (resiliente)
		// O método Enable faz a validação imperativa se o logger
		// estará habilitado ou não.. O Entry também possui
		// essa propriedade, portanto a configuração pode ser
		// granular e definida em diferentes níveis.
		// O entry só é definido no parser ou caso o usuário tenha
		// instanciado ele e utilizado ele diretamente em qualquer método de logging.
		// Do Nothing... Não fazemos literalmente nada aqui, esse trecho só está
		// aqui para informação/documentação. Deixamos o fluxo seguir para
		// que todas as capacidades e funcionalidades do logz possam ser utilizadas
		// em toda sua plenitude. (hooks, file output, notifiers, observability, metrics, etc..)
	}

	var logParts = logParts{
		entries:   make([]kbx.LogzEntry, 0),
		others:    make([]any, 0),
		jobLevel:  lvl,
		timestamp: time.Now(),
	}

	/////////////////////////////////////////////////////////////////////
	/// Aqui fazemos a primeira triagem nos argumentos recebidos.
	/// nós pegamos do array de parâmetros (rec) o que é Entry
	/// e o que não é e separamos em coleções distintas.
	///
	/// Apenas Entries serão tratadas no bloco
	/// O que for Entry é passível de seguir o fluxo principal e final
	/// da vida útil de um Entry enquanto veículo e portador da informação enquanto log.
	/// ```go
	///  //...
	///   if e, ok := r.(kbx.Entry); ok {
	///     logParts.entries = append(logParts.entries, e)
	///   }
	///   //...
	/// ```
	///
	/// As outras coisas serão tratadas no bloco
	/// Aqui o sistema irá montar uma coleção de objetos (any) para criar um novo log.
	/// O único modo de se criar um Log Entry é utilizando NewLogzEntry() ou NewEntry() que
	/// já inicializa informações imutáveis como o timestamp, etc... Todo restante é passível
	/// de atribuição através dos métodos chainable (entry.WithXXX(val XXX) Entry).
	/// O Entry é imutável, portanto cada chamada a WithXXX retorna um clone do Entry com as propriedades
	/// imutáveis e as que já foram atribuídas, porém a alocação na memória é somente para
	/// essa chamada, que quando é finalizada o GC irá liberar, permanecendo somente o clone.
	/// Não confunda o ENTRY com o Objeto Logger, que são completamente diferentes. O Entry é apenas um veículo
	/// portador da informação. Ele é descartado a cada ciclo do método Log.
	/// "if len(logParts.others) > 0" (abaixo)
	/// ```go
	///   //...
	///   else {
	///     logParts.others = append(logParts.others, r)
	///   }
	///   //...
	/// ```
	/////////////////////////////////////////////////////////////////////
	if len(rec) > 0 {
		for _, r := range rec {
			if !kbx.IsObjSafe(r, false) {
				continue
			}
			if e, ok := r.(*Entry); ok {
				logParts.entries = append(logParts.entries, e)
			} else {
				logParts.others = append(logParts.others, r)
			}
		}
	}

	lz, ok := any(l).(*LoggerZ[kbx.Entry])
	if !ok || lz == nil {
		lz = NewLoggerZ[kbx.Entry](l.Prefix(), l.GetConfig(), false)
		if lz == nil {
			// return l.Errorf("erro ao criar loggerz")
			ee := NewLogzEntry(kbx.LevelError).
				WithMessage("erro ao criar loggerz").
				WithError(errors.New("erro ao criar loggerz")).
				WithContext(context.Background()). // l.ctx é o contexto do Logger
				WithFields(map[string]any{
					"logger": "loggerz",
					"error":  errors.New("erro ao criar loggerz"),
				}).
				WithShowFields(l.GetConfig().ShowFields).
				WithShowCaller(l.GetConfig().ShowStack).
				WithShowTraceID(l.GetConfig().ShowTraceID).
				WithColor(*l.GetConfig().ShowColor).
				WithIcon(*l.GetConfig().ShowIcons).
				WithFormat(l.GetConfig().Format)

			_ = l.Log(kbx.LevelError, ee)
		}
	}
	/////////////////////////////////////////////////////////////////////
	/// 1º) Aqui percorremos todos os Entries para disparar os logs.
	/// Como cada Entry já é um veículo completo de informação,
	/// basta fazer o despacho de cada um deles individualmente.
	/// ```go
	///   //...
	///   if e, ok := r.(kbx.Entry); ok {
	///     logParts.entries = append(logParts.entries, e)
	///   }
	///   //...
	/// ```
	/////////////////////////////////////////////////////////////////////
	for pos, entry := range logParts.entries {
		// garante que o nível do job seja respeitado
		if l.Enabled(kbx.Level(entry.GetLevel())) {
			entry = entry.WithLevel(string(logParts.jobLevel))
		} else {
			continue
		}
		// garante timestamp
		if err := entry.Validate(); err != nil {
			// aqui não retornamos o erro, pois não queremos quebrar o fluxo do método log.
			// mas registramos o erro que ocorreu.
			lz.logEntryError(entry)
			logParts.entries[pos] = nil

			// O erro será registrado aqui, com level Debug.
			// lz.logEntryError(entry.(kbx.LogzEntry))
			continue
		} else {
			// Caso não haja erros de validação, disparamos os hooks e encaminhamos para o output final.
			if err := l.dispatchLogEntry(entry); err != nil {
				// aqui não retornamos o erro, pois não queremos quebrar o fluxo do método log.
				// mas registramos o erro que ocorreu.
				// l.logEntryError(entry.(*Entry)) // Esse método dispara um Log de erro com o conteúdo do entry que falhou.
				// removemos o entry da lista de entries
				logParts.entries[pos] = nil
				continue
			}
		}
	}

	/////////////////////////////////////////////////////////////////////
	/// 2º) Agora, TODOS OS OUTROS objetos que estavam na lista de argumentos
	/// Aqui o sistema irá montar uma coleção de objetos (any) para criar um novo log.
	/// O único modo de se criar um Log Entry é utilizando NewLogzEntry() ou NewEntry() que
	/// já inicializa informações imutáveis como o timestamp, etc... Todo restante é passível
	/// de atribuição através dos métodos chainable (entry.WithXXX(val XXX) Entry).
	/// O Entry é imutável, portanto cada chamada a WithXXX retorna um clone do Entry com as propriedades
	/// imutáveis e as que já foram atribuídas, porém a alocação na memória é somente para
	/// essa chamada, que quando é finalizada o GC irá liberar, permanecendo somente o clone.
	/// Não confunda o ENTRY com o Objeto Logger, que são completamente diferentes. O Entry é apenas um veículo
	/// portador da informação. Ele é descartado a cada ciclo do método Log.
	/// - utilizamos o level do logger e não da Entry em si, somente se tiver sido passado um level
	/// na chamada do método Log. Se não tivermos a informação do level
	/// - para preencher a mensagem, nós validamos o tipo de cada argumento, para que ele
	/// seja impresso da melhor forma possível, sempre mantendo o order de recebimento dos argumentos.
	/// - para preencher os campos (fields), nós validamos o tipo de cada argumento, para que ele
	/// seja impresso da melhor forma possível, sempre mantendo o order de recebimento dos argumentos.
	/// - para preencher os erros, nós validamos o tipo de cada argumento, para que ele
	/// seja impresso da melhor forma possível, sempre mantendo o order de recebimento dos argumentos.
	///
	/// ```go
	///   var entry LogzEntry = NewEntry(lvl)
	///   //...
	///   entry = entry.WithField("key", "value")
	///   //...
	///   entry = entry.WithError(errors.New("error"))
	///   //...
	///   entry = entry.WithMessage("message")
	/// ```
	///
	/// Para essas situações inesperadas e relativamente comuns, priorizamos a
	/// performance evitando parseamentos complexos e onerosos.
	/// Tentamos identificar o tipo de objeto de forma rápida e direta, sem
	/// criar overhead desnecessário. Portanto, se você precisar de um log estruturado
	/// ou personalizado, utilize o método LogEntry(), ou envie textos com level
	/// corretamente, pois como já falamos, a ordem dos parâmetros é importante.
	/// >NOTA: Os métodos de conveniência são wrappers para o método Log()
	/// e portanto seguem as mesmas regras. Se você utiliza strings como mensagem e
	/// insere o level corretamente, entrada por entrada, o logz irá montar a mensagem
	/// automaticamente, para o log level correspondente com um Entry novo para cada
	/// entrada, e isso tornará o log totalmente nativo com um Entry limpo.
	/// Esse é o método mais simples e de menor overhead para criação de logs.
	///
	/// ```go
	/// // errado
	/// logger.Error(kbx.Error, "Mensagem", errors.New("erro"))
	///
	/// // errado
	/// logger.Log(kbx.Error, errors.New("erro"), "Mensagem")
	///
	/// // certo
	/// logger.Log(kbx.Error, "Mensagem", "Erro", errors.New("erro"))
	///
	/// // ou use o método LogEntry() que é mais flexível.
	/// entry := NewEntry(kbx.Error)
	/// entry = entry.WithMessage("Mensagem 1")
	/// entry = entry.WithMessage("Mensagem 2")
	/// entry = entry.WithMessage("Mensagem 3")
	/// entry = entry.WithMessage("Mensagem 4")
	/// entry = entry.WithError(errors.New("erro"))
	/// logger.LogEntry(entry)
	/// ```
	///
	/////////////////////////////////////////////////////////////////////
	if len(logParts.others) > 0 {
		entry := NewLogzEntry(lvl).(*Entry)

		var msgParts = make([]string, 0)
		for _, other := range logParts.others {
			if str, ok := other.(string); ok {
				if str != "" {
					msgParts = append(msgParts, str)
				}
			} else if errObj, ok := other.(error); ok {
				entry = entry.WithError(errObj).(*Entry)
			} else if m, ok := other.(map[string]any); ok {
				for k, v := range m {
					entry = entry.WithField(k, v).(*Entry)
				}
			} else {
				// tenta serializar como json
				jsonBytes, err := json.MarshalIndent(other, "", "  ")
				if err == nil && len(jsonBytes) > 0 {
					msgParts = append(msgParts, string(jsonBytes))
				} else {
					// fallback simples (formatando o objeto de forma literal)
					msgParts = append(msgParts, fmt.Sprintf("%v", other))
				}
			}
		}
		entry = entry.WithMessage(fmt.Sprintf("%s", msgParts)).(*Entry)
		if err := l.dispatchLogEntry(entry); err != nil {
			return err
		}
	}
	/////////////////////////////////////////////////////////////////////
	/// Fim do pipeline de log para objetos não estruturados.
	///////////////////////////////////////////////////////////////////
	return nil
}

// LogAny is a wrapper around Log for logging arbitrary arguments.
// If len(args) == 0, it will return nil.
// Note: this method is not thread-safe. Use LoggerZ for a thread-safe logger.
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
				l.Printf("logz: panic in LogAny: %v (args=%#v)", r, args)
			}
		}
	}()
	return l.Log(level, args...)
}

// Clone returns a copy of the logger with the same configuration.
// This is useful for creating a logger with a different prefix.
func (l *LoggerZ[T]) Clone() *LoggerZ[T] {
	l.muZ.RLock()
	defer l.muZ.RUnlock()
	newOpts := l.optsZ.Clone()
	return NewLoggerZ[T](l.optsZ.Prefix, newOpts, false)
}

// SetDebugMode sets the debug mode of the logger.
// When debug=true, it shows logs of all levels (including debug and trace).
// When debug=false, it shows only logs of level info or higher.
func (l *LoggerZ[T]) SetDebugMode(debug bool) {
	if l == nil {
		return
	}
	if debug {
		l.SetMinLevel(kbx.LevelDebug)
	} else {
		l.SetMinLevel(l.GetMinLevel())
	}
}

// Debug logs a debug message.
func (l *LoggerZ[T]) Debug(msg ...any) {
	l.Log("debug", msg...)
}

// Notice logs a notice message.
func (l *LoggerZ[T]) Notice(msg ...any) {
	l.Log("notice", msg...)
}

// Info logs an info message.
func (l *LoggerZ[T]) Info(msg ...any) {
	l.Log("info", msg...)
}

// Success logs a success message.
func (l *LoggerZ[T]) Success(msg ...any) {
	l.Log("success", msg...)
}

// Warn logs a warn message.
func (l *LoggerZ[T]) Warn(msg ...any) {
	l.Log("warn", msg...)
}

// Error logs an error message and returns an error.
func (l *LoggerZ[T]) Error(msg ...any) error {
	return l.Log("error", msg...)
}

// Fatal logs a fatal message and exits the program with exit code 1.
func (l *LoggerZ[T]) Fatal(msg ...any) {
	l.Log("fatal", msg...)
	os.Exit(1)
}

// Trace logs a trace message.
func (l *LoggerZ[T]) Trace(msg ...any) {
	l.Log("trace", msg...)
}

// Critical logs a critical message.
func (l *LoggerZ[T]) Critical(msg ...any) {
	l.Log("critical", msg...)
}

// Answer logs an answer message.
func (l *LoggerZ[T]) Answer(msg ...any) {
	l.Log("answer", msg...)
}

// Alert logs an alert message.
func (l *LoggerZ[T]) Alert(msg ...any) {
	l.Log("alert", msg...)
}

// Bug logs a bug message.
func (l *LoggerZ[T]) Bug(msg ...any) {
	l.Log("bug", msg...)
}

// Panic logs a panic message.
func (l *LoggerZ[T]) Panic(msg ...any) {
	l.Log("panic", msg...)
}

// Println logs a println message.
func (l *LoggerZ[T]) Println(msg ...any) {
	l.Log("println", msg...)
}

// Printf logs a formatted message.
func (l *LoggerZ[T]) Printf(format string, args ...any) {
	l.Log("printf", fmt.Sprintf(format, args...))
}

// Debugf logs a formatted debug message.
func (l *LoggerZ[T]) Debugf(format string, args ...any) {
	l.Log("debug", fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message.
func (l *LoggerZ[T]) Infof(format string, args ...any) {
	l.Log("info", fmt.Sprintf(format, args...))
}

// Noticef logs a formatted notice message.
func (l *LoggerZ[T]) Noticef(format string, args ...any) {
	l.Log("notice", fmt.Sprintf(format, args...))
}

// Successf logs a formatted success message.
func (l *LoggerZ[T]) Successf(format string, args ...any) {
	l.Log("success", fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message.
func (l *LoggerZ[T]) Warnf(format string, args ...any) {
	l.Log("warn", fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message and returns an error.
func (l *LoggerZ[T]) Errorf(format string, args ...any) error {
	return l.Log("error", fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal message.
func (l *LoggerZ[T]) Fatalf(format string, args ...any) {
	l.Log("fatal", fmt.Sprintf(format, args...))
}

// Tracef logs a formatted trace message.
func (l *LoggerZ[T]) Tracef(format string, args ...any) {
	l.Log("trace", fmt.Sprintf(format, args...))
}

// Criticalf logs a formatted critical message.
func (l *LoggerZ[T]) Criticalf(format string, args ...any) {
	l.Log("critical", fmt.Sprintf(format, args...))
}

// Answerf logs a formatted answer message.
func (l *LoggerZ[T]) Answerf(format string, args ...any) {
	l.Log("answer", fmt.Sprintf(format, args...))
}

// Alertf logs a formatted alert message.
func (l *LoggerZ[T]) Alertf(format string, args ...any) {
	l.Log("alert", fmt.Sprintf(format, args...))
}

// Bugf logs a formatted bug message.
func (l *LoggerZ[T]) Bugf(format string, args ...any) {
	l.Log("bug", fmt.Sprintf(format, args...))
}

// Panicf logs a formatted panic message.
func (l *LoggerZ[T]) Panicf(format string, args ...any) {
	l.Log("panic", fmt.Sprintf(format, args...))
}
