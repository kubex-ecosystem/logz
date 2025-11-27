package core

import (
	"strings"

	"github.com/kubex-ecosystem/logz/interfaces"
)

// DefaultEntryDecoder cria uma função Decode pra IOBridge[*Entry],
// que transforma uma linha de texto em uma Entry simples.
//
// É uma estratégia padrão: level fixo + mensagem = linha inteira.
func DefaultEntryDecoder(defaultLevel interfaces.Level) func([]byte) (interfaces.Entry, error) {
	if defaultLevel == "" {
		defaultLevel = interfaces.LevelInfo
	}

	return func(p []byte) (interfaces.Entry, error) {
		msg := strings.TrimSpace(string(p))
		if msg == "" {
			return nil, nil
		}
		e, er := NewEntry()
		if er != nil {
			return nil, er
		}
		return e.
			WithLevel(defaultLevel).
			WithMessage(msg), nil
	}
}
