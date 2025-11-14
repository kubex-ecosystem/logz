package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	li "github.com/kubex-ecosystem/logz/internal/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// LogFormatter defines the contract for formatting log entries.
type LogFormatter interface {
	// Format converts a log entry to a formatted string.
	// Returns the formatted string and an error if formatting fails.
	Format(entry li.LogzEntry) (string, error)
}

// JSONFormatterImpl formats the log in JSON format.
type JSONFormatterImpl struct{}

// Format converts the log entry to JSON.
// Returns the JSON string and an error if marshalling fails.
func (f *JSONFormatterImpl) Format(entry li.LogzEntry) (string, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TextFormatterImpl formats the log in plain text.
type TextFormatterImpl struct{}

// Format converts the log entry to a formatted string with colors and icons.
// Returns the formatted string and an error if formatting fails.
func (f TextFormatterImpl) Format(entry li.LogzEntry) (string, error) {

	if !kbx.NoColor {
		switch kbx.LogLevel(entry.GetLevel()) {
		case kbx.NOTICE:
			kbx.LevelStr = "\033[33mNOTICE\033[0m"
		case kbx.TRACE:
			kbx.LevelStr = "\033[36mTRACE\033[0m"
		case kbx.SUCCESS:
			kbx.LevelStr = "\033[32mSUCCESS\033[0m"
		case kbx.DEBUG:
			kbx.LevelStr = "\033[34mDEBUG\033[0m"
		case kbx.INFO:
			kbx.LevelStr = "\033[32mINFO\033[0m"
		case kbx.WARN:
			kbx.LevelStr = "\033[33mWARN\033[0m"
		case kbx.ERROR:
			kbx.LevelStr = "\033[31mERROR\033[0m"
		case kbx.FATAL:
			kbx.LevelStr = "\033[35mFATAL\033[0m"
		case kbx.SILENT:
			kbx.LevelStr = ""
		case kbx.ANSWER:
			kbx.LevelStr = ""
		default:
			kbx.LevelStr = string(entry.GetLevel())
		}
	} else {
		kbx.LevelStr = string(entry.GetLevel())
	}

	systemLocale := os.Getenv("LANG")
	tag, _ := language.Parse(systemLocale)
	p := message.NewPrinter(tag)

	// Context and Metadata
	context := ""
	metadata := ""
	timestamp := ""
	if len(entry.GetMetadata()) > 0 {
		if sc, exist := entry.GetMetadata()["showContext"]; exist {
			tp := reflect.TypeOf(sc)
			if tp.Kind() == reflect.Bool {
				if sc.(bool) {
					if c, exists := entry.GetMetadata()["context"]; exists {
						context = c.(string)
					}
				}
			} else if tp.Kind() == reflect.String {
				if sc.(string) == "true" {
					metadata = fmt.Sprintf("\n%s", FormatMetadata(entry))
				}
			}

		} else if map[li.LogLevel]bool{li.DEBUG: true, li.INFO: true}[li.LogLevel(entry.GetLevel())] {
			if c, exists := entry.GetMetadata()["context"]; exists {
				context = c.(string)
			}
		}
		if smd, exist := entry.GetMetadata()["showData"]; exist {
			tp := reflect.TypeOf(smd)
			if tp.Kind() == reflect.Bool {
				if smd.(bool) {
					metadata = fmt.Sprintf("\n%s", FormatMetadata(entry))
				}
			} else if tp.Kind() == reflect.String {
				if smd.(string) == "true" {
					metadata = fmt.Sprintf("\n%s", FormatMetadata(entry))
				}
			}
		} else if entry.GetLevel() == string(li.DEBUG) {
			metadata = fmt.Sprintf("\n%s", FormatMetadata(entry))
		}
		if stp, exist := entry.GetMetadata()["showTimestamp"]; exist {
			tp := reflect.TypeOf(stp)
			if tp.Kind() == reflect.Bool {
				if stp.(bool) {
					timestamp = fmt.Sprintf("[%s]", entry.GetTimestamp().Format(p.Sprintf("%d-%m-%Y %H:%M:%S")))
				}
			} else if tp.Kind() == reflect.String {
				if stp.(string) == "true" {
					metadata = fmt.Sprintf("\n%s", FormatMetadata(entry))
				}
			}
		}
	}
	var header string
	if LevelStr != "" && icon != "" {
		// Construct the header
		header = fmt.Sprintf("%s [%s] %s %s - ", timestamp, LevelStr, context, icon)
	} else {
		header = strings.TrimSpace(fmt.Sprintf("%s %s", timestamp, context))
	}

	// Return the formatted log entry
	return fmt.Sprintf("%s%s%s", header, entry.GetMessage(), metadata), nil
}

type TableFormatterImpl struct{}

// Format formats the log entry as a table string.
// It uses the FormatRow method to get the rows and then constructs a string representation.
func (f *TableFormatterImpl) Format(entry li.LogzEntry) (string, error) {
	rows, err := f.FormatRow(entry)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, row := range rows {
		sb.WriteString(strings.Join(row, " | "))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// FormatRow formats the log entry as a table.
func (f *TableFormatterImpl) FormatRow(entry li.LogzEntry) ([][]string, error) {
	level := entry.GetLevel()
	context := entry.GetContext()
	message := entry.GetMessage()
	timestamp := entry.GetTimestamp()
	source := entry.GetSource()

	table := [][]string{
		{"Timestamp", "Level", "Context", "Source", "Message"},
		{timestamp.Format("2006-01-02 15:04:05"), string(level), context, source, message},
	}

	return table, nil
}

// FormatMetadata converts VMetadata to a JSON string.
// Returns the JSON string or an empty string if marshalling fails.
func FormatMetadata(entry li.LogzEntry) string {
	metadata := entry.GetMetadata()
	if len(metadata) == 0 {
		return ""
	}
	prefix := "Context:\n"
	for k, v := range metadata {
		if k == "showContext" {
			continue
		}
		prefix += fmt.Sprintf("  - %s: %v\n", k, v)
	}
	return prefix
}
