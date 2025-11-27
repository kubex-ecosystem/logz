package interfaces

import "strings"

// Level representa o nível semântico do log.
// Mantemos string para interoperabilidade humana/JSON.
type Level string

const (
	LevelDebug  Level = "debug"
	LevelInfo   Level = "info"
	LevelWarn   Level = "warn"
	LevelError  Level = "error"
	LevelFatal  Level = "fatal"
	LevelSilent Level = "silent"
)

func (l Level) String() string {
	return string(l)
}

// Severity retorna uma escala numérica estável, usada para filtragem.
// Quanto maior, mais grave. Silent = 0.
func (l Level) Severity() int {
	switch strings.ToLower(string(l)) {
	case "debug":
		return 10
	case "info":
		return 20
	case "warn", "warning":
		return 30
	case "error", "err":
		return 40
	case "fatal":
		return 50
	case "silent", "":
		return 0
	default:
		// fallback seguro: trata como info
		return 20
	}
}

// ParseLevel converte string em Level, com fallback para info.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error", "err":
		return LevelError
	case "fatal":
		return LevelFatal
	case "silent", "":
		return LevelSilent
	default:
		return LevelInfo
	}
}
