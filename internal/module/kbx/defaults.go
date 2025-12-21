// Package kbx has default configuration values
package kbx

const (
	DefaultLogLevel    = "info"
	DefaultLogMinLevel = "info"
	DefaultLogMaxLevel = "fatal"
	DefaultLogOutput   = "stdout"
	DefaultLogFormat   = "text"

	DefaultLogTimezone     = "UTC"
	DefaultLogLocale       = "en_US"
	DefaultLogDateFormat   = "2006-01-02"
	DefaultTimestampFormat = "2006-01-02 15:04:05"

	DefaultDebugMode   = false
	DefaultShowColor   = true
	DefaultShowIcons   = true
	DefaultShowCaller  = false
	DefaultShowTraceID = false
	DefaultShowFields  = false
	DefaultShowStack   = false

	DefaultMaxLogFileSize   = 10 // in MB
	DefaultMaxBackups       = 5
	DefaultMaxAge           = 30 // in days
	DefaultCompressLogFiles = true
	DefaultLogRotationTime  = "24h"
)

const (
	KeyringService        = "kubex"
	DefaultKubexConfigDir = "$HOME/.kubex"

	DefaultConfigDir         = "$HOME/.kubex/logz/config"
	DefaultConfigFile        = "$HOME/.kubex/logz/config.json"
	DefaultKubexDSConfigPath = "$HOME/.kubex/logz/config/config.json"
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
