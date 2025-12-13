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

	LevelSprintf   Level = "sprintf"
	LevelPrintln   Level = "println"
	LevelLog       Level = "log"
	LevelErrorf    Level = "errorf"
	LevelWarnf     Level = "warnf"
	LevelInfof     Level = "infof"
	LevelDebugf    Level = "debugf"
	LevelFatalf    Level = "fatalf"
	LevelPanicf    Level = "panicf"
	LevelAlertf    Level = "alertf"
	LevelCriticalf Level = "criticalf"
	LevelAnswerf   Level = "answerf"
	LevelBugf      Level = "bugf"
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
	case "silent", "quiet":
		return 0
	case "trace":
		return 5
	case "notice":
		return 15
	case "success":
		return 25
	case "alert":
		return 35
	case "critical":
		return 45
	case "answer":
		return 55
	case "bug":
		return 60
	case "panic":
		return 70
	default:
		// fallback seguro: trata como info
		return 20
	}
}

// ParseLevel converte string em Level, com fallback para info.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error", "err":
		return LevelError
	case "fatal":
		return LevelFatal
	case "silent", "quiet":
		return LevelSilent
	case "debug":
		return LevelDebug
	case "trace":
		return LevelTrace
	case "notice":
		return LevelNotice
	case "success":
		return LevelSuccess
	case "alert":
		return LevelAlert
	case "critical":
		return LevelCritical
	case "answer":
		return LevelAnswer
	case "bug":
		return LevelBug
	case "panic":
		return LevelPanic
	default:
		return LevelDebug
	}
}

func IsLevel(s string) bool {
	switch strings.ToLower(strings.TrimSpace(strings.ToValidUTF8(s, ""))) {
	case "debug", "info", "warn", "warning", "error", "err", "fatal", "silent":
		return true
	default:
		return false
	}
}
