package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kubex-ecosystem/logz/internal/formatters"
	"github.com/kubex-ecosystem/logz/internal/interfaces"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"

	// il "github.com/kubex-ecosystem/logz/internal/ine"
	"github.com/spf13/viper"
)

// ConfigManagerImpl implements the ConfigManager interface.
type ConfigManagerImpl struct {
	VConfig interfaces.Config
	Mu      sync.RWMutex
}

func (cm *ConfigManagerImpl) checkConfig() {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	if cm.VConfig == nil {
		cm.VConfig = &interfaces.ConfigImpl{}
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
	return kbx.DefaultServerPort
}

func (cm *ConfigManagerImpl) BindAddress() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	cm.checkConfig()
	if cm.VConfig != nil {
		return cm.VConfig.BindAddress()
	}
	return kbx.DefaultServerHost
}

func (cm *ConfigManagerImpl) Address() string {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	if cm.VConfig != nil {
		return cm.VConfig.Address()
	}
	return fmt.Sprintf("%s:%s", kbx.DefaultServerHost, kbx.DefaultServerPort)
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
	return kbx.DefaultMode
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

func (cm *ConfigManagerImpl) GetConfig() *interfaces.ConfigImpl {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	return cm.VConfig.(*interfaces.ConfigImpl)
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
		if cm.VConfig.Output() != "" && cm.VConfig.Mode() == kbx.ModeService {
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

		if cm.VConfig.Mode() == kbx.ModeService {
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

func (cm *ConfigManagerImpl) SetLevel(VLevel string) {
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

func (cm *ConfigManagerImpl) SetFormat(format string) {
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
		return &formatters.TextFormatter{}
	default:
		return &formatters.JSONFormatter{}
	}
}

// LoadConfig loads the configuration from the file and returns a Config instance.
func (cm *ConfigManagerImpl) LoadConfig() (*interfaces.ConfigImpl, error) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	configPath := cm.GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if _, createErr := os.Create(configPath); createErr != nil {
			return nil, fmt.Errorf("failed to create config file: %w", createErr)
		}
	}

	viperObj := viper.New()
	viperObj.SetConfigFile(configPath)
	viperObj.SetConfigType(getConfigType(configPath))
	viperObj.SetDefault("VMode", defaultMode)

	if readErr := viperObj.ReadInConfig(); readErr != nil {
		return nil, fmt.Errorf("failed to read VConfig: %w", readErr)
	}

	notifierManager := kbx.NewNotifierManager(nil)
	if notifierManager == nil {
		return nil, fmt.Errorf("failed to create notifier manager")
	}

	VMode := kbx.LogMode(viperObj.GetString("VMode"))
	if VMode != kbx.ModeService && VMode != kbx.ModeStandalone {
		VMode = kbx.DefaultMode
	}

	VConfig := &interfaces.ConfigImpl{
		VlPort:            kbx.GetValueOrDefaultSimple(viperObj.GetString("port"), defaultPort),
		VlBindAddress:     kbx.GetValueOrDefaultSimple(viperObj.GetString("bindAddress"), defaultBindAddress),
		VlAddress:         net.JoinHostPort(kbx.DefaultBindAddress, kbx.DefaultPort),
		VlPidFile:         viperObj.GetString("pidFile"),
		VlReadTimeout:     viperObj.GetDuration("readTimeout"),
		VlWriteTimeout:    viperObj.GetDuration("writeTimeout"),
		VlIdleTimeout:     viperObj.GetDuration("idleTimeout"),
		VlOutput:          kbx.GetValueOrDefaultSimple(viperObj.GetString("defaultLogPath"), defaultLogPath),
		VlNotifierManager: notifierManager,
		VlMode:            VMode,
	}

	cm.VConfig = VConfig

	viperObj.WatchConfig()
	viperObj.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuration changed: %s", e.Name)
		// Update Config dynamically, if necessary
	})

	return cm.VConfig.(*interfaces.ConfigImpl), nil
}
