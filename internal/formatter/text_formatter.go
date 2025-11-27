package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kubex-ecosystem/logz/interfaces"
)

type TextFormatter struct {
	DisableColor bool
	DisableIcon  bool
}

// --- ICONES E CORES ---------------------------------------------------------

var icons = map[interfaces.Level]string{
	interfaces.LevelNotice:  "📝",
	interfaces.LevelTrace:   "🔍",
	interfaces.LevelSuccess: "✅",
	interfaces.LevelDebug:   "🐛",
	interfaces.LevelInfo:    "ℹ️",
	interfaces.LevelWarn:    "⚠️",
	interfaces.LevelError:   "❌",
	interfaces.LevelFatal:   "💀",
}

var colors = map[interfaces.Level]string{
	interfaces.LevelNotice:  "\033[33m",
	interfaces.LevelTrace:   "\033[36m",
	interfaces.LevelSuccess: "\033[32m",
	interfaces.LevelDebug:   "\033[34m",
	interfaces.LevelInfo:    "\033[32m",
	interfaces.LevelWarn:    "\033[33m",
	interfaces.LevelError:   "\033[31m",
	interfaces.LevelFatal:   "\033[35m",
}

const reset = "\033[0m"

// --- CONSTRUCTOR ------------------------------------------------------------

func NewTextFormatter() interfaces.Formatter {
	return &TextFormatter{
		DisableColor: os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows",
		DisableIcon:  os.Getenv("LOGZ_NO_ICON") != "",
	}
}

// --- PUBLIC API -------------------------------------------------------------

func (f *TextFormatter) Format(e interfaces.Entry) ([]byte, error) {
	level := e.GetLevel()
	msg := strings.TrimSpace(e.GetMessage())

	// Level string
	levelStr := string(level)
	if !f.DisableColor {
		if c, ok := colors[level]; ok {
			levelStr = c + levelStr + reset
		}
	}

	// Icon
	icon := ""
	if !f.DisableIcon {
		if ic, ok := icons[level]; ok {
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
	if m := e.GetFields(); len(m) > 0 {
		b, _ := json.MarshalIndent(m, "", "  ")
		meta = "\n" + string(b)
	}

	// Line final → limpa, previsível, sem comer whitespace
	line := fmt.Sprintf("%s[%s] %s%s%s%s",
		ts,
		levelStr,
		ctx,
		icon,
		msg,
		meta,
	)

	return []byte(line), nil
}
