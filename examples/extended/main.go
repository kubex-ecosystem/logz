package main

import (
	"errors"
	"os"

	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

func main() {
	// Criar logger com formato JSON para melhor visualização de contexto
	opts := &core.LoggerOptionsImpl{
		LoggerConfig: &core.LoggerConfig{
			ID: kbx.LoggerArgs.ID,
			LogzGeneralOptions: &kbx.LogzGeneralOptions{
				Prefix: "API",
			},
			LogzFormatOptions: &kbx.LogzFormatOptions{
				Output:   os.Stdout,
				Level:    kbx.LevelDebug,
				MinLevel: kbx.LevelInfo,
				MaxLevel: kbx.LevelFatal,
			},
			LogzOutputOptions:    &kbx.LogzOutputOptions{},
			LogzRotatingOptions:  &kbx.LogzRotatingOptions{},
			LogzBufferingOptions: &kbx.LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &core.LogzAdvancedOptions{},
	}
	opts.Formatter = formatter.NewJSONFormatter(true) // pretty print

	logger := core.NewLogger("ExtendedContextExample", opts, false)

	// Exemplo 1: Log com contexto, source e trace ID
	entry1, _ := core.NewEntry(kbx.LevelInfo)
	entry1 = entry1.
		WithMessage("Requisição de autenticação recebida").
		WithContext("auth").
		WithSource("auth-service").
		WithTraceID("trace-123-456-789").
		WithField("username", "admin").
		WithField("ip", "10.0.1.55").
		WithField("method", "POST")
	logger.Log(kbx.LevelInfo, entry1)

	// Exemplo 2: Log com tags e fields
	entry2, _ := core.NewEntry(kbx.LevelWarn)
	entry2 = entry2.
		WithMessage("Tentativa de acesso negado").
		WithContext("security").
		WithSource("rbac-middleware").
		WithTraceID("trace-123-456-789").
		WithField("user_id", 12345).
		WithField("resource", "/api/admin/users").
		WithField("reason", "insufficient permissions")
	// Adicionar tags manualmente
	entry2.Tags["severity"] = "high"
	entry2.Tags["category"] = "security"
	logger.Log(kbx.LevelWarn, entry2)

	// Exemplo 3: Log com erro associado
	entry3, _ := core.NewEntry(kbx.LevelError)
	entry3 = entry3.
		WithMessage("Falha ao processar pagamento").
		WithContext("billing").
		WithSource("payment-service").
		WithTraceID("trace-987-654-321").
		WithField("amount", 99.99).
		WithField("currency", "USD").
		WithField("gateway", "stripe")
	entry3.Tags["transaction_id"] = "txn-abc-123"
	// Adicionar erro
	entry3.Error = errors.New("connection timeout: failed to reach payment gateway")
	logger.Log(kbx.LevelError, entry3)

	// Exemplo 4: Log de sucesso com métricas
	entry4, _ := core.NewEntry(kbx.LevelSuccess)
	entry4 = entry4.
		WithMessage("Sincronização de dados concluída").
		WithContext("sync").
		WithSource("data-sync-service").
		WithTraceID("trace-111-222-333").
		WithField("records_processed", 15000).
		WithField("duration_ms", 3450).
		WithField("success_rate", 99.8)
	entry4.Tags["job_id"] = "sync-20251128"
	logger.Log(kbx.LevelSuccess, entry4)
}
