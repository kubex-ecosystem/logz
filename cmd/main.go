package main

import (
	"fmt"
	"os"

	"github.com/kubex-ecosystem/logz/internal/module"
)

func main() {
	if logzErr := module.RegX().Execute(); logzErr != nil {
		fmt.Printf("ErrorCtx executing command: %v\n", logzErr)
		os.Exit(1)
	}
}
