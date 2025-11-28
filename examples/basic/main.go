package main

import (
	"os"

	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

func main() {
	// Exemplo 1: Logger básico com formato texto
	opts := &core.LoggerOptionsImpl{
		LoggerConfig: &core.LoggerConfig{
			ID: kbx.LoggerArgs.ID,
			LogzGeneralOptions: &kbx.LogzGeneralOptions{
				Prefix: "APP",
			},
			LogzFormatOptions: &kbx.LogzFormatOptions{
				Output:   os.Stdout,
				Level:    kbx.LevelInfo,
				MinLevel: kbx.LevelDebug,
				MaxLevel: kbx.LevelFatal,
			},
			LogzOutputOptions:    &kbx.LogzOutputOptions{},
			LogzRotatingOptions:  &kbx.LogzRotatingOptions{},
			LogzBufferingOptions: &kbx.LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &core.LogzAdvancedOptions{},
	}
	opts.Formatter = formatter.NewTextFormatter(false)

	logger := core.NewLogger("ExampleApp", opts, false)

	// Criar entries e logar usando métodos chainable
	entry1, _ := core.NewEntry(kbx.LevelInfo)
	entry1 = entry1.WithMessage("Aplicação iniciada com sucesso")
	logger.Log(kbx.LevelInfo, entry1)

	entry2, _ := core.NewEntry(kbx.LevelWarn)
	entry2 = entry2.WithMessage("Memória acima de 80%")
	logger.Log(kbx.LevelWarn, entry2)

	entry3, _ := core.NewEntry(kbx.LevelError)
	entry3 = entry3.WithMessage("Falha ao conectar ao banco de dados")
	logger.Log(kbx.LevelError, entry3)

	// Exemplo 2: Logger com formato JSON
	optsJSON := &core.LoggerOptionsImpl{
		LoggerConfig: &core.LoggerConfig{
			ID: kbx.LoggerArgs.ID,
			LogzGeneralOptions: &kbx.LogzGeneralOptions{
				Prefix: "API",
			},
			LogzFormatOptions: &kbx.LogzFormatOptions{
				Output:   os.Stdout,
				Level:    kbx.LevelDebug,
				MinLevel: kbx.LevelDebug,
				MaxLevel: kbx.LevelFatal,
			},
			LogzOutputOptions:    &kbx.LogzOutputOptions{},
			LogzRotatingOptions:  &kbx.LogzRotatingOptions{},
			LogzBufferingOptions: &kbx.LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &core.LogzAdvancedOptions{},
	}
	optsJSON.Formatter = formatter.NewJSONFormatter(false)

	loggerJSON := core.NewLogger("API", optsJSON, false)

	entry4, _ := core.NewEntry(kbx.LevelDebug)
	entry4 = entry4.
		WithMessage("Request recebido").
		WithField("method", "GET").
		WithField("path", "/api/users").
		WithField("ip", "192.168.1.100")
	loggerJSON.Log(kbx.LevelDebug, entry4)

	entry5, _ := core.NewEntry(kbx.LevelSuccess)
	entry5 = entry5.
		WithMessage("Response enviado").
		WithField("status", 200).
		WithField("time", "45ms")
	loggerJSON.Log(kbx.LevelSuccess, entry5)
}
