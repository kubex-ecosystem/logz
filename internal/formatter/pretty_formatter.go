package formatter

import (
	"bytes"

	"github.com/kubex-ecosystem/logz/interfaces"
)

type PrettyFormatter struct {
	TimeLayout string
	WithColors bool
}

func NewPrettyFormatter() *PrettyFormatter {
	return &PrettyFormatter{
		TimeLayout: "15:04:05.000",
		WithColors: true,
	}
}

func (f *PrettyFormatter) Format(e interfaces.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// ts := e.Timestamp.In(time.Local).Format(f.TimeLayout)
	// level := string(e.Level)

	// levelStr := level
	// if f.WithColors {
	// 	levelStr = colorForLevel(e.Level, level)
	// }

	// fmt.Fprintf(&buf, "%s  %s  %s", ts, levelStr, e.Message)
	// if e.Context != "" {
	// 	fmt.Fprintf(&buf, "  (%s)", e.Context)
	// }
	// buf.WriteByte('\n')

	// // tags e fields em linhas subsequentes
	// if len(e.Tags) > 0 || len(e.Fields) > 0 || e.Caller != "" {
	// 	keys := make([]string, 0, len(e.Tags))
	// 	for k := range e.Tags {
	// 		keys = append(keys, k)
	// 	}
	// 	sort.Strings(keys)

	// 	if len(keys) > 0 {
	// 		buf.WriteString("  tags: ")
	// 		for i, k := range keys {
	// 			if i > 0 {
	// 				buf.WriteString(" ")
	// 			}
	// 			fmt.Fprintf(&buf, "%s=%s", k, e.Tags[k])
	// 		}
	// 		buf.WriteByte('\n')
	// 	}

	// 	fkeys := make([]string, 0, len(e.Fields))
	// 	for k := range e.Fields {
	// 		fkeys = append(fkeys, k)
	// 	}
	// 	sort.Strings(fkeys)

	// 	if len(fkeys) > 0 {
	// 		buf.WriteString("  fields: ")
	// 		for i, k := range fkeys {
	// 			if i > 0 {
	// 				buf.WriteString(" ")
	// 			}
	// 			fmt.Fprintf(&buf, "%s=%v", k, e.Fields[k])
	// 		}
	// 		buf.WriteByte('\n')
	// 	}

	// 	if e.Caller != "" {
	// 		fmt.Fprintf(&buf, "  caller: %s\n", e.Caller)
	// 	}
	// }

	return buf.Bytes(), nil
}

func colorForLevel(l string, s string) string {
	const (
		reset   = "\x1b[0m"
		gray    = "\x1b[90m"
		green   = "\x1b[32m"
		yellow  = "\x1b[33m"
		red     = "\x1b[31m"
		magenta = "\x1b[35m"
	)
	switch l {
	case "debug":
		return gray + s + reset
	case "info":
		return green + s + reset
	case "warn":
		return yellow + s + reset
	case "error":
		return red + s + reset
	case "fatal":
		return magenta + s + reset
	default:
		return s
	}
}
