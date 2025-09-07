package formatters

import (
	il "github.com/kubex-ecosystem/logz/internal/core"
)

// TableFormatter formats data as a table.
type TableFormatter = il.TableFormatter

func NewTableFormatter() il.LogFormatter {
	return &TableFormatter{}
}

func NewTableFormatterExt() *il.TableFormatter {
	return &TableFormatter{}
}
