// Package cli provides the command-line interface for the Logz Logger.
package cli

import (
	"github.com/spf13/cobra"

	gl "github.com/kubex-ecosystem/logz/logger"

	"sync"
)

var (
	logCmdMap = map[string][]string{
		"debug":   []string{"dbg"},
		"notice":  []string{"not"},
		"success": []string{"suc"},
		"info":    []string{"inf"},
		"warn":    []string{"wrn"},
		"error":   []string{"err"},
		"answer":  []string{"ans"},
		"trace":   []string{"trc"},
		"fatal":   []string{"ftl"},
	}

	metaData, ctx map[string]string

	msg, output, format string

	logLevel, configFile string

	debugMode, quiet, hideBanner bool

	logger = gl.LoggerG.GetLogger()

	mu sync.RWMutex
)

func init() {
	gl.LoggerG.GetLogger()
}

// LogzCmds returns the CLI commands for different log levels and management.
func LogzCmds() *cobra.Command {
	short := "Logz Logger CLI - Your versatile logging tool"
	long := "Logz Logger CLI - Your versatile logging tool, supporting multiple log levels, formats, platforms, and log management features."

	cmd := &cobra.Command{
		Use:     "log",
		Short:   short,
		Long:    long,
		Aliases: []string{"log", "logs", "logger", "logging", "logz", "loggerz"},
		Annotations: GetDescriptions(
			[]string{
				long,
				short,
			},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	for level, aliases := range logCmdMap {
		cmd.AddCommand(newLogCmd(level, aliases))
	}

	cmd.Flags().StringVarP(&logLevel, "level", "l", "", "Log level")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file")
	cmd.Flags().StringVarP(&msg, "msg", "m", "", "Log message")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Output format")
	cmd.Flags().StringToStringVarP(&metaData, "metadata", "M", nil, "Metadata to include")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
	cmd.Flags().StringToStringVarP(&ctx, "context", "C", nil, "Context for the log")

	return cmd
}

// newLogCmd configures a command for a specific log level.
func newLogCmd(level string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     level,
		Aliases: aliases,
		Short:   "Log " + level + " level message",
		Long:    "Logs a " + level + " level message with optional metadata and context.",
		Example: "logz log " + level + " --msg 'This is a " + level + " message' --output 'logfile.log' --format 'json' --metadata key1=value1,key2=value2 --context user=admin,session=xyz",
		Annotations: GetDescriptions(
			[]string{"Logs a " + level + " level message"},
			hideBanner,
		),
		Run: func(cmd *cobra.Command, args []string) {
			if logLevel != "" {
				logger.SetLevel(logLevel)
			}
			if debugMode {
				logger.SetDebug(true)
			}
			if quiet {
				logger.SetLogLevel("error")
			}
			if format != "" {
				// lFmt := gl.LogFormat(format)
				logger.SetFormat(format)
			}

			logger.Log(level, msg, output, metaData, ctx)
		},
	}

	cmd.Flags().StringVarP(&logLevel, "level", "l", "", "Log level")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file")
	cmd.Flags().StringVarP(&msg, "msg", "m", "", "Log message")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Output format")
	cmd.Flags().StringToStringVarP(&metaData, "metadata", "M", nil, "Metadata to include")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
	cmd.Flags().StringToStringVarP(&ctx, "context", "C", nil, "Context for the log")

	return cmd
}

// // rotateLogsCmd allows manual log rotation.
// func rotateLogsCmd() *cobra.Command {
// 	var mu sync.RWMutex

// 	return &cobra.Command{
// 		Use: "rotate",
// 		Annotations: GetDescriptions(
// 			[]string{"Rotates logs that exceed the configured size"},
// 			false,
// 		),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			mu.Lock()
// 			defer mu.Unlock()

// 			configManager := il.NewConfigManager()
// 			if configManager == nil {
// 				fmt.Println("ErrorCtx initializing ConfigManager.")
// 				return
// 			}
// 			cfgMgr := *configManager

// 			config, err := cfgMgr.LoadConfig()
// 			if err != nil {
// 				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
// 				return
// 			}

// 			err = il.CheckLogSize(config)
// 			if err != nil {
// 				fmt.Printf("ErrorCtx rotating logs: %v\n", err)
// 			} else {
// 				fmt.Println("Logs rotated successfully!")
// 			}
// 		},
// 	}
// }

// // checkLogSizeCmd checks the current log size.
// func checkLogSizeCmd() *cobra.Command {
// 	var mu sync.RWMutex

// 	return &cobra.Command{
// 		Use: "check-size",
// 		Annotations: GetDescriptions(
// 			[]string{"Checks the log size without taking any action"},
// 			false,
// 		),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			mu.Lock()
// 			defer mu.Unlock()

// 			configManager := il.NewConfigManager()
// 			if configManager == nil {
// 				fmt.Println("ErrorCtx initializing ConfigManager.")
// 				return
// 			}
// 			cfgMgr := *configManager

// 			config, err := cfgMgr.LoadConfig()
// 			if err != nil {
// 				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
// 				return
// 			}

// 			logDir := config.Output()
// 			logSize, err := il.GetLogDirectorySize(filepath.Dir(logDir)) // Add this function to core
// 			if err != nil {
// 				fmt.Printf("ErrorCtx calculating log size: %v\n", err)
// 				return
// 			}

// 			sizeInMB := logSize / (1024 * 1024)

// 			fmt.Printf("The total log size in directory '%s' is: %d MB\n", filepath.Dir(logDir), sizeInMB)
// 		},
// 	}
// }

// // archiveLogsCmd allows manual log archiving.
// func archiveLogsCmd() *cobra.Command {
// 	var mu sync.RWMutex

// 	return &cobra.Command{
// 		Use: "archive",
// 		Annotations: GetDescriptions(
// 			[]string{"Manually archives all logs"},
// 			false,
// 		),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			mu.Lock()
// 			defer mu.Unlock()

// 			err := il.ArchiveLogs(nil)
// 			if err != nil {
// 				fmt.Printf("ErrorCtx archiving logs: %v\n", err)
// 			} else {
// 				fmt.Println("Logs archived successfully!")
// 			}
// 		},
// 	}
// }

// // watchLogsCmd monitors logs in real-time.
// func watchLogsCmd() *cobra.Command {
// 	var mu sync.RWMutex

// 	return &cobra.Command{
// 		Use:     "watch",
// 		Aliases: []string{"w"},
// 		Annotations: GetDescriptions(
// 			[]string{"Monitors logs in real-time"},
// 			false,
// 		),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			mu.Lock()
// 			defer mu.Unlock()

// 			configManager := il.NewConfigManager()
// 			if configManager == nil {
// 				fmt.Println("ErrorCtx initializing ConfigManager.")
// 				return
// 			}
// 			cfgMgr := *configManager

// 			config, err := cfgMgr.LoadConfig()
// 			if err != nil {
// 				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
// 				return
// 			}

// 			logFilePath := config.Output()
// 			reader := il.NewFileLogReader()
// 			stopChan := make(chan struct{})

// 			// Capture signals for interruption
// 			sigChan := make(chan os.Signal, 1)
// 			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
// 			go func() {
// 				<-sigChan
// 				close(stopChan)
// 			}()

// 			fmt.Println("Monitoring started (Ctrl+C to exit):")
// 			if err := reader.Tail(logFilePath, stopChan); err != nil {
// 				fmt.Printf("ErrorCtx monitoring logs: %v\n", err)
// 			}

// 			// Wait a small delay to finish
// 			time.Sleep(500 * time.Millisecond)
// 		},
// 	}
// }
