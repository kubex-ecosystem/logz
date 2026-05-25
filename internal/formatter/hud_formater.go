package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// HUDFormatter transforma o fluxo de logs em uma tabela contínua estilo XTUI.
type HUDFormatter struct {
	// Pool para evitar alocações pesadas no caminho crítico
	bufPool sync.Pool
	pretty  bool
}

// NewHUDFormatter creates a new HUDFormatter.
func NewHUDFormatter(pretty bool) *HUDFormatter {
	return &HUDFormatter{
		bufPool: sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		},
		pretty: pretty,
	}
}

// Name retorna o nome do formatter.
func (f *HUDFormatter) Name() string {
	return "hud"
}

// Format implementa a interface do logz (ajuste a assinatura conforme a sua struct)
func (f *HUDFormatter) Format(entry kbx.Entry) ([]byte, error) {
	buf := f.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer f.bufPool.Put(buf)

	// 1. Timestamp (Clean)
	ts := entry.GetTimestamp().Format("15:04:05.000")

	// 2. Configurações visuais dinâmicas baseadas no Level
	var levelColor, icon, levelStr string
	switch entry.GetLevel() {
	case kbx.LevelDebug.String():
		levelColor, icon, levelStr = "\033[36m", "⚙", "DBG" // Cyan
	case kbx.LevelInfo.String():
		levelColor, icon, levelStr = "\033[32m", "●", "INF" // Green
	case kbx.LevelWarn.String():
		levelColor, icon, levelStr = "\033[33m", "▲", "WRN" // Yellow
	case kbx.LevelError.String(), kbx.LevelFatal.String():
		levelColor, icon, levelStr = "\033[31m", "✖", "ERR" // Red
	default:
		levelColor, icon, levelStr = "\033[37m", "○", "LOG" // White
	}
	reset := "\033[0m"
	dim := "\033[90m"

	// 3. Extração e Sparkline (se houver latência nos campos)
	var sparkline string
	var contextFields []string

	for k, v := range entry.GetFields() {
		if k == "latency_ms" || k == "duration_ms" {
			// Gera um gráfico de barras simples baseado no valor numérico
			val := fmt.Sprintf("%v", v)
			sparkline = generateSparkline(val)
		} else {
			contextFields = append(contextFields, fmt.Sprintf("%s%s%s=%v", dim, k, reset, v))
		}
	}

	// 4. Montagem da Linha da Tabela Infinita
	// Estrutura: │ HORA │ LEVEL │ MENSAGEM (padded) │ SPARKLINE │ [CAMPOS]

	msg := entry.GetMessage()
	if len(msg) > 40 {
		msg = msg[:37] + "..." // Truncar para manter a tabela bonita
	}

	// Formatando com alinhamento fixo (% -40s garante o padding da mensagem)
	fmt.Fprintf(buf, "%s│%s %s │ %s%s %s%s │ %-40s │ %s %s %s\n",
		dim, reset, ts,
		levelColor, icon, levelStr, reset,
		msg,
		sparkline,
		strings.Join(contextFields, " "),
		reset,
	)

	// Retorna uma cópia para o fluxo I/O (seguro para o atomic writing que faremos)
	return append([]byte(nil), buf.Bytes()...), nil
}

// generateSparkline cria uma barrinha visual baseada na latência
func generateSparkline(val string) string {
	var ms int
	fmt.Sscanf(val, "%d", &ms)

	switch {
	case ms < 5:
		return "\033[32m ▂▃\033[0m    " // Rápido (Verde)
	case ms < 20:
		return "\033[33m ▂▃▄▅\033[0m  " // Aceitável (Amarelo)
	default:
		return "\033[31m ▂▃▄▅▆▇█\033[0m" // Lento (Vermelho)
	}
}
