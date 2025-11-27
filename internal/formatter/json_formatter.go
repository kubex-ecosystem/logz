// Package formatter provides a JSON formatter for log entries.
package formatter

import (
	"encoding/json"

	"github.com/kubex-ecosystem/logz/interfaces"
)

type JSONFormatter struct {
	Pretty bool
}

func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{Pretty: pretty}
}

func (f *JSONFormatter) Format(e interfaces.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if f.Pretty {
		return json.MarshalIndent(e, "", "  ")
	}
	return json.Marshal(e)
}
