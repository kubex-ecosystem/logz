package core

import (
	"strings"

	"github.com/kubex-ecosystem/logz/interfaces"
)

// IOBridge é o adaptador que IMPLEMENTA io.Writer e empurra tudo para um Logger[T].
//
// É aqui que o "modo B" (byte-first) entra no modo C híbrido:
// qualquer coisa que escreva em io.Writer pode ser redirecionada para o logger.
//
// Exemplo de uso (lá fora, na superfície pública):
//
//	entryLogger := core.NewLogger[*core.Entry](formatter, os.Stdout, core.LevelInfo)
//	bridge      := core.NewIOBridge(entryLogger, core.DefaultEntryDecoder(core.LevelInfo))
//
//	log.SetOutput(bridge)      // log stdlib
//	cmd.Stdout = bridge        // exec.Command
//	json.NewEncoder(bridge)...
type IOBridge[T interfaces.Entry] struct {
	Logger *LoggerZ[T]
	Decode func([]byte) (T, error)
}

// NewIOBridge cria a ponte genérica entre io.Writer e Logger[T].
func NewIOBridge[T interfaces.Entry](logger *LoggerZ[T], decode func([]byte) (T, error)) *IOBridge[T] {
	return &IOBridge[T]{
		Logger: logger,
		Decode: decode,
	}
}

// Write implementa io.Writer.
//
// Estratégia básica (simples e eficaz):
// - recebe chunk de bytes
// - chama Decode pra produzir um T
// - passa pro Logger.Log
//
// Se Logger ou Decode forem nil, a escrita é "dropada" mas não quebra a app.
func (b *IOBridge[T]) Write(p []byte) (int, error) {
	if b == nil || b.Logger == nil || b.Decode == nil {
		// absorve a escrita pra não quebrar o chamador.
		return len(p), nil
	}

	trimmed := strings.TrimSpace(string(p))
	if trimmed == "" {
		return len(p), nil
	}

	rec, err := b.Decode([]byte(trimmed))
	if err != nil {
		return 0, err
	}

	if err := b.Logger.Log(string(rec.GetLevel()), rec); err != nil {
		return 0, err
	}

	return len(p), nil
}
