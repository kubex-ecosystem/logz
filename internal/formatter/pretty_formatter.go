// Package formatter é um formatador de log que produz saídas legíveis por humanos.
package formatter

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type PrettyFormatter struct {
	TimeLayout string
	WithColors bool
}

func NewPrettyFormatter(pretty bool) Formatter {
	return &PrettyFormatter{
		TimeLayout: "15:04:05.000",
		WithColors: pretty,
	}
}

func (f *PrettyFormatter) Name() string {
	return "pretty"
}

func (f *PrettyFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	ts := e.GetTimestamp().In(time.Local).Format(f.TimeLayout)
	level := string(e.GetLevel())
	msg := e.GetMessage()

	levelStr := level
	if f.WithColors || e.GetShowColor() {
		levelStr = colorForLevel(e.GetLevel(), level)
	}

	fmt.Fprintf(&buf, "%s  %s  %s", ts, levelStr, msg)
	if e.GetContext() != "" {
		fmt.Fprintf(&buf, "  (%s)", e.GetContext())
	}
	buf.WriteByte('\n')

	// tags e fields em linhas subsequentes
	if e.GetShowStack() || e.GetShowCaller() || e.GetShowFields() {
		if len(e.GetTags()) > 0 || len(e.GetFields()) > 0 || e.GetCaller() != "" {
			keys := make([]string, 0, len(e.GetTags()))
			for k := range e.GetTags() {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			if len(keys) > 0 {
				buf.WriteString("  tags: ")
				for i, k := range keys {
					if i > 0 {
						buf.WriteString(" ")
					}
					fmt.Fprintf(&buf, "%s=%s", k, e.GetTags()[k])
				}
				buf.WriteByte('\n')
			}

			fkeys := make([]string, 0, len(e.GetFields()))
			for k := range e.GetFields() {
				fkeys = append(fkeys, k)
			}
			sort.Strings(fkeys)

			if len(fkeys) > 0 {
				buf.WriteString("  fields: ")
				for i, k := range fkeys {
					if i > 0 {
						buf.WriteString(" ")
					}
					fmt.Fprintf(&buf, "%s=%v", k, e.GetFields()[k])
				}
				buf.WriteByte('\n')
			}

			if e.GetCaller() != "" {
				fmt.Fprintf(&buf, "  caller: %s\n", e.GetCaller())
			}
		}
	}

	return buf.Bytes(), nil
}

func colorForLevel(l kbx.Level, s string) string {
	const (
		reset   = "\x1b[0m"
		gray    = "\x1b[90m"
		green   = "\x1b[32m"
		yellow  = "\x1b[33m"
		red     = "\x1b[31m"
		magenta = "\x1b[35m"
		cyan    = "\x1b[36m"
		white   = "\x1b[37m"
		blue    = "\x1b[34m"
		purple  = "\x1b[35m"
		redBG   = "\x1b[41m"
	)
	switch l {
	case kbx.LevelDebug:
		return gray + s + reset
	case kbx.LevelInfo:
		return green + s + reset
	case kbx.LevelWarn:
		return yellow + s + reset
	case kbx.LevelError:
		return red + s + reset
	case kbx.LevelFatal:
		return magenta + s + reset
	case kbx.LevelPanic:
		return purple + s + reset
	case kbx.LevelTrace:
		return cyan + s + reset
	case kbx.LevelNotice:
		return white + s + reset
	case kbx.LevelAlert:
		return redBG + blue + s + reset
	default:
		return s
	}
}
