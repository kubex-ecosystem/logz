package main

import (
	"github.com/kubex-ecosystem/logz"
)

func main() {
	// Teste inicial com modo debug DESABILITADO (padrão é Info)
	logz.Log("info", "=== Modo Normal (Info+) ===")
	logz.Log("debug", "Este debug NÃO deve aparecer")
	logz.Log("info", "Este info deve aparecer")
	logz.Log("warn", "Este warn deve aparecer")
	logz.Log("error", "Este error deve aparecer")

	// Agora vamos HABILITAR o modo debug
	logz.SetDebugMode(true)
	logz.Log("info", "\n=== Modo Debug ATIVADO ===")
	logz.Log("debug", "Agora este debug DEVE aparecer!")
	logz.Log("info", "Este info continua aparecendo")
	logz.Log("warn", "Este warn continua aparecendo")

	// Agora vamos DESABILITAR o modo debug novamente
	logz.SetDebugMode(false)
	logz.Log("info", "\n=== Modo Normal RESTAURADO ===")
	logz.Log("debug", "Este debug NÃO deve aparecer novamente")
	logz.Log("info", "Este info continua aparecendo")
}
