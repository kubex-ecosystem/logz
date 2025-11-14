// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type InitArgs = KBXConfig

func NewInitArgs(
	configFile string,
	configType string,
	configDBFile string,
	configDBType string,
	envFile string,
	logFile string,

	cwd string,
	tempDir string,

	// name string,
	debug bool,
	releaseMode bool,
	isConfidential bool,
	background bool,

	port string,
	bind string,
	pubCertKeyPath string,
	pubKeyPath string,
	privKeyPath string,

	timeout time.Duration,
	historySize int,
	defaultProvider string,
	defaultTemperature float32,
	providerConfigPath string,

	openAIConfig map[string]any,
	chatGPTConfig map[string]any,
	claudeConfig map[string]any,
	anthropicConfig map[string]any,
	deepSeekConfig map[string]any,
	geminiConfig map[string]any,
	ollamaConfig map[string]any,

	notificationProvider any,
	notificationTimeoutSeconds int,
) *InitArgs {
	// Prepare with default values fallbacks, just what will be used inside another variables. Otherwise,
	// we will use the direct values and respective fallbacks in the struct initialization
	configFile = GetValueOrDefaultSimple(configFile, os.ExpandEnv(DefaultGoBEConfigPath))
	configDBFile = GetValueOrDefaultSimple(configDBFile, "dbconfig.json")
	envFile = GetValueOrDefaultSimple(envFile, os.ExpandEnv(filepath.Join("$PWD", ".env")))
	logFile = GetValueOrDefaultSimple(
		logFile,
		filepath.Join(filepath.Dir(filepath.Dir(os.ExpandEnv(os.ExpandEnv(DefaultGoBEConfigPath)))), "logs", "gobe.log"),
	)
	port = GetValueOrDefaultSimple(port, "8088")
	bind = GetValueOrDefaultSimple(bind, "0.0.0.0")

	openAIConfig = GetValueOrDefaultSimple(openAIConfig, map[string]any{})
	chatGPTConfig = GetValueOrDefaultSimple(chatGPTConfig, map[string]any{})
	claudeConfig = GetValueOrDefaultSimple(claudeConfig, map[string]any{})
	anthropicConfig = GetValueOrDefaultSimple(anthropicConfig, map[string]any{})
	deepSeekConfig = GetValueOrDefaultSimple(deepSeekConfig, map[string]any{})
	geminiConfig = GetValueOrDefaultSimple(geminiConfig, map[string]any{})
	ollamaConfig = GetValueOrDefaultSimple(ollamaConfig, map[string]any{})

	return &InitArgs{
		&Reference{
			ID:   uuid.New(),
			Name: "GoBE",
		},
		&ConfigPaths{
			// Basic path and file configurations
			ConfigFile:   configFile,
			ConfigType:   filepath.Ext(configFile)[1:],
			ConfigDBFile: configDBFile,
			ConfigDBType: filepath.Ext(configDBFile)[1:],
			EnvFile:      envFile,
			LogFile:      logFile,
		},
		&ConfigScope{
			// Scope and runtime settings
			Cwd:     GetValueOrDefaultSimple(cwd, ""),
			TempDir: os.TempDir(),
		},
		&ConfigBasic{
			// Basic application settings
			// Name:           GetValueOrDefaultSimple(name, "GoBE"),
			Debug:          GetValueOrDefaultSimple(debug, false),
			ReleaseMode:    GetValueOrDefaultSimple(releaseMode, false),
			IsConfidential: GetValueOrDefaultSimple(isConfidential, false),
			Background:     background,
		},
		&ConfigNetworkSecurity{
			// Network and security settings
			Port:           port,
			Bind:           bind,
			PubCertKeyPath: GetValueOrDefaultSimple(pubCertKeyPath, os.ExpandEnv(DefaultGoBEKeyPath)),
			PubKeyPath:     GetValueOrDefaultSimple(pubKeyPath, os.ExpandEnv(DefaultGoBECertPath)),
		},
		&ConfigAIResiliency{
			Timeout:            60 * time.Second,
			HistorySize:        25,
			DefaultProvider:    "openai",
			DefaultTemperature: 0,
			ProviderConfigPath: GetValueOrDefaultSimple(providerConfigPath, ""),
		},
		&ConfigAIIntegrations{
			// AI Integrations and Settings
			OpenAIConfig:    openAIConfig,
			ChatGPTConfig:   chatGPTConfig,
			ClaudeConfig:    claudeConfig,
			AnthropicConfig: anthropicConfig,
			DeepSeekConfig:  deepSeekConfig,
			GeminiConfig:    geminiConfig,
			OllamaConfig:    ollamaConfig,
		},
		&ConfigExtended{
			// Extended settings
			NotificationProvider:       notificationProvider,
			NotificationTimeoutSeconds: GetValueOrDefaultSimple(notificationTimeoutSeconds, 5),
		},
		&LktSettings{
			ExtractOptions: &ExtractOptions{},
		},
	}
}
