package main

import (
	"fmt"
	"os"

	// gl "github.com/kubex-ecosystem/logz"

	"github.com/kubex-ecosystem/logz/internal/module"
	manifest "github.com/kubex-ecosystem/logz/internal/module/info"
)

var info, err = manifest.GetManifest()

// func init() {
// 	gl.GetLogger(info.GetName()).Info("Initializing Logz CLI Application...")
// }

func main() {
	if logzErr := module.RegX().Execute(); logzErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", logzErr)
		os.Exit(1)
	}
}
