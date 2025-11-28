// Package logz provides a global logging utility with configurable settings.
package logz

import (
	C "github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type Logger = C.LoggerZ[kbx.Entry]
type Entry = kbx.Entry

func NewEntry(level kbx.Level) (kbx.Entry, error) {
	return C.NewEntryImpl(level)
}

func NewLogger(prefix string, opts *C.LoggerOptionsImpl, withDefaults bool) *Logger {
	return C.NewLoggerZ[kbx.Entry](prefix, opts, withDefaults)
}
