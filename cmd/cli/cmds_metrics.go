package cli

import (
	il "github.com/faelmori/logz/internal/core"

	"github.com/spf13/cobra"

	"fmt"
	"strconv"
	"sync"
	"time"
)

// MetricsCmd creates the main command for managing Prometheus metrics.
func MetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "metrics",
		Annotations: GetDescriptions(
			[]string{"Manage Prometheus metrics"},
			false,
		),
	}

	cmd.AddCommand(enableMetricsCmd())
	cmd.AddCommand(disableMetricsCmd())
	cmd.AddCommand(addMetricCmd())
	cmd.AddCommand(removeMetricCmd())
	cmd.AddCommand(listMetricsCmd())
	cmd.AddCommand(watchMetricsCmd())

	return cmd
}

// enableMetricsCmd creates the command to enable Prometheus integration.
func enableMetricsCmd() *cobra.Command {
	var port string
	var mu sync.RWMutex

	enMCmd := &cobra.Command{
		Use:     "enable",
		Aliases: []string{"en"},
		Short:   "Enable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			pm := il.GetPrometheusManager()
			pm.Enable(port)
		},
	}
	enMCmd.Flags().StringVarP(&port, "port", "p", "2112", "Port to expose Prometheus metrics")
	return enMCmd
}

// disableMetricsCmd creates the command to disable Prometheus integration.
func disableMetricsCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use:     "disable",
		Aliases: []string{"dis"},
		Short:   "Disable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			pm := il.GetPrometheusManager()
			pm.Disable()
		},
	}
}

// addMetricCmd creates the command to add or update a Prometheus metric.
func addMetricCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use:     "add [name] [value]",
		Aliases: []string{"a"},
		Short:   "Add or update a Prometheus metric",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			name := args[0]
			value, valueErr := strconv.ParseFloat(args[1], 64)
			if valueErr != nil {
				fmt.Printf("Invalid metric value: %v\n", valueErr)
				return
			}
			pm := il.GetPrometheusManager()
			pm.AddMetric(name, value, nil)
		},
	}
}

// removeMetricCmd creates the command to remove a Prometheus metric.
func removeMetricCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use:     "remove [name]",
		Aliases: []string{"r"},
		Short:   "Remove a Prometheus metric",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mu.Lock()
			defer mu.Unlock()

			name := args[0]
			pm := il.GetPrometheusManager()
			pm.RemoveMetric(name)
		},
	}
}

// listMetricsCmd creates the command to list all Prometheus metrics.
func listMetricsCmd() *cobra.Command {
	var mu sync.RWMutex

	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all Prometheus metrics",
		Run: func(cmd *cobra.Command, args []string) {
			mu.RLock()
			defer mu.RUnlock()

			pm := il.GetPrometheusManager()
			metrics := pm.GetMetrics()
			if len(metrics) == 0 {
				fmt.Println("No metrics registered.")
				return
			}
			fmt.Println("Registered metrics:")
			for name, value := range metrics {
				fmt.Printf(" - %s: %f\n", name, value)
			}
		},
	}
}

// watchMetricsCmd creates the command to watch Prometheus metrics in real time.
func watchMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Watch Prometheus metrics in real time",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Watching metrics (press Ctrl+C to exit):")
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					metrics := il.GetPrometheusManager().GetMetrics()
					fmt.Println("Current Metrics:")
					if len(metrics) == 0 {
						fmt.Println("  No metrics registered.")
					} else {
						for name, value := range metrics {
							fmt.Printf(" - %s: %f\n", name, value)
						}
					}
					fmt.Println("-----")
				}
			}
		},
	}
}
