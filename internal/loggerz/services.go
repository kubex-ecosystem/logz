package loggerz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/kubex-ecosystem/logz/api/notifiers"
	"github.com/kubex-ecosystem/logz/internal/interfaces"
	li "github.com/kubex-ecosystem/logz/internal/interfaces"

	"github.com/spf13/viper"
)

const (
	// pidFile is the name of the PID file used to track the running service.
	pidFile = "logz_srv.pid"
)

var (
	lSrv    *http.Server
	lClient *http.Client
	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	//lSocket      *zmq4.Socket
	lDBus        *dbus.Conn
	globalLogger li.LogzLogger // Global core for the service
	startTime    = time.Now()
	mu           sync.RWMutex
)

// Run starts the logging service.
func Run() error {
	mu.Lock()
	defer mu.Unlock()

	// Check if the service is already running to avoid multiple instances
	if IsRunning() {
		if stopErr := shutdown(); stopErr != nil {
			return stopErr
		}
	}

	// Initialize the ConfigManager and load the configuration
	configManager := interfaces.NewConfigManager()
	if configManager == nil {
		return errors.New("failed to initialize VConfig manager")
	}
	cfgMgr := configManager

	config, err := cfgMgr.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load VConfig: %w", err)
	}

	// Initialize the global core with the configuration
	initializeGlobalLogger(config)

	// Set up the HTTP server
	mux := http.NewServeMux()
	if err := registerHandlers(mux); err != nil {
		return err
	}

	lSrv = &http.Server{
		Addr:         config.Address(),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		IdleTimeout:  config.IdleTimeout(),
	}

	// Start the HTTP server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		globalLogger.InfoCtx(fmt.Sprintf("Service running on %s", config.Address()), nil)
		if err := lSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			globalLogger.ErrorCtx(fmt.Sprintf("Service encountered an error: %v", err), nil)
		}
	}()

	<-stop
	return shutdown()
}

// Start initiates the logging service on the specified port.
func Start(port string) error {
	mu.Lock()
	defer mu.Unlock()

	if IsRunning() {
		return errors.New("service already running (pid file exists: " + getPidPath() + ")")
	}

	// Use Viper to load runtime configuration
	vpr := viper.GetViper()
	if vpr == nil {
		return errors.New("viper not initialized")
	}

	cmd := exec.Command(os.Args[0], "service", "spawn", "-c", vpr.ConfigFileUsed())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	file, err := os.OpenFile(getPidPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open PID file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	//if runtime.GOOS != "windows" {
	//	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
	//		return errors.New("another process is writing to the PID file")
	//	}
	//}

	pid := cmd.Process.Pid
	pidData := fmt.Sprintf("%d\n%s", pid, port)
	if _, writeErr := file.Write([]byte(pidData)); writeErr != nil {
		return fmt.Errorf("failed to write PID data: %w", writeErr)
	}

	globalLogger.InfoCtx(fmt.Sprintf("Service started with pid %d", pid), nil)
	return nil
}

// Stop terminates the running logging service.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	pid, port, pidPath, err := GetServiceInfo()
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	time.Sleep(1 * time.Second)
	if err := os.Remove(pidPath); err != nil {
		return err
	}

	globalLogger.InfoCtx(fmt.Sprintf("Service with pid %d and port %s stopped", pid, port), nil)
	return nil
}

// Server returns the HTTP server instance.
func Server() *http.Server {
	return lSrv
}

// Client returns the HTTP client instance.
func Client() *http.Client {
	if lClient == nil {
		lClient = &http.Client{}
	}
	return lClient
}

// Temporarily disabled due to external dependency on zmq4
// Uncomment and ensure the required libraries are installed if needed in the future
// Socket returns the ZMQ socket instance.
//func Socket() *zmq4.Socket {
//	if lSocket == nil {
//		lSocket, _ = zmq4.NewSocket(zmq4.PUB)
//	}
//	return lSocket
//}

