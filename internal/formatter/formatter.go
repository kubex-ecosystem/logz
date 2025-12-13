package formatter

import (
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type Formatter interface {
	// Name retorna o nome do formatter.
	Name() string
	Format(e kbx.Entry) ([]byte, error)
}

// FormatterFunc é uma função que implementa a interface Formatter.
type FormatterFunc func(e kbx.Entry) ([]byte, error)

// Name retorna o nome do formatter.
func (f FormatterFunc) Name() string {
	return "custom_formatter_func"
}

// Format formata a entry.
func (f FormatterFunc) Format(e kbx.Entry) ([]byte, error) {
	return f(e)
}

func ParseFormatter(format string, pretty bool) Formatter {
	switch format {
	case "json":
		return NewJSONFormatter(pretty)
	case "text":
		return NewTextFormatter(pretty)
	case "yaml":
		return NewYamlFormatter(pretty)
	case "csv":
		return NewCSVFormatter(pretty)
	case "xml":
		return NewXMLFormatter(pretty)
	default:
		return NewTextFormatter(pretty)
	}
}
