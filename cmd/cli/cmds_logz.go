package cli

import (
	il "github.com/faelmori/logz/internal/core"

	"github.com/spf13/cobra"

	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// LogzCmds returns the CLI commands for different log levels and management.
func LogzCmds() []*cobra.Command {
	return []*cobra.Command{
		newLogCmd("debug", []string{"dbg"}),
		newLogCmd("notice", []string{"not"}),
		newLogCmd("success", []string{"suc"}),
		newLogCmd("info", []string{"inf"}),
		newLogCmd("warn", []string{"wrn"}),
		newLogCmd("error", []string{"err"}),
		newLogCmd("fatal", []string{"ftl"}),
		watchLogsCmd(),
		startServiceCmd(),
		stopServiceCmd(),
		rotateLogsCmd(),
		checkLogSizeCmd(),
		archiveLogsCmd(),
	}
}

// newLogCmd configures a command for a specific log level.
func newLogCmd(level string, aliases []string) *cobra.Command {
	var metaData, ctx map[string]string
	var msg, output, format string
	var mu sync.RWMutex

	cmd := &cobra.Command{
		Use:     level,
		Aliases: aliases,
		Annotations: GetDescriptions(
			[]string{"Logs a " + level + " level message"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			configManager := il.NewConfigManager()
			if configManager == nil {
				fmt.Println("ErrorCtx initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
				return
			}

			if format != "" {
				config.SetFormat(il.LogFormat(format))
			}

			if output != "" {
				config.SetOutput(output)
			}
			logr := il.NewLogger("logz")
			for k, v := range metaData {
				logr.SetMetadata(k, v)
			}
			ctxInterface := make(map[string]interface{})
			for k, v := range ctx {
				ctxInterface[k] = v
				mu.Lock()
				defer mu.Unlock()
			}
			switch level {
			case "debug":
				logr.DebugCtx(msg, ctxInterface)
			case "notice":
				logr.NoticeCtx(msg, ctxInterface)
			case "info":
				logr.InfoCtx(msg, ctxInterface)
			case "success":
				logr.SuccessCtx(msg, ctxInterface)
			case "warn":
				logr.WarnCtx(msg, ctxInterface)
			case "error":
				logr.ErrorCtx(msg, ctxInterface)
			case "fatal":
				logr.FatalCtx(msg, ctxInterface)
			default:
				logr.InfoCtx(msg, ctxInterface)
			}
		},
	}

	cmd.Flags().StringVarP(&msg, "msg", "M", "", "Log message")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Output format")
	cmd.Flags().StringToStringVarP(&metaData, "metadata", "m", nil, "Metadata to include")
	cmd.Flags().StringToStringVarP(&ctx, "context", "c", nil, "Context for the log")

	return cmd
}

// rotateLogsCmd allows manual log rotation.
func rotateLogsCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use: "rotate",
		Annotations: GetDescriptions(
			[]string{"Rotates logs that exceed the configured size"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			configManager := il.NewConfigManager()
			if configManager == nil {
				fmt.Println("ErrorCtx initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
				return
			}

			err = il.CheckLogSize(config)
			if err != nil {
				fmt.Printf("ErrorCtx rotating logs: %v\n", err)
			} else {
				fmt.Println("Logs rotated successfully!")
			}
		},
	}
}

// checkLogSizeCmd checks the current log size.
func checkLogSizeCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use: "check-size",
		Annotations: GetDescriptions(
			[]string{"Checks the log size without taking any action"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			configManager := il.NewConfigManager()
			if configManager == nil {
				fmt.Println("ErrorCtx initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
				return
			}

			logDir := config.Output()
			logSize, err := il.GetLogDirectorySize(filepath.Dir(logDir)) // Add this function to core
			if err != nil {
				fmt.Printf("ErrorCtx calculating log size: %v\n", err)
				return
			}

			sizeInMB := logSize / (1024 * 1024)

			fmt.Printf("The total log size in directory '%s' is: %d MB\n", filepath.Dir(logDir), sizeInMB)
		},
	}
}

// archiveLogsCmd allows manual log archiving.
func archiveLogsCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use: "archive",
		Annotations: GetDescriptions(
			[]string{"Manually archives all logs"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			err := il.ArchiveLogs(nil)
			if err != nil {
				fmt.Printf("ErrorCtx archiving logs: %v\n", err)
			} else {
				fmt.Println("Logs archived successfully!")
			}
		},
	}
}

// watchLogsCmd monitors logs in real-time.
func watchLogsCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Annotations: GetDescriptions(
			[]string{"Monitors logs in real-time"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			configManager := il.NewConfigManager()
			if configManager == nil {
				fmt.Println("ErrorCtx initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("ErrorCtx loading configuration: %v\n", err)
				return
			}

			logFilePath := config.Output()
			reader := il.NewFileLogReader()
			stopChan := make(chan struct{})

			// Capture signals for interruption
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				close(stopChan)
			}()

			fmt.Println("Monitoring started (Ctrl+C to exit):")
			if err := reader.Tail(logFilePath, stopChan); err != nil {
				fmt.Printf("ErrorCtx monitoring logs: %v\n", err)
			}

			// Wait a small delay to finish
			time.Sleep(500 * time.Millisecond)
		},
	}
}
