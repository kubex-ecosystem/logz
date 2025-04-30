package core

import (
	//"fmt"
	//c "github.com/faelmori/kubex-interfaces/config"

	//c "github.com/faelmori/kubex-interfaces/config"
	//"os"
	//"path/filepath"
	"time"
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

//func (lc *LogzConfig) Validate() error {
//	if lc.LogLevel == "" {
//		return fmt.Errorf("LogLevel cannot be empty")
//	}
//	if lc.LogFilePath == "" {
//		return fmt.Errorf("LogFilePath cannot be empty")
//	}
//	return nil
//}
//
//// LogzConfigManager with validation support
//type LogzConfigManager struct {
//	*LogzConfig
//	config *LogzConfig
//}

//func NewLogzConfigManager() *LogzConfigManager {
//	return NewConfigManager()
//	//return &LogzConfigManager{
//	//	LogzConfig: &cfg
//	//}
//}

//func (lcm *LogzConfigManager) ValidateConfig() error {
//	if config, ok := any(lcm).(c.Configurable); ok {
//		return config.Validate()
//	}
//	return fmt.Errorf("invalid configuration or validation not supported")
//}
//
//func (lcm *LogzConfigManager) GetConfig() *LogzConfig { return lcm.config }
//
//// GetPidPath returns the path to the PID file.
//func (lcm *LogzConfigManager) GetPidPath() string {
//	cacheDir, cacheDirErr := os.UserCacheDir()
//	if cacheDirErr != nil {
//		cacheDir = "/tmp"
//	}
//	cacheDir = filepath.Join(cacheDir, "logz", lcm.config.PidFile)
//	if mkdirErr := os.MkdirAll(filepath.Dir(cacheDir), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
//		return ""
//	}
//	return cacheDir
//}
//
//// GetConfigPath returns the path to the configuration file.
//func (lcm *LogzConfigManager) GetConfigPath() string {
//	if lcm.config != nil {
//		if lcm.Output() != "" {
//			return lcm.Output()
//		}
//	}
//
//	home, homeErr := os.UserHomeDir()
//	if homeErr != nil {
//		home, homeErr = os.UserConfigDir()
//		if homeErr != nil {
//			home, homeErr = os.UserCacheDir()
//			if homeErr != nil {
//				home = "/tmp"
//			}
//		}
//	}
//	configPath := filepath.Join(home, ".kubex", "logz", "VConfig.json")
//	if mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
//		return ""
//	}
//	return configPath
//}
//
//// SetOutput sets the path to the default log file.
//func (lcm *LogzConfigManager) SetOutput(output string) {
//	if lcm.config != nil {
//		lcm.SetOutput(output)
//	} else {
//		fmt.Println("Cannot set output in standalone VMode")
//	}
//}
//
//// Output returns the path to the configuration file.
//func (lcm *LogzConfigManager) Output() string {
//	if lcm.config != nil {
//		if lcm.Output() != "" {
//			return lcm.Output()
//		}
//	}
//	home, homeErr := os.UserHomeDir()
//	if homeErr != nil {
//		home, homeErr = os.UserConfigDir()
//		if homeErr != nil {
//			home, homeErr = os.UserCacheDir()
//			if homeErr != nil {
//				home = "/tmp"
//			}
//		}
//	}
//	logPath := filepath.Join(home, ".kubex", "logz", "logz.log")
//	if mkdirErr := os.MkdirAll(filepath.Dir(logPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
//		return ""
//	}
//	if _, err := os.Stat(logPath); os.IsNotExist(err) {
//		if _, createErr := os.Create(logPath); createErr != nil {
//			return ""
//		}
//	}
//	return logPath
//}
//
//func (lcm *LogzConfigManager) SetLevel(level LogLevel) { lcm.config.LogLevel = string(level) }
//
//func (lcm *LogzConfigManager) SetFormat(format LogFormat) {
//	lcm.config.LogFormat = string(format)
//}
//
//// GetFormatter returns the formatter for the core.
//func (lcm *LogzConfigManager) GetFormatter() LogFormatter {
//	switch lcm.config.LogFormat {
//	case "text":
//		return &TextFormatter{}
//	default:
//		return &JSONFormatter{}
//	}
//}

// LoadConfig loads the configuration from the file and returns a Config instance.
// func (lcm *LogzConfigManager) LoadConfig() (*LogzConfig, error) {
// 	cfg, err := InitConfigManager()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize the configuration manager: %w", err)
// 	}

// 	// Load the configuration from the file
// 	cf := cfg.GetConfig()

// 	// Validate the configuration
// 	if err := cf.Validate(); err != nil {
// 		return nil, fmt.Errorf("failed to validate the configuration: %w", err)
// 	}

// 	// Save the configuration
// 	lcm.VConfig = &cf

// 	return lcm.VConfig, nil
// }
