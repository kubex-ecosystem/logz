package main

import (
	"github.com/kubex-ecosystem/logz"
)

// Este arquivo pode ser executado separadamente para testar o Fatal
// go run fatal_test.go
// O programa deve exibir a mensagem fatal e encerrar com exit code 1

func main() {
	logz.Info("Iniciando teste do Fatal...")
	logz.Warn("O programa vai encerrar após a próxima mensagem")
	logz.Fatal("Erro crítico detectado!", "Encerrando aplicação...")
	logz.Info("Esta mensagem NUNCA deve aparecer")
}
