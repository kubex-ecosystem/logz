package formatter

import (
	"encoding/json"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type JSONFormatter struct {
	Pretty bool
}

func NewJSONFormatter(pretty bool) Formatter {
	return &JSONFormatter{Pretty: pretty}
}

func (f *JSONFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if f.Pretty {
		return json.MarshalIndent(e, "", "  ")
	}
	return json.Marshal(e)
}
