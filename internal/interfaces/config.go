package interfaces

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kubex-ecosystem/logz/internal/loggerz/formatters"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/spf13/viper"
)

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
	SetLevel(VLevel string)
	Format() string
	SetFormat(logFormat string)
	GetInt(key string, value int) int
	GetFormatter() interface{}
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

// NewConfigManager creates a new instance of ConfigManager.
func NewConfigManager() ConfigManager {
	cfgMgr := &

	if cfg, err := cfgMgr.LoadConfig(); err != nil || cfg == nil {
		log.Printf("ErrorCtx loading configuration: %v\n", err)
		return nil
	}

	var cfgM ConfigManagerImpl = cfgMgr

	return &cfgM
}

const (
	defaultPort        = "9999"
	defaultBindAddress = "0.0.0.0"
	defaultMode        = apiC.ModeStandalone
)

// LogzConfig specific to Logz
type LogzConfig struct {
	LogLevel     string
	LogFormat    string
	LogFilePath  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PidFile      string
}

var defaultLogPath = "stdout"

// ConfigImpl implements the Config interface and holds the configuration values.
type ConfigImpl struct {
	// Config is a constraint to implement Config interface
	*LogzConfig

	VlLevel           kbx.LogLevel
	VlFormat          kbx.LogFormat
	VlPort            string
	VlBindAddress     string
	VlAddress         string
	VlPidFile         string
	VlReadTimeout     time.Duration
	VlWriteTimeout    time.Duration
	VlIdleTimeout     time.Duration
	VlOutput          string
	VlNotifierManager NotifierManager
	VlMode            kbx.LogMode
}

func (c *ConfigImpl) GetFormatter() interface{} {
	switch c.Format() {
	case "json":
		return &formatters.NewJSONFormatterImpl().Format
	default:
		return &formatters.NewTableFormatterImpl().Format
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
func (c *ConfigImpl) SetLevel(VLevel string)       { c.VlLevel = kbx.LogLevel(VLevel) }
func (c *ConfigImpl) Format() string               { return strings.ToLower(string(c.VlFormat)) }
func (c *ConfigImpl) SetFormat(format string)      { c.VlFormat = kbx.LogFormat(format) }
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
