package kbx

import "strings"

// Level representa o nível semântico do log.
// Mantemos string para interoperabilidade humana/JSON.
type Level string

const (
	LevelNotice   Level = "notice"
	LevelDebug    Level = "debug"
	LevelTrace    Level = "trace"
	LevelSuccess  Level = "success"
	LevelInfo     Level = "info"
	LevelWarn     Level = "warn"
	LevelError    Level = "error"
	LevelFatal    Level = "fatal"
	LevelSilent   Level = "silent"
	LevelAlert    Level = "alert"
	LevelCritical Level = "critical"
	LevelAnswer   Level = "answer"
	LevelBug      Level = "bug"
	LevelPanic    Level = "panic"

	LevelSprintf Level = "sprintf"
	LevelPrintln Level = "println"
)

func (l Level) String() string { return string(l) }

// Severity retorna uma escala numérica estável, usada para filtragem.
// Quanto maior, mais grave. Silent = 0.
func (l Level) Severity() int {
	switch strings.ToLower(strings.ToValidUTF8(string(l), "")) {
	case string(LevelSilent):
		return 0
	case string(LevelAnswer):
		return 1
	case string(LevelTrace):
		return 5
	case string(LevelNotice):
		return 10
	case string(LevelDebug):
		return 15
	case string(LevelInfo):
		return 20
	case string(LevelSuccess):
		return 25
	case string(LevelWarn):
		return 30
	case string(LevelAlert):
		return 35
	case string(LevelError):
		return 40
	case string(LevelCritical):
		return 45
	case string(LevelFatal):
		return 50
	case string(LevelBug):
		return 60
	case string(LevelPanic):
		return 70
	default:
		return 5
	}
}

// ParseLevel converte string em Level, com fallback para info.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(strings.ToValidUTF8(s, ""))) {
	case string(LevelInfo):
		return LevelInfo
	case string(LevelWarn):
		return LevelWarn
	case string(LevelError):
		return LevelError
	case string(LevelFatal):
		return LevelFatal
	case string(LevelSilent):
		return LevelSilent
	case string(LevelDebug):
		return LevelDebug
	case string(LevelTrace):
		return LevelTrace
	case string(LevelNotice):
		return LevelNotice
	case string(LevelSuccess):
		return LevelSuccess
	case string(LevelAlert):
		return LevelAlert
	case string(LevelCritical):
		return LevelCritical
	case string(LevelAnswer):
		return LevelAnswer
	case string(LevelBug):
		return LevelBug
	case string(LevelPanic):
		return LevelPanic
	default:
		return LevelSilent
	}
}

func IsLevel(s string) bool {
	switch strings.ToLower(strings.TrimSpace(strings.ToValidUTF8(s, ""))) {
	case string(LevelInfo),
		string(LevelWarn),
		string(LevelError),
		string(LevelFatal),
		string(LevelSilent),
		string(LevelDebug),
		string(LevelTrace),
		string(LevelNotice),
		string(LevelSuccess),
		string(LevelAlert),
		string(LevelCritical),
		string(LevelAnswer),
		string(LevelBug),
		string(LevelPanic):
		return true
	default:
		return false
	}
}
