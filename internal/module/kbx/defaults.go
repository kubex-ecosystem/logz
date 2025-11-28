// Package kbx has default configuration values
package kbx

const (
	KeyringService        = "canalize"
	DefaultKubexConfigDir = "$HOME/.canalize"

	DefaultCanalizeBEKeyPath    = "$HOME/.canalize/gobe/gobe-key.pem"
	DefaultCanalizeBECertPath   = "$HOME/.canalize/gobe/gobe-cert.pem"
	DefaultCanalizeBECAPath     = "$HOME/.canalize/gobe/ca-cert.pem"
	DefaultCanalizeBEConfigPath = "$HOME/.canalize/gobe/config/config.json"

	DefaultConfigDir            = "$HOME/.canalize/logz/config"
	DefaultConfigFile           = "$HOME/.canalize/logz/config.json"
	DefaultCanalizeDSConfigPath = "$HOME/.canalize/logz/config/config.json"
)

const (
	DefaultVolumesDir     = "$HOME/.canalize/volumes"
	DefaultRedisVolume    = "$HOME/.canalize/volumes/redis"
	DefaultPostgresVolume = "$HOME/.canalize/volumes/postgresql"
	DefaultMongoDBVolume  = "$HOME/.canalize/volumes/mongodb"
	DefaultMongoVolume    = "$HOME/.canalize/volumes/mongo"
	DefaultRabbitMQVolume = "$HOME/.canalize/volumes/rabbitmq"
)

const (
	DefaultRateLimitLimit  = 100
	DefaultRateLimitBurst  = 100
	DefaultRequestWindow   = 1 * 60 * 1000 // 1 minute
	DefaultRateLimitJitter = 0.1
)

const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = 1 * 1000 // 1 second
)

const (
	DefaultMaxIdleConns          = 100
	DefaultMaxIdleConnsPerHost   = 100
	DefaultIdleConnTimeout       = 90 * 1000 // 90 seconds
	DefaultTLSHandshakeTimeout   = 10 * 1000 // 10 seconds
	DefaultExpectContinueTimeout = 1 * 1000  // 1 second
	DefaultResponseHeaderTimeout = 5 * 1000  // 5 seconds
	DefaultTimeout               = 30 * 1000 // 30 seconds
	DefaultKeepAlive             = 30 * 1000 // 30 seconds
	DefaultMaxConnsPerHost       = 100
)

const (
	DefaultLLMProvider    = "gemini"
	DefaultLLMModel       = "gemini-2.0-flash"
	DefaultLLMMaxTokens   = 1024
	DefaultLLMTemperature = 0.3
)

const (
	DefaultApprovalRequireForResponses = false
	DefaultApprovalTimeoutMinutes      = 15
)

const (
	DefaultServerPort = "5000"
	DefaultServerHost = "0.0.0.0"
)

type ValidationError struct {
	Field   string
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}
func (v *ValidationError) FieldError() map[string]string {
	return map[string]string{v.Field: v.Message}
}
func (v *ValidationError) FieldsError() map[string]string {
	return map[string]string{v.Field: v.Message}
}
func (v *ValidationError) ErrorOrNil() error {
	return v
}

var (
	ErrUsernameRequired = &ValidationError{Field: "username", Message: "Username is required"}
	ErrPasswordRequired = &ValidationError{Field: "password", Message: "Password is required"}
	ErrEmailRequired    = &ValidationError{Field: "email", Message: "Email is required"}
	ErrDBNotProvided    = &ValidationError{Field: "db", Message: "Database not provided"}
	ErrModelNotFound    = &ValidationError{Field: "model", Message: "Model not found"}
)
