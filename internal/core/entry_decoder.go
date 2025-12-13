package core

import (
	"strings"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// DefaultEntryDecoder cria uma função Decode pra IOBridge[*Entry],
// que transforma uma linha de texto em uma Entry simples.
//
// É uma estratégia padrão: level fixo + mensagem = linha inteira.
func DefaultEntryDecoder(defaultLevel kbx.Level) func([]byte) (kbx.LogzEntry, error) {
	if defaultLevel == "" {
		defaultLevel = kbx.LevelInfo
	}

	return func(p []byte) (kbx.LogzEntry, error) {
		msg := strings.TrimSpace(string(p))
		if msg == "" {
			return nil, nil
		}

		return NewLogzEntry(defaultLevel).
			WithMessage(msg), nil
	}
}
