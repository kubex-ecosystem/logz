package logger

import (
	//"errors"
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	"github.com/godbus/dbus/v5"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type LogzService struct {
	Server    *http.Server
	Client    *http.Client
	DBusConn  *dbus.Conn
	Logger    *LogzCoreImpl
	StartTime time.Time
}


var (
	lSrv    *http.Server
	lClient *http.Client
	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	//lSocket      *zmq4.Socket
	lDBus        *dbus.Conn
	globalLogger *LogzCoreImpl // Global logger for the service
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

	// Caso contrário, forneça uma configuração padrão
	defaultConfig := LogzConfig{
		LogLevel:     "INFO",
		LogFilePath:  "/tmp/logz.log",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	cfg := c.NewConfigManager(defaultConfig)
	if cfg == nil {
		return nil, fmt.Errorf("falha ao criar o gerenciador de configuração")
	}

	// Salve a configuração padrão na instância global do logger
	if loggerInstance != nil {
		loggerInstance.Config = cfg
	}

	return cfg, nil
}

func IsRunning() bool {
	lgConfig := configManager.GetConfig()
	_, err := os.Stat(lgConfig.PidFile)
	return err == nil
}

func Run() error {
	// Verifica se já está rodando
	if IsRunning() {
		if err := Stop(); err != nil {
			return fmt.Errorf("erro ao parar instância antiga: %w", err)
		}
	}

	// Inicializa o ConfigManager
	_, err := InitConfigManager()
	if err != nil {
		return fmt.Errorf("falha ao inicializar gerenciador de configuração: %w", err)
	}

	//config := cfg.GetConfig()

	// Configura o servidor HTTP
	mux := http.NewServeMux()
	if err := registerHandlers(mux); err != nil {
		return fmt.Errorf("erro ao registrar handlers HTTP: %w", err)
	}

	service.Server = &http.Server{
		//Addr:         config.Address(),
		//Handler:      loggingMiddleware(mux),
		//ReadTimeout:  config.ReadTimeout(),
		//WriteTimeout: config.WriteTimeout(),
		//IdleTimeout:  config.IdleTimeout(),
	}

	// Gerenciamento de sinal para interrupção
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//
	//go func() {
	//	service.Logger.Info(fmt.Sprintf("Serviço rodando em %s", config.Address()), nil)
	//	if err := service.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		service.Logger.Error(fmt.Sprintf("Erro no servidor: %v", err), nil)
	//	}
	//}()
	//
	//<-stop
	return Stop()
}

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

	cmd := exec.Command(os.Args[0], "service", "spawn", "-p", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("falha ao iniciar serviço: %w", err)
	}

	file, err := os.OpenFile(pidPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo PID: %w", err)
	}
	defer file.Close()

	pid := cmd.Process.Pid
	if _, writeErr := file.WriteString(fmt.Sprintf("%d\n%s", pid, port)); writeErr != nil {
		return fmt.Errorf("falha ao escrever PID: %w", writeErr)
	}

	service.Logger.Info(fmt.Sprintf("Serviço iniciado com PID %d na porta %s", pid, port), nil)
	return nil
}

func Stop() error {

	mu.Lock()
	defer mu.Unlock()

	pid, port, pidPath, err := GetServiceInfo()

	if err != nil {
		return fmt.Errorf("não foi possível obter informações do serviço: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("processo não encontrado: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("erro ao enviar SIGTERM: %w", err)
	}

	if err := os.Remove(pidPath); err != nil {
		return fmt.Errorf("erro ao remover arquivo PID: %w", err)
	}

	service.Logger.Info(fmt.Sprintf("Serviço com PID %d foi encerrado.", pid), nil)
	return nil
}

func GetServiceInfo() (int, string, string, error) {
	pidPath := getPidPath()
	file, err := os.Open(pidPath)
	if err != nil {
		return 0, "", "", fmt.Errorf("erro ao abrir arquivo PID: %w", err)
	}
	defer file.Close()

	var pid int
	var port string
	if _, err := fmt.Fscanf(file, "%d\n%s", &pid, &port); err != nil {
		return 0, "", "", fmt.Errorf("erro ao ler arquivo PID: %w", err)
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

	globalLogger.Info(fmt.Sprintf("Callback received: %v", payload), nil)
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
	pm := GetPrometheusManager()
	if (!pm.IsEnabled()) {
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
			fmt.Println(fmt.Sprintf("Error writing metric '%s': %v", name, err))
		}
	}

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.Logger.Info(fmt.Sprintf("Requisição recebida: %s %s", r.Method, r.URL.Path), nil)
		next.ServeHTTP(w, r)
	})
}

func registerHandlers(mux *http.ServeMux) error {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Método não permitido"))
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, writeErr := w.Write([]byte("Erro ao registrar log"))
			if writeErr != nil {
				return
			}
			return
		}

// initializeGlobalLogger initializes the global logger with the provided configuration.
func initializeGlobalLogger(config Config) {
	mu.Lock()
	defer mu.Unlock()

	if globalLogger == nil {
		globalLogger = NewLogger(config)
	}
}
