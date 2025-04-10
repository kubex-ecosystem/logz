package main

import (
	"fmt"
	"os"
)

func main() {
	if logzErr := RegX().Execute(); logzErr != nil {
		fmt.Printf("ErrorCtx executing command: %v\n", logzErr)
		os.Exit(1)
	}
}
