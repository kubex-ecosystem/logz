package logger

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
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
	NotifierManager() NotifierManager
	Mode() LogMode
	Level() string
	SetLevel(level LogLevel)
	Format() string
	SetFormat(LogFormat LogFormat)
	GetInt(key string, value int) int
	GetFormatter() LogFormatter
}

// ConfigImpl implements the Config interface and holds the configuration values.
type ConfigImpl struct {
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

func (c *ConfigImpl) GetFormatter() LogFormatter {
	switch c.Format() {
	case "json":
		return &JSONFormatter{}
	default:
		return &TextFormatter{}
	}
}
func (c *ConfigImpl) Port() string                     { return c.VlPort }
func (c *ConfigImpl) BindAddress() string              { return c.VlBindAddress }
func (c *ConfigImpl) Address() string                  { return c.VlAddress }
func (c *ConfigImpl) PidFile() string                  { return c.VlPidFile }
func (c *ConfigImpl) ReadTimeout() time.Duration       { return c.VlReadTimeout }
func (c *ConfigImpl) WriteTimeout() time.Duration      { return c.VlWriteTimeout }
func (c *ConfigImpl) IdleTimeout() time.Duration       { return c.VlIdleTimeout }
func (c *ConfigImpl) NotifierManager() NotifierManager { return c.VlNotifierManager }
func (c *ConfigImpl) Mode() LogMode                    { return c.VlMode }
func (c *ConfigImpl) Level() string                    { return strings.ToUpper(string(c.VlLevel)) }
func (c *ConfigImpl) SetLevel(level LogLevel)          { c.VlLevel = level }
func (c *ConfigImpl) Format() string                   { return strings.ToLower(string(c.VlFormat)) }
func (c *ConfigImpl) SetFormat(format LogFormat)       { c.VlFormat = format }
func (c *ConfigImpl) Output() string {
	if c.VlOutput == "" {
		return os.Stdout.Name()
	}
	return c.VlOutput
	//if c.VlOutput != "" {
	//	return c.VlOutput
	//}
	//home, homeErr := os.UserHomeDir()
	//if homeErr != nil {
	//	home, homeErr = os.UserConfigDir()
	//	if homeErr != nil {
	//		home, homeErr = os.UserCacheDir()
	//		if homeErr != nil {
	//			home = "/tmp"
	//		}
	//	}
	//}
	//logPath := filepath.Join(home, ".kubex", "logz", "logz.log")
	//if mkdirErr := os.MkdirAll(filepath.Dir(logPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
	//	return ""
	//}
	//if _, err := os.Stat(logPath); os.IsNotExist(err) {
	//	if _, createErr := os.Create(logPath); createErr != nil {
	//		return ""
	//	}
	//}
	//return logPath
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

// ConfigManagerImpl implements the ConfigManager interface.
type ConfigManagerImpl struct {
	config Config
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
			return fmt.Errorf("failed to create default config: %w", writeErr)
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
