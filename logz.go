// Package logz provides a global logging utility with configurable settings.
package logz

import (
	I "github.com/kubex-ecosystem/logz/interfaces"
	C "github.com/kubex-ecosystem/logz/internal/core"
)

type Logger = C.LoggerZ[I.Entry]
type Entry = I.Entry

func NewEntry() (I.Entry, error) {
	return C.NewEntryImpl()
}

func NewLogger(prefix string, opts *C.LoggerOptionsImpl, withDefaults bool) *Logger {
	return C.NewLoggerZ[I.Entry](prefix, opts, withDefaults)
}
