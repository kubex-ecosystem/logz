// Package cli implements the command-line interface for Logz.
package cli

import (
	"os"
	"strings"
	"time"

	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/module/info"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/spf13/cobra"
)

var initArgs = &kbx.InitArgs{}

func LogzCmd() *cobra.Command {
	logzCmd := &cobra.Command{
		Use:   "logz",
		Short: "Logz related commands",
		Long:  `Commands related to the Logz logging library.`,
		Annotations: info.GetDescriptions(
			[]string{
				"Logz related commands",
				`Commands related to the Logz logging library.`,
			},
			os.Getenv("LOGZ_HIDE_BANNER") == "true",
		),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	logzCmd.AddCommand(LoggerCmd())

	return logzCmd
}

func LoggerCmd() *cobra.Command {
	var Output, Format, Level, MinLevel, MaxLevel string
	var Message []string

	short := "Logger related operations"
	long := `Perform various logger related operations such as setting log levels, formats, and outputs.
You can configure the logger to suit your application's needs.`

	loggerCmd := &cobra.Command{
		Use:     "logger",
		Aliases: []string{"log", "lg"},
		Short:   short,
		Long:    long,
		Annotations: info.GetDescriptions(
			[]string{short, long}, os.Getenv("LOGZ_HIDE_BANNER") == "true",
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			opts := core.NewLoggerOptions(initArgs)

			entry, err := core.NewEntry()

			logger := core.NewLogger("LogzCLI", opts, false)
			if err != nil {
				return err
			}
			entry.Message = strings.TrimSpace(
				strings.ToValidUTF8(
					strings.Join(
						Message,
						" ",
					), "",
				),
			)
			entry.Level = opts.Level
			entry.Timestamp = time.Now().UTC()
			entry.Fields = make(map[string]any)
			// Add any additional fields to the entry
			entry.Fields["module"] = "logger_cmd"
			entry.Fields["timestamp"] = entry.Timestamp.Format(time.RFC3339)
			// Add any additional fields to the entry
			entry.Fields["caller"] = "logger_cmd"
			entry.Fields["file"] = "cmd/cli/cmds_logz.go"
			// Add any additional fields to the entry
			entry.Fields["args"] = args
			entry.Fields["message"] = entry.Message

			return logger.Log("error", entry)

		},
	}

	loggerCmd.Flags().BoolVarP(&initArgs.Debug, "debug", "d", false, "Enable debug mode")
	loggerCmd.Flags().StringVarP(&Level, "level", "l", "info", "Set the logging level (e.g., debug, info, warn, error)")
	loggerCmd.Flags().StringVarP(&MinLevel, "min-level", "m", "debug", "Set the minimum logging level")
	loggerCmd.Flags().StringVarP(&MaxLevel, "max-level", "M", "fatal", "Set the maximum logging level")
	loggerCmd.Flags().StringVarP(&Output, "output", "o", "stdout", "Set the logging output (e.g., stdout, file)")
	loggerCmd.Flags().StringVarP(&Format, "format", "f", "text", "Set the logging format (e.g., json, text)")
	loggerCmd.Flags().StringArrayVarP(&Message, "message", "e", []string{}, "Log message parts")

	return loggerCmd
}
