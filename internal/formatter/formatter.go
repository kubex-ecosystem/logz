package formatter

import (
	"github.com/kubex-ecosystem/logz/interfaces"
)

func ParseFormatter(format string, pretty bool) interfaces.Formatter {
	switch format {
	case "json":
		return NewJSONFormatter(pretty)
	case "text":
		return NewTextFormatter(pretty)
	default:
		return NewTextFormatter(pretty)
	}
}
