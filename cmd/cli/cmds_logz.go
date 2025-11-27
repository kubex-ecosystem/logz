package cli

import (
	"io"
	"os"

	"github.com/kubex-ecosystem/logz/interfaces"
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
		Run: func(cmd *cobra.Command, args []string) {
			opts := core.NewLoggerOptions()
			opts.Debug = kbx.BoolPtr(initArgs.Debug)
			opts.Level = interfaces.ParseLevel(initArgs.Level)
			opts.MinLevel = interfaces.ParseLevel(initArgs.MinLevel)
			opts.MaxLevel = interfaces.ParseLevel(initArgs.MaxLevel)
			opts.Output = kbx.GetValueOrDefaultSimple[io.Writer](any(initArgs.Output).(io.Writer), os.Stdout)
			opts.Formatter = interfaces.ParseFormatter(initArgs.Format)

			logger, entry := SetupLogger(
				opts,
				initArgs.Message,
			)

			logger.Log(initArgs.Level, entry)

		},
	}

	loggerCmd.Flags().BoolVarP(&initArgs.Debug, "debug", "d", false, "Enable debug mode")
	loggerCmd.Flags().StringVarP(&initArgs.Level, "level", "l", "info", "Set the logging level (e.g., debug, info, warn, error)")
	loggerCmd.Flags().StringVarP(&initArgs.MinLevel, "min-level", "m", "", "Set the minimum logging level")
	loggerCmd.Flags().StringVarP(&initArgs.MaxLevel, "max-level", "M", "", "Set the maximum logging level")
	loggerCmd.Flags().StringVarP(&initArgs.Output, "output", "o", "stdout", "Set the logging output (e.g., stdout, file)")
	loggerCmd.Flags().StringVarP(&initArgs.Format, "format", "f", "json", "Set the logging format (e.g., json, text)")
	loggerCmd.Flags().StringArrayVarP(&initArgs.Message, "message", "e", []string{}, "Log message parts")

	return loggerCmd
}

func SetupLogger(loggerOptions *core.LoggerOptionsImpl, messages []string) (*core.Logger, interfaces.Entry) {
	logger := core.NewLogger("LogzCLI", loggerOptions, false)
	entry := core.NewEntry()
	for _, msgPart := range messages {
		entry.WithMessage(entry.GetMessage() + msgPart + " ")
	}
	entry.WithMessage(entry.GetMessage()[:len(entry.GetMessage())-1]) // Remove trailing space
	return logger, entry
}
