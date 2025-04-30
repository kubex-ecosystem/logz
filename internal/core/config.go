package core

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	// "github.com/fsnotify/fsnotify"
	// "github.com/spf13/viper"
	//
	// "encoding/json"
	// "fmt"
	// "log"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "sync"
	// "time"
	"strings"
	"time"
)

const (
	defaultPort        = "9999"
	defaultBindAddress = "0.0.0.0"
	defaultMode        = ModeStandalone
)

var defaultLogPath = "stdout"

// Config interface defines the methods to access configuration settings.
type Config interface {
	Port() string
	BindAddress() string
	Address() string
	PidFile() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	Output() string
	SetOutput(configPath string)
	NotifierManager() interface{}
	Mode() interface{}
	Level() string
	SetLevel(VLevel interface{})
	Format() string
	SetFormat(LogFormat interface{})
	GetInt(key string, value int) int
	GetFormatter() interface{}
}

// ConfigImpl implements the Config interface and holds the configuration values.
type ConfigImpl struct {
	// Config is a constraint to implement Config interface
	Config

	VlLevel           LogLevel
	VlFormat          LogFormat
	VlPort            string
	VlBindAddress     string
	VlAddress         string
	VlPidFile         string
	VlReadTimeout     time.Duration
	VlWriteTimeout    time.Duration
	VlIdleTimeout     time.Duration
	VlOutput          string
	VlNotifierManager NotifierManager
	VlMode            LogMode
}

func (c *ConfigImpl) GetFormatter() interface{} {
	switch c.Format() {
	case "json":
		return &JSONFormatter{}
	default:
		return &TextFormatter{}
	}
}
func (c *ConfigImpl) Port() string                 { return c.VlPort }
func (c *ConfigImpl) BindAddress() string          { return c.VlBindAddress }
func (c *ConfigImpl) Address() string              { return c.VlAddress }
func (c *ConfigImpl) PidFile() string              { return c.VlPidFile }
func (c *ConfigImpl) ReadTimeout() time.Duration   { return c.VlReadTimeout }
func (c *ConfigImpl) WriteTimeout() time.Duration  { return c.VlWriteTimeout }
func (c *ConfigImpl) IdleTimeout() time.Duration   { return c.VlIdleTimeout }
func (c *ConfigImpl) NotifierManager() interface{} { return c.VlNotifierManager }
func (c *ConfigImpl) Mode() interface{}            { return c.VlMode }
func (c *ConfigImpl) Level() string                { return strings.ToUpper(string(c.VlLevel)) }
func (c *ConfigImpl) SetLevel(VLevel interface{})  { c.VlLevel = LogLevel(VLevel.(string)) }
func (c *ConfigImpl) Format() string               { return strings.ToLower(string(c.VlFormat)) }
func (c *ConfigImpl) SetFormat(format interface{}) { c.VlFormat = LogFormat(format.(string)) }
func (c *ConfigImpl) Output() string {
	if c.VlOutput != "" {
		return c.VlOutput
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	logPath := filepath.Join(home, ".kubex", "logz", "logz.log")
	if mkdirErr := os.MkdirAll(filepath.Dir(logPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, createErr := os.Create(logPath); createErr != nil {
			return ""
		}
	}
	return logPath
}
func (c *ConfigImpl) SetOutput(configPath string) {
	c.VlOutput = configPath
}
func (c *ConfigImpl) GetInt(key string, defaultValue int) int {
	viperInstance := viper.GetViper()

	// Primeiro tenta buscar via Viper, se disponível
	if viperInstance != nil {
		// Obtém o valor como string para lidar com chaves configuradas em diferentes formatos
		rawValue := viperInstance.GetString(key)
		if rawValue != "" {
			parsedVal, err := strconv.Atoi(rawValue) // Converte o valor para inteiro
			if err == nil {
				return parsedVal
			}
		}
	}

	// Caso não encontre ou a conversão falhe, retorna o valor padrão
	return defaultValue
}
func (c *ConfigImpl) GetString(key string, defaultValue string) string {
	viperInstance := viper.GetViper()

	// Primeiro tenta buscar via Viper, se disponível
	if viperInstance != nil {
		// Obtém o valor como string para lidar com chaves configuradas em diferentes formatos
		rawValue := viperInstance.GetString(key)
		if rawValue != "" {
			return rawValue
		}
	}

	// Caso não encontre ou a conversão falhe, retorna o valor padrão
	return defaultValue
}

// ConfigManager interface defines methods to manage configuration.
type ConfigManager interface {
	GetConfig() Config
	GetPidPath() string
	GetConfigPath() string
	Output() string
	SetOutput(configPath string)
	LoadConfig() (Config, error)
}

// ConfigManagerImpl implements the ConfigManager interface.
type ConfigManagerImpl struct {
	VConfig Config
	Mu      sync.RWMutex
}

func (cm *ConfigManagerImpl) checkConfig() {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	if cm.VConfig == nil {
		cm.VConfig = &ConfigImpl{}
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	configPath := filepath.Join(home, ".kubex", "logz", "VConfig.json")
	if mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if _, createErr := os.Create(configPath); createErr != nil {
			return
		}
	}
}

func (cm *ConfigManagerImpl) Port() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	cm.checkConfig()
	if cm.VConfig != nil {
		return cm.VConfig.Port()
	}
	return defaultPort
}

func (cm *ConfigManagerImpl) BindAddress() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	cm.checkConfig()
	if cm.VConfig != nil {
		return cm.VConfig.BindAddress()
	}
	return defaultBindAddress
}

