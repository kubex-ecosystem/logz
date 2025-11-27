package kbx

import (
	"io"
	"time"

	"github.com/kubex-ecosystem/logz/interfaces"
)

type LogzGeneralOptions struct {
	// General options
	Debug     *bool `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	ShowColor *bool `json:"show_color,omitempty" yaml:"show_color,omitempty" mapstructure:"show_color,omitempty"`
	ShowIcons *bool `json:"show_icons,omitempty" yaml:"show_icons,omitempty" mapstructure:"show_icons,omitempty"`
}

type LogzFormatOptions struct {
	// Format        string
	Formatter interfaces.Formatter `json:"formatter,omitempty" yaml:"formatter,omitempty" mapstructure:"formatter,omitempty"`
	Output    io.Writer            `json:"output,omitempty" yaml:"output,omitempty" mapstructure:"output,omitempty"`
	MinLevel  interfaces.Level     `json:"min_level,omitempty" yaml:"min_level,omitempty" mapstructure:"min_level,omitempty"`
	MaxLevel  interfaces.Level     `json:"max_level,omitempty" yaml:"max_level,omitempty" mapstructure:"max_level,omitempty"`
	Level     interfaces.Level     `json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
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

type LogzAdvancedOptions struct {
	// Hooks
	Hooks    []interfaces.Hook     `json:"hooks,omitempty" yaml:"hooks,omitempty" mapstructure:"hooks,omitempty"`
	LHooks   interfaces.LHook[any] `json:"l_hooks,omitempty" yaml:"l_hooks,omitempty" mapstructure:"l_hooks,omitempty"`
	Metadata map[string]any        `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}
