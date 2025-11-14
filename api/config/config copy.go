package config

import "time"

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
