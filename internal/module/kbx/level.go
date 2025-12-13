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
	LevelErrorf    Level = "error"
	LevelWarnf     Level = "warn"
	LevelInfof     Level = "info"
	LevelDebugf    Level = "debug"
	LevelFatalf    Level = "fatal"
	LevelPanicf    Level = "panic"
	LevelAlertf    Level = "alert"
	LevelCriticalf Level = "critical"
	LevelAnswerf   Level = "answer"
	LevelBugf      Level = "bug"
)

func (l Level) String() string {
	return string(l)
}

// Severity retorna uma escala numérica estável, usada para filtragem.
// Quanto maior, mais grave. Silent = 0.
func (l Level) Severity() int {
	switch strings.ToLower(string(l)) {
	case "debug","debug*":
		return 10
	case "info","info*":
		return 20
	case "warn", "warn*", "warning":
		return 30
	case "error", "error*", "err":
		return 40
	case "fatal", "fatal*":
		return 50
	case "silent", "silent*", "quiet", "quiet*":
		return 0
	case "trace", "trace*":
		return 5
	case "notice", "notice*":
		return 15
	case "success", "success*":
		return 25
	case "alert", "alert*":
		return 35
	case "critical", "critical*":
		return 45
	case "answer", "answer*":
		return 55
	case "bug", "bug*":
		return 60
	case "panic", "panic*":
		return 70
	default:
		// fallback seguro: trata como info
		return 20
	}
}

// ParseLevel converte string em Level, com fallback para info.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info", "info*":
		return LevelInfo
	case "warn", "warn*", "warning":
		return LevelWarn
	case "error", "error*", "err":
		return LevelError
	case "fatal", "fatal*":
		return LevelFatal
	case "silent", "silent*", "quiet", "quiet*":
		return LevelSilent
	case "debug", "debug*":
		return LevelDebug
	case "trace", "trace*":
		return LevelTrace
	case "notice", "notice*":
		return LevelNotice
	case "success", "success*":
		return LevelSuccess
	case "alert", "alert*":
		return LevelAlert
	case "critical", "critical*":
		return LevelCritical
	case "answer", "answer*":
		return LevelAnswer
	case "bug", "bug*":
		return LevelBug
	case "panic", "panic*":
		return LevelPanic
	default:
		return LevelInfo
	}
}

func IsLevel(s string) bool {
	switch strings.ToLower(strings.TrimSpace(strings.ToValidUTF8(s, ""))) {
	case "debug", "debug*", "info", "info*", "warn", "warn*", "warning", "error",
		"error*", "err", "fatal", "fatal*", "silent", "silent*", "quiet", "quiet*",
		"trace", "trace*", "notice", "notice*", "success", "success*", "alert", "alert*",
		"critical", "critical*", "answer", "answer*", "bug", "bug*", "panic", "panic*":
		return true
	default:
		return false
	}
}
