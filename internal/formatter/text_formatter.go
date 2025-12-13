package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type TextFormatter struct {
	DisableColor bool
	DisableIcon  bool
}

// --- ICONES E CORES ---------------------------------------------------------

var icons = map[kbx.Level]string{
	kbx.LevelAlert:    "🚨",
	kbx.LevelAnswer:   "💡",
	kbx.LevelNotice:   "📝",
	kbx.LevelTrace:    "🔍",
	kbx.LevelSuccess:  "✅",
	kbx.LevelDebug:    "🐛",
	kbx.LevelInfo:     "ℹ️",
	kbx.LevelWarn:     "⚠️",
	kbx.LevelError:    "❌",
	kbx.LevelFatal:    "💀",
	kbx.LevelPanic:    "🔥",
	kbx.LevelBug:      "🐞",
	kbx.LevelCritical: "❗",
}

var colors = map[kbx.Level]string{
	kbx.LevelAlert:    "\033[31m",
	kbx.LevelAnswer:   "\033[34m",
	kbx.LevelNotice:   "\033[33m",
	kbx.LevelTrace:    "\033[36m",
	kbx.LevelSuccess:  "\033[32m",
	kbx.LevelDebug:    "\033[34m",
	kbx.LevelInfo:     "\033[32m",
	kbx.LevelWarn:     "\033[33m",
	kbx.LevelError:    "\033[31m",
	kbx.LevelFatal:    "\033[35m",
	kbx.LevelPanic:    "\033[31m",
	kbx.LevelBug:      "\033[31m",
	kbx.LevelCritical: "\033[31m",
}

const reset = "\033[0m"

// --- CONSTRUCTOR ------------------------------------------------------------

func NewTextFormatter(pretty bool) Formatter {
	return &TextFormatter{
		DisableColor: os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows" || !pretty,
		DisableIcon:  os.Getenv("LOGZ_NO_ICON") != "" || !pretty,
	}
}

// --- PUBLIC API -------------------------------------------------------------

func (f *TextFormatter) Name() string {
	return "text"
}

func (f *TextFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	msg := strings.TrimSpace(e.GetMessage())

	// Level string
	levelStr := string(e.GetLevel())
	if !f.DisableColor && e.GetShowColor() {
		if c, ok := colors[e.GetLevel()]; ok {
			levelStr = c + levelStr + reset
		}
	}

	// Icon
	icon := ""
	if !f.DisableIcon && e.GetShowIcon() {
		if ic, ok := icons[e.GetLevel()]; ok {
			icon = ic + " "
		}
	}

	// Timestamp (opcional)
	ts := ""
	if e.GetTimestamp().Unix() != 0 {
		ts = e.GetTimestamp().Format("2006-01-02 15:04:05")
		ts = "[" + ts + "] "
	}

	// Context
	ctx := ""
	if c := e.GetContext(); c != "" {
		ctx = "(" + c + ") "
	}

	// Metadata (bonitinho e opcional)
	meta := ""
	if m := e.GetTags(); len(m) > 0 && e.GetShowStack() {
		b, _ := json.MarshalIndent(m, "", "  ")
		meta = "\n" + string(b)
	}

	// Fields (bonitinho e opcional)
	if m := e.GetFields(); len(m) > 1 && e.GetShowFields() {
		b, _ := json.MarshalIndent(m, "", "  ")
		meta += "\n" + string(b)
	}

	// TraceID (opcional)
	if e.GetTraceID() != "" && e.GetShowTraceID() {
		meta += fmt.Sprintf("\nTraceID: %s", e.GetTraceID())
	}

	// Caller (opcional)
	if e.GetShowCaller() || e.GetShowStack() {
		meta += fmt.Sprintf("\nCaller: %s", e.GetCaller())
	}

	// Line final → limpa, previsível, sem comer whitespace
	line := fmt.Sprintf("%s[%s] %s%s %s%s",
		ts,
		levelStr,
		ctx,
		icon,
		msg,
		meta,
	)

	return []byte(line), nil
}