func (cm *ConfigManagerImpl) Address() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.Address()
	}
	return fmt.Sprintf("%s:%s", defaultBindAddress, defaultPort)
}

func (cm *ConfigManagerImpl) PidFile() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.PidFile()
	}
	return "logz_srv.pid"
}

func (cm *ConfigManagerImpl) ReadTimeout() time.Duration {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.ReadTimeout()
	}
	return time.Duration(15 * time.Second)
}

func (cm *ConfigManagerImpl) WriteTimeout() time.Duration {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.WriteTimeout()
	}
	return time.Duration(15 * time.Second)
}

func (cm *ConfigManagerImpl) IdleTimeout() time.Duration {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.IdleTimeout()
	}
	return time.Duration(60 * time.Second)
}

func (cm *ConfigManagerImpl) NotifierManager() interface{} {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.NotifierManager()
	}
	return nil
}

func (cm *ConfigManagerImpl) Mode() interface{} {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.Mode()
	}
	return defaultMode
}

func (cm *ConfigManagerImpl) Level() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.Level()
	}
	return strings.ToUpper(string(cm.VConfig.Level()))
}

func (cm *ConfigManagerImpl) Format() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.Format()
	}
	return strings.ToLower(string(cm.VConfig.Format()))
}

func (cm *ConfigManagerImpl) GetInt(key string, value int) int {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	viperInstance := viper.GetViper()
	if viperInstance != nil {
		rawValue := viperInstance.GetString(key)
		if rawValue != "" {
			parsedVal, err := strconv.Atoi(rawValue)
			if err == nil {
				return parsedVal
			}
		}
	}
	return value
}

func (cm *ConfigManagerImpl) GetConfig() Config {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	return cm.VConfig
}

