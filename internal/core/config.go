package core

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/formatter"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

const (
	LevelNotice  kbx.Level = "notice"
	LevelDebug   kbx.Level = "debug"
	LevelTrace   kbx.Level = "trace"
	LevelSuccess kbx.Level = "success"
	LevelInfo    kbx.Level = "info"
	LevelWarn    kbx.Level = "warn"
	LevelError   kbx.Level = "error"
	LevelFatal   kbx.Level = "fatal"
	LevelSilent  kbx.Level = "silent"
)

type LogzAdvancedOptions struct {
	// Hooks
	Formatter formatter.Formatter   `json:"formatter,omitempty" yaml:"formatter,omitempty" mapstructure:"formatter,omitempty"`
	Hooks     []interfaces.Hook     `json:"hooks,omitempty" yaml:"hooks,omitempty" mapstructure:"hooks,omitempty"`
	LHooks    interfaces.LHook[any] `json:"l_hooks,omitempty" yaml:"l_hooks,omitempty" mapstructure:"l_hooks,omitempty"`
	Metadata  map[string]any        `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

type LoggerConfig = kbx.InitArgs

type LoggerOptionsImpl struct {
	*LoggerConfig        `json:",inline" yaml:",inline" mapstructure:",squash"`
	*LogzAdvancedOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

func NewLoggerOptions(initArgs *kbx.InitArgs) *LoggerOptionsImpl {
	if initArgs != nil {
		return &LoggerOptionsImpl{
			LoggerConfig:        initArgs,
			LogzAdvancedOptions: &LogzAdvancedOptions{},
		}
	}
	return &LoggerOptionsImpl{
		LoggerConfig: &LoggerConfig{
			ID:                   uuid.New(),
			LogzGeneralOptions:   &kbx.LogzGeneralOptions{},
			LogzFormatOptions:    &kbx.LogzFormatOptions{},
			LogzOutputOptions:    &kbx.LogzOutputOptions{},
			LogzRotatingOptions:  &kbx.LogzRotatingOptions{},
			LogzBufferingOptions: &kbx.LogzBufferingOptions{},
		},
		LogzAdvancedOptions: &LogzAdvancedOptions{},
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
		&kbx.InitArgs{
			ID:                   o.ID,
			LogzGeneralOptions:   o.LogzGeneralOptions,
			LogzFormatOptions:    o.LogzFormatOptions,
			LogzOutputOptions:    o.LogzOutputOptions,
			LogzRotatingOptions:  o.LogzRotatingOptions,
			LogzBufferingOptions: o.LogzBufferingOptions,
		},
		o.LogzAdvancedOptions,
	}
}

// ConfigSetter aplica uma opção no logger.
type ConfigSetter func(l *Logger, value any)

var loggerSetters = map[string]ConfigSetter{}

func RegisterOptionSetter(key string, setter ConfigSetter) {
	loggerSetters[key] = setter
}

func init() {

	// ---- Níveis / Formatação / Writer ----

	RegisterOptionSetter("min_level", func(l *Logger, v any) {
		if lvl, ok := v.(kbx.Level); ok {
			l.SetMinLevel(lvl)
		}
	})

	RegisterOptionSetter("formatter", func(l *Logger, v any) {
		if f, ok := v.(formatter.Formatter); ok {
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
		o.MinLevel = value.(kbx.Level)

	case "formatter":
		o.Formatter = value.(formatter.Formatter)

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
