package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/logger"
	"github.com/spf13/cobra"
)

// ServiceCmd creates the main command for managing the web service.
func ServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "service",
		Annotations: GetDescriptions(
			[]string{"Start, stop, and get information about the web service"},
			false,
		),
	}
	cmd.AddCommand(startServiceCmd())
	cmd.AddCommand(stopServiceCmd())
	cmd.AddCommand(getServiceCmd())
	cmd.AddCommand(spawnServiceCmd())
	return cmd
}

// startServiceCmd creates the command to start the web service.
func startServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the web service",
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Error initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			_, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				return
			}

			if err := logger.Start("9999"); err != nil {
				fmt.Printf("Error starting service: %v\n", err)
			} else {
				fmt.Println("Service started successfully.")
			}
		},
	}
}

// stopServiceCmd creates the command to stop the web service.
func stopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the web service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := logger.Stop(); err != nil {
				fmt.Printf("Error stopping service: %v\n", err)
			} else {
				fmt.Println("Service stopped successfully.")
			}
		},
	}
}

// getServiceCmd creates the command to get information about the running web service.
func getServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get information about the running service",
		Run: func(cmd *cobra.Command, args []string) {
			pid, port, pidPath, err := logger.GetServiceInfo()
			if err != nil {
				fmt.Println("Service is not running")
			} else {
				fmt.Printf("Service running with PID %d on port %s\n", pid, port)
				fmt.Printf("PID file: %s\n", pidPath)
			}
		},
	}
}

// spawnServiceCmd creates the command to spawn a new instance of the web service.
func spawnServiceCmd() *cobra.Command {
	var configPath string
	spCmd := &cobra.Command{
		Use:    "spawn",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	spCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the service configuration file")
	return spCmd
}
