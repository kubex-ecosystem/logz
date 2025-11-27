package main

import (
	"fmt"
	"os"

	"github.com/kubex-ecosystem/logz/internal/module"
	manifest "github.com/kubex-ecosystem/logz/internal/module/info"
)

var info, err = manifest.GetManifest()

func init() {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading manifest: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if logzErr := module.RegX().Execute(); logzErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", logzErr)
		os.Exit(1)
	}
}
