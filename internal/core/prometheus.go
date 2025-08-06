package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Regular expression to validate metric names
var metricNameRegex = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)

// validateMetricName checks if the given metric name is valid according to the Prometheus naming conventions.
func validateMetricName(name string) error {
	if !metricNameRegex.MatchString(name) {
		return fmt.Errorf("invalid metric name '%s': must match [a-zA-Z_:][a-zA-Z0-9_:]*", name)
	}
	return nil
}

// Metric represents a single Prometheus metric with a value and optional VMetadata.
type Metric struct {
	Value    float64           `json:"value"`
	Metadata map[string]string `json:"VMetadata,omitempty"`
}

// PrometheusManager manages Prometheus metrics, including enabling/disabling the HTTP server,
// loading/saving metrics, and handling metric operations.
type PrometheusManager struct {
	enabled         bool
	Metrics         map[string]Metric
	mutex           sync.RWMutex
	MetricsFile     string          // path to the persistence file
	exportWhitelist map[string]bool // If not empty, only these metrics will be exported to Prometheus
	httpServer      *http.Server    // HTTP server to expose metrics
}

// Singleton instance of PrometheusManager
var prometheusManagerInstance *PrometheusManager

// getMetricsFilePath returns the path to the metrics persistence file, using an environment variable if set,
// or a default location in the user's cache directory.
func getMetricsFilePath() string {
	if envPath := os.Getenv("LOGZ_METRICS_FILE"); envPath != "" {
		return envPath
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "/tmp"
	}
	dir := filepath.Join(cacheDir, "kubex", "logz")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "metrics.json")
}

// GetPrometheusManager returns the singleton instance of PrometheusManager, initializing it if necessary.
func GetPrometheusManager() *PrometheusManager {
	if prometheusManagerInstance == nil {
		prometheusManagerInstance = &PrometheusManager{
			enabled:         false,
			Metrics:         make(map[string]Metric),
			MetricsFile:     getMetricsFilePath(),
			exportWhitelist: make(map[string]bool),
		}
		if err := prometheusManagerInstance.LoadMetrics(); err != nil {
			fmt.Printf("Warning: could not load metrics: %v\n", err)
		}
	}
	return prometheusManagerInstance
}

// LoadMetrics loads metrics from the persistence file into the PrometheusManager instance.
func (pm *PrometheusManager) LoadMetrics() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	data, err := os.ReadFile(pm.MetricsFile)
	if err != nil {
		if os.IsNotExist(err) {
			pm.Metrics = make(map[string]Metric)
			return nil
		}
		return err
	}
	var loaded map[string]Metric
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	pm.Metrics = loaded
	return nil
}

// SaveMetrics saves the current metrics to the persistence file.
func (pm *PrometheusManager) SaveMetrics() error {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	data, err := json.MarshalIndent(pm.Metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pm.MetricsFile, data, 0644)
}

// Enable starts the Prometheus HTTP server on the specified port to expose metrics.
func (pm *PrometheusManager) Enable(port string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if pm.enabled {
		fmt.Println("Prometheus metrics are already enabled.")
		return
	}
	pm.enabled = true

	// Start the HTTP server to expose metrics
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := pm.GetMetrics()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		for name, value := range metrics {
			_, fPrintFErr := fmt.Fprintf(w, "# TYPE %s gauge\n%s %f\n", name, name, value)
			if fPrintFErr != nil {
				return
			}
		}
	})
	pm.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	go func() {
		if err := pm.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("ErrorCtx starting Prometheus metrics server: %v\n", err)
		}
	}()
	fmt.Println("Prometheus metrics enabled.")
}

// Disable stops the Prometheus HTTP server and disables metric exposure.
func (pm *PrometheusManager) Disable() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if !pm.enabled {
		fmt.Println("Prometheus metrics are already disabled.")
		return
	}
	pm.enabled = false
	if pm.httpServer != nil {
		_ = pm.httpServer.Close()
	}
	fmt.Println("Prometheus metrics disabled.")
}

// GetMetrics returns the current metrics, filtered by the export whitelist if defined.
func (pm *PrometheusManager) GetMetrics() map[string]float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	filteredMetrics := make(map[string]float64)
	for name, metric := range pm.Metrics {
		// Respect the exportWhitelist, if defined
		if len(pm.exportWhitelist) > 0 && !pm.exportWhitelist[name] {
			continue
		}
		filteredMetrics[name] = metric.Value
	}
	return filteredMetrics
}

// SetExportWhitelist sets the list of metrics that are allowed to be exported to Prometheus.
func (pm *PrometheusManager) SetExportWhitelist(metrics []string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.exportWhitelist = make(map[string]bool)
	for _, m := range metrics {
		pm.exportWhitelist[m] = true
	}
	fmt.Println("Export whitelist updated for Prometheus metrics.")
}

// IsEnabled returns whether the Prometheus metrics exposure is enabled.
func (pm *PrometheusManager) IsEnabled() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.enabled
}

