// Package api provides the API for logging to Logz.io.
package api

import (
	il "github.com/rafa-mori/logz/internal/core"
)

// LogzEntry is an alias for the internal LogzEntry type.
type LogzEntry = il.LogzEntry

// NewLogEntry creates a new LogzEntry.
func NewLogEntry() LogzEntry {
	return il.NewLogEntry()
}
