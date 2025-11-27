package formatter

import (
	core "github.com/kubex-ecosystem/logz/internal/core"
)

// Formatter transforma uma Entry em bytes prontos pra IO.
// NÃO escreve em lugar nenhum. Só formata.
type Formatter interface {
	Format(*core.Entry) ([]byte, error)
}