// AddMetric adds or updates a metric with the given name, value, and VMetadata.
func (pm *PrometheusManager) AddMetric(name string, value float64, metadata map[string]string) {
	if err := validateMetricName(name); err != nil {
		fmt.Printf("ErrorCtx adding metric: %v\n", err)
		return
	}
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.Metrics[name] = Metric{
		Value:    value,
		Metadata: metadata,
	}
	fmt.Printf("Metric '%s' added/updated with value: %f\n", name, value)
	if err := pm.SaveMetrics(); err != nil {
		fmt.Printf("ErrorCtx saving metrics: %v\n", err)
	}
}

// RemoveMetric removes a metric with the given name.
func (pm *PrometheusManager) RemoveMetric(name string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.Metrics, name)
	fmt.Printf("Metric '%s' removed.\n", name)
	if err := pm.SaveMetrics(); err != nil {
		fmt.Printf("ErrorCtx saving metrics: %v\n", err)
	}
}

// IncrementMetric increments the value of a metric by the given delta.
func (pm *PrometheusManager) IncrementMetric(name string, delta float64) {
	if err := validateMetricName(name); err != nil {
		fmt.Printf("ErrorCtx incrementing metric: %v\n", err)
		return
	}
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	metric, exists := pm.Metrics[name]
	if !exists {
		metric = Metric{Value: 0, Metadata: nil}
	}
	metric.Value += delta
	pm.Metrics[name] = metric
	fmt.Printf("Metric '%s' incremented by %f, new value: %f\n", name, delta, metric.Value)
	if err := pm.SaveMetrics(); err != nil {
		fmt.Printf("ErrorCtx saving metrics: %v\n", err)
	}
}

// ListMetrics prints all registered metrics to the console.
func (pm *PrometheusManager) ListMetrics() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	if len(pm.Metrics) == 0 {
		fmt.Println("No metrics registered.")
		return
	}
	fmt.Println("Registered metrics:")
	for name, metric := range pm.Metrics {
		fmt.Printf("- %s: %f", name, metric.Value)
		if len(metric.Metadata) > 0 {
			metadataJSON, _ := json.Marshal(metric.Metadata)
			fmt.Printf(" (VMetadata: %s)", string(metadataJSON))
		}
		fmt.Println()
	}
}

// setPrometheusSysConfig configures the Prometheus system to scrape metrics from this application.
func (pm *PrometheusManager) setPrometheusSysConfig() error {
	// Logz specific configuration for Prometheus
	prometheusConfig := `
scrape_configs:
  - job_name: 'logz'
    static_configs:
      - targets: ['localhost:2112']
`

	configFilePath := "/etc/prometheus/prometheus.yml"

	// Check if the configuration file exists
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// Create the configuration file if it does not exist
		if err := os.WriteFile(configFilePath, []byte(prometheusConfig), 0644); err != nil {
			return fmt.Errorf("failed to create Prometheus configuration file: %w", err)
		}
		fmt.Println("Prometheus configuration file created successfully.")
	} else if err == nil {
		// Check if there is already a configuration for 'logz'
		configContent, readErr := os.ReadFile(configFilePath)
		if readErr != nil {
			return fmt.Errorf("failed to read Prometheus configuration file: %w", readErr)
		}

		if strings.Contains(string(configContent), "job_name: 'logz'") {
			fmt.Println("Prometheus configuration for 'logz' already exists.")
		} else {
			// Add the configuration to the existing file
			f, openErr := os.OpenFile(configFilePath, os.O_APPEND|os.O_WRONLY, 0644)
			if openErr != nil {
				return fmt.Errorf("failed to open Prometheus configuration file: %w", openErr)
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			if _, writeErr := f.WriteString(prometheusConfig); writeErr != nil {
				return fmt.Errorf("failed to append to Prometheus configuration file: %w", writeErr)
			}
			fmt.Println("Prometheus configuration for 'logz' added successfully.")
		}
	} else {
		return fmt.Errorf("failed to check Prometheus configuration file: %w", err)
	}

	return nil
}

// initPrometheus initializes the Prometheus metrics and system configuration.
func (pm *PrometheusManager) initPrometheus() error {
	if !pm.IsEnabled() {
		return fmt.Errorf("prometheus is not enabled")
	}

	defaultMetrics := []string{"infoCount", "warnCount", "errorCount", "debugCount", "successCount"}
	for _, metric := range defaultMetrics {
		if err := validateMetricName(metric); err != nil {
			fmt.Printf("ErrorCtx initializing metric '%s': %v\n", metric, err)
			continue
		}
		pm.AddMetric(metric, 0, nil) // Initialize with value 0 and no VMetadata
	}

	if err := pm.setPrometheusSysConfig(); err != nil {
		return fmt.Errorf("failed to configure Prometheus system: %w", err)
	}

	fmt.Println("Prometheus initialized successfully with default metrics.")
	return nil
}
