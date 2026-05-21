// Package kbx provides utilities for working with configuration options.
package kbx

import (
	"io"
	"time"
)

// TelemetryConfig holds configuration for telemetry.
type TelemetryConfig struct {
	// The endpoint URL for the telemetry service
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint" json:"endpoint"`
	// API key for authentication
	APIKey string `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
	// Timeout for telemetry requests
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout" json:"timeout"`
	// Interval for sending telemetry data
	Interval time.Duration `mapstructure:"interval" yaml:"interval" json:"interval"`
	// Enable or disable telemetry
	Enabled bool `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
}

// LogzObservabilityOptions represents the observability options for the logger.
type LogzObservabilityOptions struct {
	EnablePrometheus *bool            `json:"enable_prometheus,omitempty" yaml:"enable_prometheus,omitempty" mapstructure:"enable_prometheus,omitempty"`
	MetricsPrefix    string           `json:"metrics_prefix,omitempty" yaml:"metrics_prefix,omitempty" mapstructure:"metrics_prefix,omitempty"`       // ex: "gnyx_" ou "domus_"
	PushGatewayURL   string           `json:"push_gateway_url,omitempty" yaml:"push_gateway_url,omitempty" mapstructure:"push_gateway_url,omitempty"` // Para jobs curtos (opcional)
	TelemetryConfig  *TelemetryConfig `json:"telemetry_config,omitempty" yaml:"telemetry_config,omitempty" mapstructure:"telemetry_config,omitempty"`
}

// LogzGeneralOptions represents the general options for the logger.
type LogzGeneralOptions struct {
	// General options
	Prefix      string `json:"prefix,omitempty" yaml:"prefix,omitempty" mapstructure:"prefix,omitempty"`
	Debug       bool   `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	ShowColor   *bool  `json:"show_color,omitempty" yaml:"show_color,omitempty" mapstructure:"show_color,omitempty"`
	ShowIcons   *bool  `json:"show_icons,omitempty" yaml:"show_icons,omitempty" mapstructure:"show_icons,omitempty"`
	ShowTraceID bool   `json:"show_trace_id,omitempty" yaml:"show_trace_id,omitempty" mapstructure:"show_trace_id,omitempty"`
	ShowStack   bool   `json:"show_caller,omitempty" yaml:"show_caller,omitempty" mapstructure:"show_caller,omitempty"`
	ShowFields  bool   `json:"show_stack,omitempty" yaml:"show_stack,omitempty" mapstructure:"show_stack,omitempty"`
}

// LogzFormatOptions represents the format options for the logger.
type LogzFormatOptions struct {
	Output   io.Writer `json:"output,omitempty" yaml:"output,omitempty" mapstructure:"output,omitempty"`
	MinLevel Level     `json:"min_level,omitempty" yaml:"min_level,omitempty" mapstructure:"min_level,omitempty"`
	MaxLevel Level     `json:"max_level,omitempty" yaml:"max_level,omitempty" mapstructure:"max_level,omitempty"`
	Level    Level     `json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
	Format   string    `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
}

// LogzOutputOptions represents the output options for the logger.
type LogzOutputOptions struct {
	// Output options
	OutputTTY    *bool   `json:"output_tty,omitempty" yaml:"output_tty,omitempty" mapstructure:"output_tty,omitempty"`
	OutputFile   *string `json:"output_file,omitempty" yaml:"output_file,omitempty" mapstructure:"output_file,omitempty"`
	OutputSyslog *string `json:"output_syslog,omitempty" yaml:"output_syslog,omitempty" mapstructure:"output_syslog,omitempty"`

	// Add any additional options here
	StackTrace *bool `json:"stack_trace,omitempty" yaml:"stack_trace,omitempty" mapstructure:"stack_trace,omitempty"`
}

// LogzRotatingOptions represents the rotating options for the logger.
type LogzRotatingOptions struct {
	// Rotation
	Rotate        *bool  `json:"rotate,omitempty" yaml:"rotate,omitempty" mapstructure:"rotate,omitempty"`
	RotateMaxSize *int64 `json:"rotate_max_size,omitempty" yaml:"rotate_max_size,omitempty" mapstructure:"rotate_max_size,omitempty"`
	RotateMaxBack *int64 `json:"rotate_max_back,omitempty" yaml:"rotate_max_back,omitempty" mapstructure:"rotate_max_back,omitempty"`
	RotateMaxAge  *int64 `json:"rotate_max_age,omitempty" yaml:"rotate_max_age,omitempty" mapstructure:"rotate_max_age,omitempty"`
	Compress      *bool  `json:"compress,omitempty" yaml:"compress,omitempty" mapstructure:"compress,omitempty"`
}

// LogzBufferingOptions represents the buffering options for the logger.
type LogzBufferingOptions struct {
	// Buffering
	Buffer        []byte         `json:"buffer,omitempty" yaml:"buffer,omitempty" mapstructure:"buffer,omitempty"`
	BufferSize    *int           `json:"buffer_size,omitempty" yaml:"buffer_size,omitempty" mapstructure:"buffer_size,omitempty"`
	FlushInterval *time.Duration `json:"flush_interval,omitempty" yaml:"flush_interval,omitempty" mapstructure:"flush_interval,omitempty"`
}

// LogzConfig is an alias for InitArgs.
type LogzConfig = InitArgs

// NewConfig returns a pointer to the LoggerArgs.
func NewConfig() *LogzConfig { return LoggerArgs }
