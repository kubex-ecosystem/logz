package logger

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}
}

func TestSetMetadata(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}

	logger.SetMetadata("key", "value")
	if logger.metadata["key"] != "value" {
		t.Errorf("Expected metadata 'key' to be 'value', got '%v'", logger.metadata["key"])
	}
}

func TestLogMethods(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}

	var buf bytes.Buffer
	logger.SetWriter(NewDefaultWriter(&buf, &TextFormatter{}))

	logger.Debug("debug message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("DEBUG")) {
		t.Errorf("Expected 'DEBUG' log entry, got '%s'", buf.String())
	}

	logger.Info("info message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("INFO")) {
		t.Errorf("Expected 'INFO' log entry, got '%s'", buf.String())
	}

	logger.Warn("warn message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("WARN")) {
		t.Errorf("Expected 'WARN' log entry, got '%s'", buf.String())
	}

	logger.Error("error message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("ERROR")) {
		t.Errorf("Expected 'ERROR' log entry, got '%s'", buf.String())
	}
}

func TestConcurrentAccess(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			logger.SetMetadata("key", i)
			logger.Debug("concurrent message", nil)
		}(i)
	}
	wg.Wait()
}

func TestLogRotation(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}

	logFile := config.Output()
	defer os.Remove(logFile)

	for i := 0; i < 1000; i++ {
		logger.Info("log rotation test message", nil)
	}

	err = CheckLogSize(config)
	if err != nil {
		t.Fatalf("Error checking log size: %v", err)
	}

	_, err = os.Stat(logFile + ".tar.gz")
	if os.IsNotExist(err) {
		t.Errorf("Expected log rotation to create archive, but it did not.")
	}
}

func TestDynamicConfigChanges(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("Error initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("Error creating logger.")
	}

	viper.Set("logLevel", "DEBUG")
	time.Sleep(1 * time.Second)

	if logger.GetLevel() != DEBUG {
		t.Errorf("Expected log level to be 'DEBUG', got '%v'", logger.GetLevel())
	}
}
