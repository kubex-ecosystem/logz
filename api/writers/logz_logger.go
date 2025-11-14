// Package api provides the API for logging to Logz.io.
package writers

import (
	"github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/loggerz"
)

// LogzEntry is an alias for the internal LogzEntry type.
type LogzEntryImpl = loggerz.LogEntry
type LogzEntryZ = interfaces.LogzEntry

// NewLogEntry creates a new LogzEntry.
func NewLogEntry() LogzEntryZ {
	return loggerz.NewLogEntry()
}
