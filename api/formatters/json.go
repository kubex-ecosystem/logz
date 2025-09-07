// Package formatters provides functionality to format data as JSON.
package formatters

import (
	il "github.com/kubex-ecosystem/logz/internal/core"
)

// JSONFormatter formats log entries as JSON.
type JSONFormatter = il.JSONFormatter

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}
