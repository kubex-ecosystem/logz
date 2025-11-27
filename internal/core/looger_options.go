package core

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// ConfigSetter aplica uma opção no logger.
type ConfigSetter func(l *Logger, value any)

var loggerSetters = map[string]ConfigSetter{}

func RegisterOptionSetter(key string, setter ConfigSetter) {
	loggerSetters[key] = setter
}

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

type LoggerOptionsImpl struct {
	ID     uuid.UUID `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`
	Prefix string    `json:"prefix,omitempty" yaml:"prefix,omitempty" mapstructure:"prefix,omitempty"`

	*LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzAdvancedOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

func NewLoggerOptions() *LoggerOptionsImpl {
	return &LoggerOptionsImpl{
		ID:                   uuid.New(),
		LogzGeneralOptions:   &LogzGeneralOptions{},
		LogzFormatOptions:    &LogzFormatOptions{},
		LogzOutputOptions:    &LogzOutputOptions{},
		LogzRotatingOptions:  &LogzRotatingOptions{},
		LogzBufferingOptions: &LogzBufferingOptions{},
		LogzAdvancedOptions:  &LogzAdvancedOptions{},
	}
}

func (o *LoggerOptionsImpl) ApplyOptions(logger *Logger) {

	for key, setter := range loggerSetters {
		if val := o.Get(key); val != nil {
			setter(logger, val)
		}
	}
}

func (o *LoggerOptionsImpl) Merge(override *LoggerOptionsImpl) *LoggerOptionsImpl {
	out := o.Clone()

	for key := range loggerSetters {
		if val := override.Get(key); val != nil {
			out.Set(key, val)
		}
	}

	return out
}

func (o *LoggerOptionsImpl) WithDefaults(def *LoggerOptionsImpl) *LoggerOptionsImpl {
	out := o.Clone()

	for key := range loggerSetters {
		if out.Get(key) == nil {
			out.Set(key, def.Get(key))
		}
	}

	return out
}

func (o *LoggerOptionsImpl) Inherit(parent *LoggerOptionsImpl) *LoggerOptionsImpl {
	return parent.WithDefaults(o)
}

func (o *LoggerOptionsImpl) Hydrate(prefix string) *LoggerOptionsImpl {
	out := o.Clone()

	for key := range loggerSetters {
		envKey := strings.ToUpper(prefix + "_" + key)
		if val := LoadFromEnvTyped(envKey, out.Get(key)); val != nil {
			out.Set(key, val)
		}
	}

	return out
}

func (o *LoggerOptionsImpl) Clone() *LoggerOptionsImpl {
	return &LoggerOptionsImpl{
		ID:                   o.ID,
		LogzGeneralOptions:   o.LogzGeneralOptions,
		LogzFormatOptions:    o.LogzFormatOptions,
		LogzOutputOptions:    o.LogzOutputOptions,
		LogzRotatingOptions:  o.LogzRotatingOptions,
		LogzBufferingOptions: o.LogzBufferingOptions,
		LogzAdvancedOptions:  o.LogzAdvancedOptions,
	}
}

func init() {

	// ---- Níveis / Formatação / Writer ----

	RegisterOptionSetter("min_level", func(l *Logger, v any) {
		if lvl, ok := v.(interfaces.Level); ok {
			l.SetMinLevel(lvl)
		}
	})

	RegisterOptionSetter("formatter", func(l *Logger, v any) {
		if f, ok := v.(interfaces.Formatter); ok {
			l.SetFormatter(f)
		}
	})

	RegisterOptionSetter("output", func(l *Logger, v any) {
		if w, ok := v.(io.Writer); ok {
			l.SetOutput(w)
		}
	})

	// ---- Rotação ----

	RegisterOptionSetter("rotate", func(l *Logger, v any) {
		if b, ok := v.(bool); ok {
			l.SetRotate(b)
		}
	})

	RegisterOptionSetter("rotate_max_size", func(l *Logger, v any) {
		if n, ok := v.(int64); ok {
			l.SetRotateMaxSize(n)
		}
	})

	RegisterOptionSetter("rotate_max_back", func(l *Logger, v any) {
		if n, ok := v.(int64); ok {
			l.SetRotateMaxBack(n)
		}
	})

	RegisterOptionSetter("rotate_max_age", func(l *Logger, v any) {
		if n, ok := v.(int64); ok {
			l.SetRotateMaxAge(n)
		}
	})

	RegisterOptionSetter("compress", func(l *Logger, v any) {
		if b, ok := v.(bool); ok {
			l.SetCompress(b)
		}
	})

	// ---- Buffer / Flush ----

	RegisterOptionSetter("buffer_size", func(l *Logger, v any) {
		if n, ok := v.(int); ok {
			l.SetBufferSize(n)
		}
	})

	RegisterOptionSetter("flush_interval", func(l *Logger, v any) {
		if d, ok := v.(time.Duration); ok {
			l.SetFlushInterval(d)
		}
	})

	// ---- Hooks ----

	RegisterOptionSetter("hooks", func(l *Logger, v any) {
		if h, ok := v.([]interfaces.Hook); ok {
			l.SetHooks(h)
		}
	})

	RegisterOptionSetter("lhooks", func(l *Logger, v any) {
		if h, ok := v.(interfaces.LHook[any]); ok {
			l.SetLHooks(h)
		}
	})

	// ---- Metadata ----

	RegisterOptionSetter("metadata", func(l *Logger, v any) {
		if m, ok := v.(map[string]any); ok {
			l.SetMetadata(m)
		}
	})
}

func (o *LoggerOptionsImpl) Get(key string) any {
	switch key {
	case "min_level":
		return o.MinLevel
	case "formatter":
		return o.Formatter
	case "output":
		return o.Output

	case "rotate":
		return deref(o.Rotate)
	case "rotate_max_size":
		return deref(o.RotateMaxSize)
	case "rotate_max_back":
		return deref(o.RotateMaxBack)
	case "rotate_max_age":
		return deref(o.RotateMaxAge)
	case "compress":
		return deref(o.Compress)

	case "buffer_size":
		return deref(o.BufferSize)
	case "flush_interval":
		return derefDuration(o.FlushInterval)

	case "hooks":
		return o.Hooks
	case "lhooks":
		return o.LHooks
	case "metadata":
		return o.Metadata
	}
	return nil
}
func (o *LoggerOptionsImpl) Set(key string, value any) {
	switch key {
	case "min_level":
		o.MinLevel = value.(interfaces.Level)

	case "formatter":
		o.Formatter = value.(interfaces.Formatter)

	case "output":
		o.Output = value.(io.Writer)

	case "rotate":
		o.Rotate = kbx.BoolPtr(value.(bool))
	case "rotate_max_size":
		o.RotateMaxSize = kbx.PtrInt64(value.(int64))
	case "rotate_max_back":
		o.RotateMaxBack = kbx.PtrInt64(value.(int64))
	case "rotate_max_age":
		o.RotateMaxAge = kbx.PtrInt64(value.(int64))
	case "compress":
		o.Compress = kbx.BoolPtr(value.(bool))

	case "buffer_size":
		o.BufferSize = kbx.PtrInt(value.(int64))
	case "flush_interval":
		o.FlushInterval = kbx.PtrDuration(value.(time.Duration))

	case "hooks":
		o.Hooks = value.([]interfaces.Hook)
	case "lhooks":
		o.LHooks = value.(interfaces.LHook[any])
	case "metadata":
		o.Metadata = value.(map[string]any)
	}
}

func LoadFromEnvTyped(key string, defaultValue any) any {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	switch defaultValue.(type) {
	case int:
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	case bool:
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	case time.Duration:
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}

	return val
}

func deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}
func derefDuration(ptr *time.Duration) time.Duration {
	if ptr != nil {
		return *ptr
	}
	return 0
}
