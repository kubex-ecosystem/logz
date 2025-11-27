// Package kbx provides utilities for working with initialization arguments.
package kbx

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
	Level    string
	MinLevel string
	MaxLevel string
	Output   string
	Format   string
	Message  []string

	ConfigFile     string
	ConfigType     string
	EnvFile        string
	LogFile        string
	Name           string
	Image          string
	Debug          bool
	ReleaseMode    bool
	IsConfidential bool
	DryRun         bool
	Force          bool
	Reset          bool
	FailFast       bool
	BatchMode      bool
	NoColor        bool
	RootMode       bool
	Host           string
	Command        string
	Subcommand     string
	MaxProcs       int
	TimeoutMS      int
	EnvVars        map[string]string
	Ports          map[string]string
	Volumes        map[string]string
	Bind           string
	Address        string
	PubCertKeyPath string
	PubKeyPath     string
	Pwd            string

	Enabled  bool
	Port     string
	Hostname string
	Username string
	Password string
	Database string
	DSN      string
}

func Log(level string, payload ...any) {
	gl.Log(level, payload...)
}

type MigrationInfo struct {
	// MigrationPath é o caminho para os arquivos de migração.
	MigrationPath string `json:"migration_path,omitempty" yaml:"migration_path,omitempty" mapstructure:"migration_path,omitempty"`
	// Enabled indica se a migração está habilitada.
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"false"`
	// Options são opções adicionais para a migração.
	Options map[string]any `json:"options,omitempty" yaml:"options,omitempty" mapstructure:"options,omitempty"`
	// Auto indica se a migração deve ser executada automaticamente.
	Auto *bool `json:"auto,omitempty" yaml:"auto,omitempty" mapstructure:"auto,omitempty" default:"false"`
	// Version indica a versão da migração.
	Version string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`
	// Bootstrap indica se a migração deve ser inicializada.
	Bootstrap *bool `json:"bootstrap,omitempty" yaml:"bootstrap,omitempty" mapstructure:"bootstrap,omitempty" default:"false"`
	// DryRun indica se a migração deve ser executada em modo de simulação.
	DryRun *bool `json:"dry_run,omitempty" yaml:"dry_run,omitempty" mapstructure:"dry_run,omitempty" default:"false"`
	// Reset indica se a migração deve redefinir o banco de dados antes de aplicar as migrações.
	Reset *bool `json:"reset,omitempty" yaml:"reset,omitempty" mapstructure:"reset,omitempty" default:"false"`
	// Force indica se a migração deve forçar a aplicação de todas as migrações.
	Force *bool `json:"force,omitempty" yaml:"force,omitempty" mapstructure:"force,omitempty" default:"false"`
}

// DBConfig is the pure, serializable configuration for a single database.
type DBConfig struct {
	ID         string         `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`
	Name       string         `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	IsDefault  bool           `json:"is_default,omitempty" yaml:"is_default,omitempty" mapstructure:"is_default,omitempty"`
	Enabled    *bool          `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
	Debug      bool           `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	Type       DBType         `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"Type,omitempty" validate:"required,oneof=postgres rabbitmq redis mongodb mysql mssql sqlite oracle"`
	Host       string         `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host,omitempty"`
	Port       string         `json:"port,omitempty" yaml:"port,omitempty" mapstructure:"port,omitempty"`
	User       string         `json:"user,omitempty" yaml:"user,omitempty" mapstructure:"user,omitempty"`
	Pass       string         `json:"pass,omitempty" yaml:"pass,omitempty" mapstructure:"pass,omitempty"`
	TLSEnabled bool           `json:"tls_enabled,omitempty" yaml:"tls_enabled,omitempty" mapstructure:"tls_enabled,omitempty"`
	DBName     string         `json:"db_name,omitempty" yaml:"db_name,omitempty" mapstructure:"db_name,omitempty"`
	Schema     string         `json:"schema,omitempty" yaml:"schema,omitempty" mapstructure:"schema,omitempty"`
	DSN        string         `json:"dsn,omitempty" yaml:"dsn,omitempty" mapstructure:"dsn,omitempty"`
	Options    map[string]any `json:"options,omitempty" yaml:"options,omitempty" mapstructure:"options,omitempty"`
	Backend    string         `json:"backend,omitempty" yaml:"backend,omitempty" mapstructure:"backend,omitempty"`
	Migration  *MigrationInfo `json:"migration,omitempty" yaml:"migration,omitempty" mapstructure:"migration,omitempty" validate:"omitempty,object"`
}

// RootConfig representa o arquivo de configuração do DS.
type RootConfig struct {
	Name      string      `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	FilePath  string      `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Enabled   *bool       `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
	Databases []*DBConfig `json:"databases,omitempty" yaml:"databases,omitempty" mapstructure:"databases,omitempty"`
}
