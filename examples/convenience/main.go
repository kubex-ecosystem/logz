package main

import (
	"github.com/kubex-ecosystem/logz"
)

func main() {
	// Habilita debug mode para ver todos os logs
	logz.SetDebugMode(true)

	logz.Info("=== Testando métodos de conveniência do Logz ===")

	// Debug
	logz.Debug("Mensagem de debug", "variável x =", 42)

	// Notice
	logz.Notice("Notificação importante", "sistema iniciado")

	// Info
	logz.Info("Informação geral", "processo em execução")

	// Success
	logz.Success("Operação completada com sucesso!", "registros:", 150)

	// Warn
	logz.Warn("Atenção!", "Memória em", "85%")

	// Error (retorna error)
	if err := logz.Error("Erro ao processar", "código:", 500); err != nil {
		logz.Info("Error retornou:", err)
	}

	// Fatal (este será o último - programa encerra com exit 1)
	// Descomente a linha abaixo para testar:
	// logz.Fatal("Erro crítico!", "Encerrando aplicação...")

	logz.Info("Fim dos testes (Fatal está comentado para não encerrar)")
}
