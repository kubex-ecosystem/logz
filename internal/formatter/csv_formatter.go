package formatter

import (
	"fmt"
	"strings"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type CSVFormatter struct {
	Pretty bool
}

func NewCSVFormatter(pretty bool) Formatter {
	return &CSVFormatter{Pretty: pretty}
}

func (f *CSVFormatter) Name() string {
	return "csv"
}

type csvOutput struct {
	Headers []string   `csv:"headers"`
	Entries [][]string `csv:"entries"`
}

func (f *CSVFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	table := csvOutput{
		Headers: []string{"ID", "Message", "Timestamp", "LogLevel", "AdditionalField"},
		Entries: [][]string{{e.GetTraceID(), e.GetMessage(), e.GetTimestamp().Format("2006-01-02T15:04:05Z07:00"), string(e.GetLevel()), fieldsToString(e.GetFields())}},
	}

	if f.Pretty {
		return marshalCSVPretty(table)
	}
	return marshalCSV(table)
}

func marshalCSVPretty(data csvOutput) ([]byte, error) {
	var b strings.Builder
	b.WriteString("CSV Output:\n")
	b.WriteString("Headers:\n")
	for _, header := range data.Headers {
		b.WriteString(" - " + header + "\n")
	}
	b.WriteString("Entries:\n")
	for _, entry := range data.Entries {
		b.WriteString(" - ")
		for i, value := range entry {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(value)
		}
		b.WriteString("\n")
	}
	return []byte(b.String()), nil
}

func marshalCSV(data csvOutput) ([]byte, error) {
	var b strings.Builder
	b.WriteString("ID,Message\n")
	for _, entry := range data.Entries {
		b.WriteString(entry[0] + "," + entry[1] + "\n")
	}
	return []byte(b.String()), nil
}

func fieldsToString(fields map[string]interface{}) string {
	var sb strings.Builder
	for k, v := range fields {
		sb.WriteString(k + "=" + toString(v) + " ")
	}
	return strings.TrimSpace(sb.String())
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int32, int64, float32, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return "unknown"
	}
}
