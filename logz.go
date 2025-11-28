// Package logz provides a global logging utility with configurable settings.
package logz

import (
	"os"

	C "github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// defaultLogger cria um logger padrão configurado para uso global
func defaultLogger() *C.Logger {
	opts := &C.LoggerOptionsImpl{
		LoggerConfig: &C.LoggerConfig{
			ID: kbx.LoggerArgs.ID,
			LogzGeneralOptions: &kbx.LogzGeneralOptions{
				Prefix: "",
			},
			LogzFormatOptions: &kbx.LogzFormatOptions{
				Output:   os.Stdout,
				Level:    kbx.LevelInfo,
				MinLevel: kbx.LevelInfo,
				MaxLevel: kbx.LevelFatal,
			},
			LogzOutputOptions:    &kbx.LogzOutputOptions{},
			LogzRotatingOptions:  &kbx.LogzRotatingOptions{},
			LogzBufferingOptions: &kbx.LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &C.LogzAdvancedOptions{},
	}
	opts.Formatter = formatter.NewTextFormatter(false)

	return C.NewLogger("", opts, false)
}

// Logger é a instância global padrão do logger
var Logger = defaultLogger()

type LoggerZ = C.LoggerZ[kbx.Entry]
type Entry = kbx.Entry

func NewEntry(level kbx.Level) (kbx.Entry, error) {
	return C.NewEntryImpl(level)
}

func NewLogger(prefix string, opts *C.LoggerOptionsImpl, withDefaults bool) *C.Logger {
	return C.NewLogger(prefix, opts, withDefaults)
}

func NewLoggerZ(prefix string, opts *C.LoggerOptionsImpl, withDefaults bool) *LoggerZ {
	return C.NewLoggerZ[kbx.Entry](prefix, opts, withDefaults)
}

// Log é a função global mais simples para logging.
// Aceita um level como string e mensagens variádicas.
// Uso: logz.Log("info", "mensagem", "mais", "dados")
func Log(level string, msg ...any) error {
	if Logger == nil {
		return nil
	}
	lvl := kbx.ParseLevel(level)
	return Logger.Log(lvl, msg...)
}

// LogAny é uma variante que aceita qualquer tipo como mensagem
func LogAny(level string, msg any) error {
	if Logger == nil {
		return nil
	}
	lvl := kbx.ParseLevel(level)
	return Logger.LogAny(lvl, msg)
}

// SetDebugMode habilita ou desabilita o modo debug do logger global.
// Quando debug=true, mostra logs de todos os níveis (incluindo debug e trace).
// Quando debug=false, mostra apenas logs de nível info ou superior.
func SetDebugMode(debug bool) {
	if Logger == nil {
		return
	}
	if debug {
		Logger.SetMinLevel(kbx.LevelDebug)
	} else {
		Logger.SetMinLevel(kbx.LevelInfo)
	}
}

// Debug loga uma mensagem de debug
func Debug(msg ...any) {
	Log("debug", msg...)
}

// Notice loga uma mensagem de notice
func Notice(msg ...any) {
	Log("notice", msg...)
}

// Info loga uma mensagem informativa
func Info(msg ...any) {
	Log("info", msg...)
}

// Success loga uma mensagem de sucesso
func Success(msg ...any) {
	Log("success", msg...)
}

// Warn loga um aviso
func Warn(msg ...any) {
	Log("warn", msg...)
}

// Error loga um erro e retorna error
func Error(msg ...any) error {
	return Log("error", msg...)
}

// Fatal loga uma mensagem fatal e encerra o programa com exit code 1
func Fatal(msg ...any) {
	Log("fatal", msg...)
	os.Exit(1)
}
