package core

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// LogFormatter defines the contract for formatting log entries.
type LogFormatter interface {
	// Format converts a log entry to a formatted string.
	// Returns the formatted string and an error if formatting fails.
	Format(entry LogzEntry) (string, error)
}

// JSONFormatter formats the log in JSON format.
type JSONFormatter struct{}

// Format converts the log entry to JSON.
// Returns the JSON string and an error if marshalling fails.
func (f *JSONFormatter) Format(entry LogzEntry) (string, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TextFormatter formats the log in plain text.
type TextFormatter struct{}

// Format converts the log entry to a formatted string with colors and icons.
// Returns the formatted string and an error if formatting fails.
func (f *TextFormatter) Format(entry LogzEntry) (string, error) {

	// Check for environment variables
	noColor := os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows"
	noIcon := os.Getenv("LOGZ_NO_ICON") != ""

	icon, levelStr := "", ""

	if !noIcon {
		switch entry.GetLevel() {
		case NOTICE:
			icon = "\033[33m📝\033[0m "
		case TRACE:
			icon = "\033[36m🔍\033[0m "
		case SUCCESS:
			icon = "\033[32m✅\033[0m "
		case DEBUG:
			icon = "\033[34m🐛\033[0m "
		case INFO:
			icon = "\033[32mℹ️\033[0m "
		case WARN:
			icon = "\033[33m⚠️\033[0m "
		case ERROR:
			icon = "\033[31m❌\033[0m "
		case FATAL:
			icon = "\033[35m💀\033[0m "
		default:
			icon = ""
		}
	} else {
		icon = ""
	}

	// Configure colors and icons by VLevel
	if !noColor {
		switch entry.GetLevel() {
		case NOTICE:
			levelStr = "\033[33mNOTICE\033[0m"
		case TRACE:
			levelStr = "\033[36mTRACE\033[0m"
		case SUCCESS:
			levelStr = "\033[32mSUCCESS\033[0m"
		case DEBUG:
			levelStr = "\033[34mDEBUG\033[0m"
		case INFO:
			levelStr = "\033[32mINFO\033[0m"
		case WARN:
			levelStr = "\033[33mWARN\033[0m"
		case ERROR:
			levelStr = "\033[31mERROR\033[0m"
		case FATAL:
			levelStr = "\033[35mFATAL\033[0m"
		default:
			levelStr = string(entry.GetLevel())
		}
	} else {
		levelStr = string(entry.GetLevel())
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
					metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
				}
			}

		} else if map[LogLevel]bool{DEBUG: true, INFO: true}[entry.GetLevel()] {
			if c, exists := entry.GetMetadata()["context"]; exists {
				context = c.(string)
			}
		}
		if smd, exist := entry.GetMetadata()["showData"]; exist {
			tp := reflect.TypeOf(smd)
			if tp.Kind() == reflect.Bool {
				if smd.(bool) {
					metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
				}
			} else if tp.Kind() == reflect.String {
				if smd.(string) == "true" {
					metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
				}
			}
		} else if entry.GetLevel() == DEBUG {
			metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
		}
		if stp, exist := entry.GetMetadata()["showTimestamp"]; exist {
			tp := reflect.TypeOf(stp)
			if tp.Kind() == reflect.Bool {
				if stp.(bool) {
					timestamp = fmt.Sprintf("[%s]", entry.GetTimestamp().Format(p.Sprintf("%d-%m-%Y %H:%M:%S")))
				}
			} else if tp.Kind() == reflect.String {
				if stp.(string) == "true" {
					metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
				}
			}
		}
	}

	// Construct the header
	header := fmt.Sprintf("%s [%s] %s %s - ", timestamp, levelStr, context, icon)

	// Return the formatted log entry
	return fmt.Sprintf("%s%s%s", header, entry.GetMessage(), metadata), nil
}

type TableFormatter struct{}

// Format formats the log entry as a table string.
// It uses the FormatRow method to get the rows and then constructs a string representation.
func (f *TableFormatter) Format(entry LogzEntry) (string, error) {
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
func (f *TableFormatter) FormatRow(entry LogzEntry) ([][]string, error) {
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
