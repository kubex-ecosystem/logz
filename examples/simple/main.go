package main

import (
	"github.com/kubex-ecosystem/logz"
)

func main() {
	// Uso mais simples possível do logz
	// Sem configuração, sem constructor, sem nada!

	if err := logz.Log("info", "testando o logz...", "pra", "ver", "Funcionando!!"); err != nil {
		panic(err)
	}

	if err := logz.Log("warn", "Atenção:", "Este é um aviso"); err != nil {
		panic(err)
	}

	if err := logz.Log("error", "Algo deu errado!", "Código:", 500); err != nil {
		panic(err)
	}

	if err := logz.Log("debug", "Debug info", "variável x =", 42); err != nil {
		panic(err)
	}
}
