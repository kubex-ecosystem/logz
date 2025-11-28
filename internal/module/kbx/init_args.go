// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/kubex-ecosystem/logz/internal/writer"
)

type glgr interface {
	Log(level string, parts ...any)
}

var gl glgr

func SetLogger(logger glgr) {
	gl = logger
}

type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeRabbitMQ DBType = "rabbitmq"
	DBTypeRedis    DBType = "redis"
	DBTypeMongoDB  DBType = "mongodb"
	DBTypeMySQL    DBType = "mysql"
	DBTypeMSSQL    DBType = "mssql"
	DBTypeSQLite   DBType = "sqlite"
	DBTypeOracle   DBType = "oracle"
)

type InitArgs struct {
	ID uuid.UUID

	*LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	// *LogzAdvancedOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// RootConfig representa o arquivo de configuração do DS.
type RootConfig struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	FilePath string `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
}

var LoggerArgs *InitArgs = &InitArgs{
	ID:                   uuid.New(),
	LogzGeneralOptions:   &LogzGeneralOptions{},
	LogzFormatOptions:    &LogzFormatOptions{},
	LogzOutputOptions:    &LogzOutputOptions{},
	LogzRotatingOptions:  &LogzRotatingOptions{},
	LogzBufferingOptions: &LogzBufferingOptions{},
}

func ParseLoggerArgs(level string, minLevel string, maxLevel string, output string) *InitArgs {
	LoggerArgs.Level = Level(GetValueOrDefaultSimple(level, "info"))
	LoggerArgs.MinLevel = Level(GetValueOrDefaultSimple(minLevel, "debug"))
	LoggerArgs.MaxLevel = Level(GetValueOrDefaultSimple(maxLevel, "fatal"))
	LoggerArgs.Output = GetValueOrDefaultSimple[io.Writer](writer.ParseWriter(output), os.Stdout)
	return LoggerArgs
}

func init() {
	if LoggerArgs == nil {
		LoggerArgs = &InitArgs{
			ID:                   uuid.New(),
			LogzGeneralOptions:   &LogzGeneralOptions{},
			LogzFormatOptions:    &LogzFormatOptions{},
			LogzOutputOptions:    &LogzOutputOptions{},
			LogzRotatingOptions:  &LogzRotatingOptions{},
			LogzBufferingOptions: &LogzBufferingOptions{},
		}
	}
}
