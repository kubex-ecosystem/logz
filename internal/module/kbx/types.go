package kbx

import (
	"time"

	"github.com/google/uuid"
)

// Reference is the internal struct that holds the server startup unique identifier with name.
type Reference struct {
	// refID is the unique identifier for this context.
	ID uuid.UUID `json:"id" yaml:"id" mapstructure:"id"`
	// refName is the name of the context.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// ConfigPaths holds the configuration file paths.
type ConfigPaths struct {
	// Basic path and file configurations
	ConfigFile   string `yaml:"config_file,omitempty" json:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	ConfigType   string `yaml:"config_type,omitempty" json:"config_type,omitempty" mapstructure:"config_type,omitempty"`
	ConfigDBFile string `yaml:"config_db_file,omitempty" json:"config_db_file,omitempty" mapstructure:"config_db_file,omitempty"`
	ConfigDBType string `yaml:"config_db_type,omitempty" json:"config_db_type,omitempty" mapstructure:"config_db_type,omitempty"`
	EnvFile      string `yaml:"env_file,omitempty" json:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	LogFile      string `yaml:"log_file,omitempty" json:"log_file,omitempty" mapstructure:"log_file,omitempty"`
}

// ConfigScope holds the scope and runtime settings.
type ConfigScope struct {
	// Scope and runtime settings
	Cwd     string `yaml:"cwd,omitempty" json:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	TempDir string `yaml:"temp_dir,omitempty" json:"temp_dir,omitempty" mapstructure:"temp_dir,omitempty"`
}

type ConfigBasic struct {
	// Basic application settings
	// Name           string `yaml:"name,omitempty" json:"name,omitempty" mapstructure:"name,omitempty"`
	Debug          bool `yaml:"debug,omitempty" json:"debug,omitempty" mapstructure:"debug,omitempty"`
	ReleaseMode    bool `yaml:"release_mode,omitempty" json:"release_mode,omitempty" mapstructure:"release_mode,omitempty"`
	IsConfidential bool `yaml:"is_confidential,omitempty" json:"is_confidential,omitempty" mapstructure:"is_confidential,omitempty"`
	Background     bool `yaml:"background,omitempty" json:"background,omitempty" mapstructure:"background,omitempty"`
}

type ConfigNetworkSecurity struct {
	// Network and security settings
	Port           string `yaml:"port,omitempty" json:"port,omitempty" mapstructure:"port,omitempty"`
	Bind           string `yaml:"bind,omitempty" json:"bind,omitempty" mapstructure:"bind,omitempty"`
	PubCertKeyPath string `yaml:"pub_cert_key_path,omitempty" json:"pub_cert_key_path,omitempty" mapstructure:"pub_cert_key_path,omitempty"`
	PubKeyPath     string `yaml:"pub_key_path,omitempty" json:"pub_key_path,omitempty" mapstructure:"pub_key_path,omitempty"`
	PrivKeyPath    string `yaml:"priv_key_path,omitempty" json:"priv_key_path,omitempty" mapstructure:"priv_key_path,omitempty"`
}

type ConfigAIResiliency struct {
	// AI Resiliency Configurations
	Timeout            time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" mapstructure:"timeout,omitempty"`
	HistorySize        int           `yaml:"history_size,omitempty" json:"history_size,omitempty" mapstructure:"history_size,omitempty"`
	DefaultProvider    string        `yaml:"default_provider,omitempty" json:"default_provider,omitempty" mapstructure:"default_provider,omitempty"`
	DefaultTemperature float32       `yaml:"default_temperature,omitempty" json:"default_temperature,omitempty" mapstructure:"default_temperature,omitempty"`
	ProviderConfigPath string        `yaml:"provider_config_path,omitempty" json:"provider_config_path,omitempty" mapstructure:"provider_config_path,omitempty"`
}

// ConfigAIIntegrations holds configuration settings for various AI integrations.
// Each field represents the configuration for a specific AI provider, stored as a map of key-value pairs.
// The supported providers include OpenAI, ChatGPT, Claude, Anthropic, DeepSeek, Gemini, and Ollama.
// Fields are tagged for YAML, JSON, and mapstructure serialization.
type ConfigAIIntegrations struct {
	// AI integrations and Settings
	OpenAIConfig    map[string]any `yaml:"openai_config,omitempty" json:"openai_config,omitempty" mapstructure:"openai_config,omitempty"`
	ChatGPTConfig   map[string]any `yaml:"chatgpt_config,omitempty" json:"chatgpt_config,omitempty" mapstructure:"chatgpt_config,omitempty"`
	ClaudeConfig    map[string]any `yaml:"claude_config,omitempty" json:"claude_config,omitempty" mapstructure:"claude_config,omitempty"`
	AnthropicConfig map[string]any `yaml:"anthropic_config,omitempty" json:"anthropic_config,omitempty" mapstructure:"anthropic_config,omitempty"`
	DeepSeekConfig  map[string]any `yaml:"deep_seek_config,omitempty" json:"deep_seek_config,omitempty" mapstructure:"deep_seek_config,omitempty"`
	GeminiConfig    map[string]any `yaml:"gemini_config,omitempty" json:"gemini_config,omitempty" mapstructure:"gemini_config,omitempty"`
	OllamaConfig    map[string]any `yaml:"ollama_config,omitempty" json:"ollama_config,omitempty" mapstructure:"ollama_config,omitempty"`
}

type ConfigExtended struct {
	// Extended configurations
	NotificationProvider       any `yaml:"notification_provider,omitempty" json:"notification_provider,omitempty" mapstructure:"notification_provider,omitempty"`
	NotificationTimeoutSeconds int `yaml:"notification_timeout_seconds,omitempty" json:"notification_timeout_seconds,omitempty" mapstructure:"notification_timeout_seconds,omitempty"`
}

// ExtractOptions defines options for file extraction.
type ExtractOptions struct {
	Overwrite  bool `json:"overwrite,omitempty" yaml:"overwrite,omitempty" mapstructure:"overwrite"`
	CreateDirs bool `json:"createDirs,omitempty" yaml:"createDirs,omitempty" mapstructure:"createDirs"`
	DryRun     bool `json:"dryRun,omitempty" yaml:"dryRun,omitempty" mapstructure:"dryRun"`
}

// LktSettings is the Lookatni File Markers specific settings
type LktSettings struct {
	*ExtractOptions `json:"extract_options,omitempty" yaml:"extract_options,omitempty" mapstructure:"extract_options,omitempty"`
}

// KBXConfig is the main configuration struct for the application.
type KBXConfig struct {
	*Reference             `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigPaths           `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigScope           `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigBasic           `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigNetworkSecurity `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigAIResiliency    `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigAIIntegrations  `yaml:",inline" json:",inline" mapstructure:",inline"`
	*ConfigExtended        `yaml:",inline" json:",inline" mapstructure:",inline"`
	*LktSettings           `yaml:",inline" json:",inline" mapstructure:",inline"`
}
