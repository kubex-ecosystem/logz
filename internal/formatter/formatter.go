package formatter

import "github.com/kubex-ecosystem/logz/interfaces"

// Formatter transforma uma Entry em bytes prontos pra IO.
// NÃO escreve em lugar nenhum. Só formata.
type Formatter interface {
	Format(interfaces.Entry) ([]byte, error)
}
