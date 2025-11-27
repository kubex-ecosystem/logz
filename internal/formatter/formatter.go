package formatter

import (
	"github.com/kubex-ecosystem/logz/interfaces"
)

// Formatter transforma uma Entry em bytes prontos pra IO.
// NÃO escreve em lugar nenhum. Só formata.
type Formatter interface {
	Format(interfaces.Entry) ([]byte, error)
}

func ParseFormatter(format string, pretty bool) interfaces.Formatter {
	switch format {
	case "json":
		return NewJSONFormatter(pretty)
	case "text":
		return NewTextFormatter()
	default:
		return NewTextFormatter()
	}
}
