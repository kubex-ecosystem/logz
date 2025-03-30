package logger

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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

	// Configure colors and icons by level
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

// LogWriter defines the contract for writing logs.
type LogWriter interface {
	// Write writes a formatted log entry.
	// Returns an error if writing fails.
	Write(entry LogzEntry) error
}

// DefaultWriter implements LogWriter using an io.Writer and a LogFormatter.
type DefaultWriter struct {
	out       io.Writer
	formatter LogFormatter
}

// NewDefaultWriter creates a new instance of DefaultWriter.
// Takes an io.Writer and a LogFormatter as parameters.
func NewDefaultWriter(out io.Writer, formatter LogFormatter) *DefaultWriter {
	return &DefaultWriter{
		out:       out,
		formatter: formatter,
	}
}

// Write formats the entry and writes it to the configured destination.
// Returns an error if formatting or writing fails.
func (w *DefaultWriter) Write(entry LogzEntry) error {
	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w.out, formatted)
	return err
}

func Writer(module string) io.Writer {
	currentPid := os.Getpid()
	logFileName := strings.Join([]string{module, "logz_", strconv.Itoa(currentPid), ".log"}, "")
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil || cacheDir == "" {
		cacheDir = os.TempDir()
	}
	logFilePath := filepath.Join(cacheDir, logFileName)
	if logFileStatErr := os.Remove(logFilePath); logFileStatErr == nil {
		cmdRm := fmt.Sprintf("rm -f %s", logFilePath)
		if _, cmdRmErr := exec.Command("bash", "-c", cmdRm).Output(); cmdRmErr != nil {
			fmt.Println(cmdRmErr)
			return os.Stdout
		}
	}
	if logFilePathErr := os.MkdirAll(filepath.Dir(logFilePath), 0777); logFilePathErr != nil {
		fmt.Println(logFilePathErr)
		return os.Stdout
	}
	logFile, logFileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if logFileErr != nil {
		fmt.Println(logFileErr)
		return os.Stdout
	}
	return logFile
}

// formatMetadata converts metadata to a JSON string.
// Returns the JSON string or an empty string if marshalling fails.
func formatMetadata(entry LogzEntry) string {
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
