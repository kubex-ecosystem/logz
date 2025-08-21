package core

import (
	"testing"
)

// TestLoggerCreation testa criação do logger
func TestLoggerCreation(t *testing.T) {
	logger := NewLogger("test")

	if logger == nil {
		t.Fatal("Failed to create logger")
	}

	t.Log("Logger created successfully")
}

// TestLoggerBasicLogging testa logging básico
func TestLoggerBasicLogging(t *testing.T) {
	logger := NewLogger("basic-test")

	// Testar diferentes níveis
	logger.InfoCtx("Info message", map[string]interface{}{"test": true})
	logger.WarnCtx("Warning message", map[string]interface{}{"test": true})
	logger.ErrorCtx("Error message", map[string]interface{}{"test": true})

	t.Log("Basic logging test completed")
}

// TestMetadataSet testa metadata
func TestMetadataSet(t *testing.T) {
	logger := NewLogger("metadata-test")

	logger.SetMetadata("service", "logz")
	logger.InfoCtx("Message with global metadata", map[string]interface{}{
		"local": "data",
	})

	t.Log("Metadata test completed")
}
