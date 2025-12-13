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

	gl "github.com/kubex-ecosystem/logz"
)

var initArgs = kbx.NewConfig()

func LogzCmd() *cobra.Command {
	cmd := LoggerCmd()

	return cmd
}

func LoggerCmd() *cobra.Command {
	var Output, Format, Level, MinLevel, MaxLevel string
	var DisableColors, ShowTraceID, ShowFields, ShowStack, DisableIcons bool

	short := "Logger related operations"
	long := `Perform various logger related operations such as setting log levels, formats, and outputs.
You can configure the logger to suit your application's needs.`

	loggerCmd := &cobra.Command{
		Use:     "logger",
		Aliases: []string{"log", "lg", "lz", "logz", "l"},
		Short:   short,
		Long:    long,
		Annotations: info.GetDescriptions(
			[]string{short, long}, os.Getenv("LOGZ_HIDE_BANNER") == "true",
		),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Checa se o primeiro argumento é um nível de log válido
				if kbx.IsLevel(args[0]) {
					kbx.LoggerArgs.Level = kbx.ParseLevel(args[0])
					args = args[1:] // Remove o nível do log dos argumentos
				}
				// Junta os argumentos restantes como mensagem
				if len(args) > 0 {
					kbx.LoggerArgs.Messages = append(kbx.LoggerArgs.Messages, args...)
				}
			}

			kbx.LoggerArgs.Level = kbx.ParseLevel(Level)

			// Configurar argumentos do logger com valores padrão se não especificados

			kbx.LoggerArgs.Format = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.Format, kbx.GetValueOrDefaultSimple(Format, "text"))
			kbx.LoggerArgs.Output = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.Output, gl.ParseWriter(Output))
			kbx.LoggerArgs.Level = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.Level, gl.ParseLevel(Level))
			kbx.LoggerArgs.MinLevel = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.MinLevel, gl.ParseLevel(MinLevel))
			kbx.LoggerArgs.MaxLevel = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.MaxLevel, gl.ParseLevel(MaxLevel))
			kbx.LoggerArgs.ShowColor = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.ShowColor, kbx.BoolPtr(!DisableColors))
			kbx.LoggerArgs.ShowIcons = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.ShowIcons, kbx.BoolPtr(!DisableIcons))
			kbx.LoggerArgs.ShowTraceID = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.ShowTraceID, ShowTraceID)
			kbx.LoggerArgs.ShowFields = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.ShowFields, ShowFields)
			kbx.LoggerArgs.ShowStack = kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.ShowStack, ShowStack)

			// Criar opções do logger
			opts := core.NewLoggerOptions(kbx.LoggerArgs)

			// Aplicar opções avançadas se especificadas
			opts.MinLevel = kbx.LoggerArgs.MinLevel
			opts.MaxLevel = kbx.LoggerArgs.MaxLevel
			opts.Output = kbx.LoggerArgs.Output
			opts.ShowTraceID = kbx.LoggerArgs.ShowTraceID
			opts.ShowFields = kbx.LoggerArgs.ShowFields
			opts.ShowIcons = kbx.LoggerArgs.ShowIcons
			opts.ShowColor = kbx.LoggerArgs.ShowColor
			opts.ShowStack = kbx.LoggerArgs.ShowStack
			opts.StackTrace = kbx.BoolPtr(kbx.LoggerArgs.ShowStack)

			// Aplicar metadata se especificado

			if len(kbx.LoggerArgs.Metadata) > 0 {
				opts.LoggerConfig.Metadata = kbx.LoggerArgs.Metadata
				fields := make(map[string]any)
				for k, v := range kbx.LoggerArgs.Metadata {
					fields[k] = v
				}
				opts.LogzAdvancedOptions.Metadata = fields
			}

			// Aplicar formato se especificado
			if kbx.LoggerArgs.Format != "" {
				opts.Formatter = formatter.ParseFormatter(kbx.LoggerArgs.Format, kbx.DefaultTrue(kbx.LoggerArgs.ShowColor))
			} else {
				// Usar formatter padrão se não especificado
				opts.Formatter = formatter.NewTextFormatter(kbx.DefaultTrue(kbx.LoggerArgs.ShowColor))
			}

			// Criar logger
			logger := core.NewLogger(kbx.GetValueOrDefaultSimple(kbx.LoggerArgs.Prefix, "LogzCLI"), opts, false)
			if logger == nil {
				return fmt.Errorf("failed to create logger")
			}

			// Criar entrada de log
			entry := core.NewLogzEntry(kbx.LoggerArgs.Level).
				WithMessage(strings.TrimSpace(strings.ToValidUTF8(strings.Join(kbx.LoggerArgs.Messages, " "), ""))).
				WithColor(kbx.DefaultTrue(kbx.LoggerArgs.ShowColor)).
				WithIcon(kbx.DefaultTrue(kbx.LoggerArgs.ShowIcons)).
				WithFields(opts.LogzAdvancedOptions.Metadata).
				WithTraceID(kbx.LoggerArgs.ID.String()).
				WithShowTraceID(kbx.LoggerArgs.ShowTraceID).
				WithShowCaller(kbx.LoggerArgs.ShowStack).
				WithShowFields(kbx.LoggerArgs.ShowFields).
				WithStack(kbx.LoggerArgs.ShowStack).
				WithCaller("CLI")

			// Usar o level correto ao invés de "error" fixo
			return logger.Log(kbx.LoggerArgs.Level, entry)
		},
	}

	loggerCmd.Flags().BoolVarP(&kbx.LoggerArgs.Debug, "debug", "D", false, "Enable debug mode")
	loggerCmd.Flags().StringVarP(&Level, "level", "l", "info", "Set the logging level (e.g., debug, info, warn, error)")
	loggerCmd.Flags().StringVarP(&MinLevel, "min-level", "L", "debug", "Set the minimum logging level")
	loggerCmd.Flags().StringVarP(&MaxLevel, "max-level", "U", "fatal", "Set the maximum logging level")
	loggerCmd.Flags().StringVarP(&Output, "output", "o", "stdout", "Set the logging output (e.g., stdout, file)")
	loggerCmd.Flags().StringVarP(&Format, "format", "f", "text", "Set the logging format (e.g., json, text)")
	loggerCmd.Flags().StringArrayVarP(&kbx.LoggerArgs.Messages, "message", "m", []string{}, "Log message parts")
	loggerCmd.Flags().StringToStringVarP(&kbx.LoggerArgs.Metadata, "metadata", "M", map[string]string{}, "Set metadata key-value pairs for the log entry")
	loggerCmd.Flags().BoolVarP(&DisableColors, "disableColors", "c", false, "Enable colored output")
	loggerCmd.Flags().BoolVarP(&DisableIcons, "disableIcons", "i", false, "Enable icons in the log entry")
	loggerCmd.Flags().BoolVarP(&ShowTraceID, "showTraceID", "t", false, "Include trace ID in the log entry")
	loggerCmd.Flags().BoolVarP(&initArgs.ShowStack, "showStack", "S", false, "Include stack trace in the log entry")
	loggerCmd.Flags().BoolVarP(&kbx.LoggerArgs.ShowFields, "showFields", "F", false, "Include fields in the log entry")
	loggerCmd.Flags().StringVarP(&kbx.LoggerArgs.Prefix, "prefix", "p", "LogzCLI", "Set the log message prefix")

	loggerCmd.MarkFlagFilename("output")

	return loggerCmd
}
