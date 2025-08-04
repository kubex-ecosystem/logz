package formatters

import (
	il "github.com/rafa-mori/logz/internal/core"
)

// TextFormatter formats log entries in plain text.
type TextFormatter = il.TextFormatter

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}
