// Package cli implements the command-line interface for Logz.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubex-ecosystem/logz/internal/core"
	"github.com/kubex-ecosystem/logz/internal/formatter"
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
			kbx.ParseLoggerArgs(Level, MinLevel, MaxLevel, Output)

			// Criar opções do logger
			opts := core.NewLoggerOptions(kbx.LoggerArgs)

			// Aplicar formato se especificado
			if Format != "" {
				opts.Formatter = formatter.ParseFormatter(Format, false)
			} else {
				// Usar formatter padrão se não especificado
				opts.Formatter = formatter.NewTextFormatter(false)
			}

			logger := core.NewLogger("LogzCLI", opts, false)
			if logger == nil {
				return fmt.Errorf("failed to create logger")
			}

			entry, err := core.NewEntry(kbx.LoggerArgs.Level)
			if err != nil {
				return err
			}
			entry.Message = strings.TrimSpace(strings.ToValidUTF8(strings.Join(Message, " "), ""))

			// Usar o level correto ao invés de "error" fixo
			return logger.Log(kbx.LoggerArgs.Level, entry)

		},
	}

	loggerCmd.Flags().BoolVarP(&kbx.LoggerArgs.Debug, "debug", "d", false, "Enable debug mode")
	loggerCmd.Flags().StringVarP(&Level, "level", "l", "info", "Set the logging level (e.g., debug, info, warn, error)")
	loggerCmd.Flags().StringVarP(&MinLevel, "min-level", "m", "debug", "Set the minimum logging level")
	loggerCmd.Flags().StringVarP(&MaxLevel, "max-level", "M", "fatal", "Set the maximum logging level")
	loggerCmd.Flags().StringVarP(&Output, "output", "o", "stdout", "Set the logging output (e.g., stdout, file)")
	loggerCmd.Flags().StringVarP(&Format, "format", "f", "text", "Set the logging format (e.g., json, text)")
	loggerCmd.Flags().StringArrayVarP(&Message, "message", "e", []string{}, "Log message parts")

	return loggerCmd
}
