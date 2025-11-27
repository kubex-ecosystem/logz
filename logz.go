// Package logz provides a global logging utility with configurable settings.
package logz

// 🧪 Como testar rápido

// Depois de jogar isso no repo:

// package main

// import (
// 	"os"

// 	"github.com/kubex-ecosystem/logz"
// )

// func main() {
// 	logz.Info("hello kubex",
// 		logz.FContext("boot"),
// 		logz.FTag("env", "dev"),
// 		logz.FField("attempt", 1),
// 	)

// 	logz.Entry().
// 		Level(logz.LevelInfo).
// 		Context("auth.login").
// 		Tag("user", "42").
// 		Field("duration_ms", 91).
// 		Msg("user logged in")

// 	f, _ := os.Create("logs.json")
// 	logz.SetGlobalWriter(
// 		logz.NewMultiWriter(
// 			logz.NewIOWriter(os.Stdout),
// 			logz.NewIOWriter(f),
// 		),
// 	)
// 	logz.SetGlobalFormatter(logz.NewJSONFormatter(false))

// 	logz.Info("now in json",
// 		logz.FContext("switch"),
// 		logz.FTag("mode", "json"),
// 	)
// }