// DBus returns the DBus connection instance.
func DBus() *dbus.Conn {
	if lDBus == nil {
		lDBus, _ = dbus.SystemBus()
	}
	return lDBus
}

// getPidPath returns the path to the PID file.
func getPidPath() string {
	if envPath := os.Getenv("LOGZ_PID_PATH"); envPath != "" {
		return envPath
	}
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil {
		cacheDir = "/tmp"
	}
	cacheDir = filepath.Join(cacheDir, "logz", pidFile)
	if mkdirErr := os.MkdirAll(filepath.Dir(cacheDir), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return cacheDir
}

// IsRunning checks if the service is currently running.
func IsRunning() bool {
	_, err := os.Stat(getPidPath())
	return err == nil
}

// GetServiceInfo retrieves the PID, port, and PID file path of the running service.
func GetServiceInfo() (int, string, string, error) {
	pidPath := getPidPath()

	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, "", "", os.ErrNotExist
	}

	lines := strings.Split(string(data), "\n")
	pid, pidErr := strconv.Atoi(lines[0])
	if pidErr != nil {
		return 0, "", "", os.ErrInvalid
	}

	port := "unknown"
	if len(lines) > 1 {
		port = lines[1]
	}

	return pid, port, pidPath, nil
}

// registerHandlers registers HTTP handlers for the service.
func registerHandlers(mux *http.ServeMux) error {
	integrations := viper.GetStringMap("integrations")
	if integrations == nil {
		return errors.New("no integrations configured")
	}

	for path := range integrations {
		if !viper.GetBool("integrations." + path + ".enabled") {
			continue
		}

		healthPath, _ := url.JoinPath("/", path, "/health")
		metricsPath, _ := url.JoinPath("/", path, "/metrics")
		callbackPath, _ := url.JoinPath("/", path, "/receive")

		mux.HandleFunc(healthPath, healthHandler)
		mux.HandleFunc(metricsPath, metricsHandler)
		mux.HandleFunc(callbackPath, callbackHandler)
	}

	return nil
}

// callbackHandler handles incoming callback requests.
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit the payload size to prevent abuse
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if _, ok := payload["message"]; !ok {
		http.Error(w, "Missing 'message' in payload", http.StatusBadRequest)
		return
	}

	globalLogger.InfoCtx(fmt.Sprintf("Callback received: %v", payload), nil)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success","message":"Callback processed"}`))
}

// healthHandler handles health check requests.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(startTime).String()
	response := fmt.Sprintf("OK\nUptime: %s\n", uptime)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

// metricsHandler handles metrics requests.
func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	pm := notifiers.GetPrometheusManager()
	if !pm.IsEnabled() {
		http.Error(w, "Prometheus integration is not enabled", http.StatusForbidden)
		return
	}

	metrics := pm.GetMetrics()
	if len(metrics) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for name, value := range metrics {
		if _, err := fmt.Fprintf(w, "# HELP %s Custom metric from Logz\n# TYPE %s gauge\n%s %f\n", name, name, name, value); err != nil {
			fmt.Println(fmt.Printf("ErrorCtx writing metric '%s': %v", name, err))
		}
	}
}

// loggingMiddleware logs incoming HTTP requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// shutdown gracefully shuts down the service.
func shutdown() error {
	globalLogger.InfoCtx("Shutting down service gracefully...", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := lSrv.Shutdown(ctx); err != nil {
		globalLogger.ErrorCtx(fmt.Sprintf("Service shutdown failed: %v", err), nil)
		return fmt.Errorf("shutdown process failed: %w", err)
	}

	globalLogger.InfoCtx("Service stopped gracefully.", nil)
	return nil
}

// initializeGlobalLogger initializes the global core with the provided configuration.
func initializeGlobalLogger(config li.Config) {
	mu.Lock()
	defer mu.Unlock()

	if globalLogger == nil {
		globalLogger = NewLogger("Logz")
	}
}
