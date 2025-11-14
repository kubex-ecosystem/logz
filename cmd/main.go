package main

import (
	"os"

	gl "github.com/kubex-ecosystem/logz"

	"github.com/kubex-ecosystem/logz/internal/module"
	manifest "github.com/kubex-ecosystem/logz/internal/module/info"
)

var info, err = manifest.GetManifest()

func init() {
	gl.GetLogger(info.GetName()).Info("Initializing Logz CLI Application...")
}

func main() {
	if logzErr := module.RegX().Execute(); logzErr != nil {
		gl.GetLogger(info.GetName()).Error("ErrorCtx executing command: %v", logzErr)
		os.Exit(1)
	}
}