// GetPidPath returns the path to the PID file.
func (cm *ConfigManagerImpl) GetPidPath() string {
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil {
		cacheDir = "/tmp"
	}
	cacheDir = filepath.Join(cacheDir, "logz", cm.VConfig.PidFile())
	if mkdirErr := os.MkdirAll(filepath.Dir(cacheDir), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return cacheDir
}

// GetConfigPath returns the path to the configuration file.
func (cm *ConfigManagerImpl) GetConfigPath() string {
	if cm.VConfig != nil {
		if cm.VConfig.Output() != "" && cm.VConfig.Mode() == ModeService {
			return cm.VConfig.Output()
		}
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	configPath := filepath.Join(home, ".kubex", "logz", "VConfig.json")
	if mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return configPath
}

// SetOutput sets the path to the default log file.
func (cm *ConfigManagerImpl) SetOutput(output string) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	if cm.VConfig != nil {
		cm.VConfig.SetOutput(output)
	} else {

		if cm.VConfig.Mode() == ModeService {
			VConfig, configErr := cm.LoadConfig()
			if configErr != nil {
				log.Printf("ErrorCtx loading configuration: %v\n", configErr)
				return
			}
			VConfig.SetOutput(output)
			cm.VConfig = VConfig
		} else {
			log.Printf("Cannot set output in standalone VMode\n")
		}

	}
}

// Output returns the path to the configuration file.
func (cm *ConfigManagerImpl) Output() string {
	if cm.VConfig != nil {
		if cm.VConfig.Output() != "" {
			return cm.VConfig.Output()
		}
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	logPath := filepath.Join(home, ".kubex", "logz", "logz.log")
	if mkdirErr := os.MkdirAll(filepath.Dir(logPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, createErr := os.Create(logPath); createErr != nil {
			return ""
		}
	}
	return logPath
}

func (cm *ConfigManagerImpl) SetLevel(VLevel interface{}) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	if cm.VConfig != nil {
		cm.VConfig.SetLevel(VLevel)
	} else {
		VConfig, configErr := cm.LoadConfig()
		if configErr != nil {
			log.Printf("ErrorCtx loading configuration: %v\n", configErr)
			return
		}
		VConfig.SetLevel(VLevel)
		cm.VConfig = VConfig
	}
}

func (cm *ConfigManagerImpl) SetFormat(format interface{}) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	if cm.VConfig != nil {
		cm.VConfig.SetFormat(format)
	} else {
		VConfig, configErr := cm.LoadConfig()
		if configErr != nil {
			log.Printf("ErrorCtx loading configuration: %v\n", configErr)
			return
		}
		VConfig.SetFormat(format)
		cm.VConfig = VConfig
	}
}

// GetFormatter returns the formatter for the core.
func (cm *ConfigManagerImpl) GetFormatter() interface{} {
	switch cm.VConfig.Format() {
	case "text":
		return &TextFormatter{}
	default:
		return &JSONFormatter{}
	}
}

// LoadConfig loads the configuration from the file and returns a Config instance.
func (cm *ConfigManagerImpl) LoadConfig() (Config, error) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	configPath := cm.GetConfigPath()
	if err := ensureConfigExists(configPath); err != nil {
		return nil, fmt.Errorf("failed to ensure VConfig exists: %w", err)
	}

	viperObj := viper.New()
	viperObj.SetConfigFile(configPath)
	viperObj.SetConfigType(getConfigType(configPath))

	if readErr := viperObj.ReadInConfig(); readErr != nil {
		return nil, fmt.Errorf("failed to read VConfig: %w", readErr)
	}

	notifierManager := NewNotifierManager(nil)
	if notifierManager == nil {
		return nil, fmt.Errorf("failed to create notifier manager")
	}

	VMode := LogMode(viperObj.GetString("VMode"))
	if VMode != ModeService && VMode != ModeStandalone {
		VMode = defaultMode
	}

	VConfig := ConfigImpl{
		VlPort:            getOrDefault(viperObj.GetString("port"), defaultPort),
		VlBindAddress:     getOrDefault(viperObj.GetString("bindAddress"), defaultBindAddress),
		VlAddress:         fmt.Sprintf("%s:%s", defaultBindAddress, defaultPort),
		VlPidFile:         viperObj.GetString("pidFile"),
		VlReadTimeout:     viperObj.GetDuration("readTimeout"),
		VlWriteTimeout:    viperObj.GetDuration("writeTimeout"),
		VlIdleTimeout:     viperObj.GetDuration("idleTimeout"),
		VlOutput:          getOrDefault(viperObj.GetString("defaultLogPath"), defaultLogPath),
		VlNotifierManager: notifierManager,
		VlMode:            VMode,
	}

	cm.VConfig = &VConfig

	viperObj.WatchConfig()
	viperObj.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuration changed: %s", e.Name)
		// Update Config dynamically, if necessary
	})

	return cm.VConfig, nil
}

// NewConfigManager creates a new instance of ConfigManager.
func NewConfigManager() *ConfigManager {
	cfgMgr := &ConfigManagerImpl{}

	if cfg, err := cfgMgr.LoadConfig(); err != nil || cfg == nil {
		log.Printf("ErrorCtx loading configuration: %v\n", err)
		return nil
	}

	var cfgM ConfigManager
	cfgM = cfgMgr

	return &cfgM
}

// ensureConfigExists checks if the configuration file exists, and creates it with default values if it does not.
func ensureConfigExists(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := ConfigImpl{
			VlPort:            defaultPort,
			VlBindAddress:     defaultBindAddress,
			VlAddress:         fmt.Sprintf("%s:%s", defaultBindAddress, defaultPort),
			VlPidFile:         "logz_srv.pid",
			VlReadTimeout:     15 * time.Second,
			VlWriteTimeout:    15 * time.Second,
			VlIdleTimeout:     60 * time.Second,
			VlOutput:          defaultLogPath,
			VlNotifierManager: NewNotifierManager(nil),
			VlMode:            defaultMode,
		}
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		if writeErr := os.WriteFile(configPath, data, 0644); writeErr != nil {
			return fmt.Errorf("failed to create default VConfig: %w", writeErr)
		}
	}
	return nil
}

func getConfigType(configPath string) string {
	configType := filepath.Ext(configPath)
	switch configType {
	case ".yaml":
		return "yaml"
	case ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	default:
		return "json"
	}

}

// getOrDefault returns the value if it is not empty, otherwise returns the default value.
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
