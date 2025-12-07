package kbx

import (
	"io"
	"time"
)

type LogzGeneralOptions struct {
	// General options
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty" mapstructure:"prefix,omitempty"`
	Debug     bool   `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	ShowColor *bool  `json:"show_color,omitempty" yaml:"show_color,omitempty" mapstructure:"show_color,omitempty"`
	ShowIcons *bool  `json:"show_icons,omitempty" yaml:"show_icons,omitempty" mapstructure:"show_icons,omitempty"`
}

type LogzFormatOptions struct {
	Output   io.Writer `json:"output,omitempty" yaml:"output,omitempty" mapstructure:"output,omitempty"`
	MinLevel Level     `json:"min_level,omitempty" yaml:"min_level,omitempty" mapstructure:"min_level,omitempty"`
	MaxLevel Level     `json:"max_level,omitempty" yaml:"max_level,omitempty" mapstructure:"max_level,omitempty"`
	Level    Level     `json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
	Format   string    `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
}

type LogzOutputOptions struct {
	// Output options
	OutputTTY    *bool   `json:"output_tty,omitempty" yaml:"output_tty,omitempty" mapstructure:"output_tty,omitempty"`
	OutputFile   *string `json:"output_file,omitempty" yaml:"output_file,omitempty" mapstructure:"output_file,omitempty"`
	OutputSyslog *string `json:"output_syslog,omitempty" yaml:"output_syslog,omitempty" mapstructure:"output_syslog,omitempty"`

	// Add any additional options here
	StackTrace *bool `json:"stack_trace,omitempty" yaml:"stack_trace,omitempty" mapstructure:"stack_trace,omitempty"`
}

type LogzRotatingOptions struct {
	// Rotation
	Rotate        *bool  `json:"rotate,omitempty" yaml:"rotate,omitempty" mapstructure:"rotate,omitempty"`
	RotateMaxSize *int64 `json:"rotate_max_size,omitempty" yaml:"rotate_max_size,omitempty" mapstructure:"rotate_max_size,omitempty"`
	RotateMaxBack *int64 `json:"rotate_max_back,omitempty" yaml:"rotate_max_back,omitempty" mapstructure:"rotate_max_back,omitempty"`
	RotateMaxAge  *int64 `json:"rotate_max_age,omitempty" yaml:"rotate_max_age,omitempty" mapstructure:"rotate_max_age,omitempty"`
	Compress      *bool  `json:"compress,omitempty" yaml:"compress,omitempty" mapstructure:"compress,omitempty"`
}

type LogzBufferingOptions struct {
	// Buffering
	Buffer        []byte         `json:"buffer,omitempty" yaml:"buffer,omitempty" mapstructure:"buffer,omitempty"`
	BufferSize    *int           `json:"buffer_size,omitempty" yaml:"buffer_size,omitempty" mapstructure:"buffer_size,omitempty"`
	FlushInterval *time.Duration `json:"flush_interval,omitempty" yaml:"flush_interval,omitempty" mapstructure:"flush_interval,omitempty"`
}
