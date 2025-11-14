// Package logger provides a unified interface for logging with various configurations and formats.
package logger

// import (
// 	"os"

// 	apiW "github.com/kubex-ecosystem/logz/api/writers"
// 	il "github.com/kubex-ecosystem/logz/"
// )

// // type ILogWriter[T any, W il.WriteLog[T]] interface{ il.LogWriter[T,W] }

// // type LogFormat = LogFormat
// type LogWriter[T any, W apiW.WriteLogz[T]] = il.LogzWriter[T, W]
// type LogzWriter[E any, W WriteLogz[E]] = il.LogzWriter[E, W]
// type Writer = il.Writer
// type Config = il.Config
// type LogzEntry = il.LogzEntry
// type LogFormatter = il.LogFormatter
// type logxLogger = il.LogzCoreImpl
// type LogzLogger = il.LogzLogger

// func NewLogger(prefix string) LogzLogger { return il.NewLogger(prefix) }
// func NewDefaultWriter(out *os.File, formatter LogFormatter) *il.DefaultWriter[any] {
// 	return il.NewDefaultWriter[any](out, formatter)
// }
