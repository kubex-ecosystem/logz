package core

/*func TestNewLogger(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}
}

func TestSetMetadata(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}

	logger.SetMetadata("key", "value")
	if logger.metadata["key"] != "value" {
		t.Errorf("Expected VMetadata 'key' to be 'value', got '%v'", logger.metadata["key"])
	}
}

func TestLogMethods(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}

	var buf bytes.Buffer
	logger.SetWriter(NewDefaultWriter[any](&buf, &TextFormatter{}))

	logger.DebugCtx("debug message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("DEBUG")) {
		t.Errorf("Expected 'DEBUG' log entry, got '%s'", buf.String())
	}

	logger.InfoCtx("info message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("INFO")) {
		t.Errorf("Expected 'INFO' log entry, got '%s'", buf.String())
	}

	logger.WarnCtx("warn message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("WARN")) {
		t.Errorf("Expected 'WARN' log entry, got '%s'", buf.String())
	}

	logger.ErrorCtx("error message", nil)
	if !bytes.Contains(buf.Bytes(), []byte("ERROR")) {
		t.Errorf("Expected 'ERROR' log entry, got '%s'", buf.String())
	}
}

func TestConcurrentAccess(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			logger.SetMetadata("key", i)
			logger.DebugCtx("concurrent message", nil)
		}(i)
	}
	wg.Wait()
}

func TestLogRotation(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}

	logFile := config.Output()
	defer os.Remove(logFile)

	for i := 0; i < 1000; i++ {
		logger.InfoCtx("log rotation test message", nil)
	}

	err = CheckLogSize(config)
	if err != nil {
		t.Fatalf("ErrorCtx checking log size: %v", err)
	}

	_, err = os.Stat(logFile + ".tar.gz")
	if os.IsNotExist(err) {
		t.Errorf("Expected log rotation to create archive, but it did not.")
	}
}

func TestDynamicConfigChanges(t *testing.T) {
	configManager := NewConfigManager()
	if configManager == nil {
		t.Fatal("ErrorCtx initializing ConfigManager.")
	}
	cfgMgr := *configManager
	config, err := cfgMgr.LoadConfig()
	if err != nil {
		t.Fatalf("ErrorCtx loading configuration: %v", err)
	}
	logger := NewLogger(config)
	if logger == nil {
		t.Fatal("ErrorCtx creating core.")
	}

	viper.Set("logLevel", "DEBUG")
	time.Sleep(1 * time.Second)

	if logger.GetLevel() != DEBUG {
		t.Errorf("Expected log VLevel to be 'DEBUG', got '%v'", logger.GetLevel())
	}
}
*/
