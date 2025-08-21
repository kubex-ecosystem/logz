package core

//func TestIntegrationLogger(t *testing.T) {
//	configManager := NewConfigManager()
//	if configManager == nil {
//		t.Fatal("ErrorCtx initializing ConfigManager.")
//	}
//	cfgMgr := *configManager
//	config, err := cfgMgr.LoadConfig()
//	if err != nil {
//		t.Fatalf("ErrorCtx loading configuration: %v", err)
//	}
//	logger := NewLogger(config)
//	if logger == nil {
//		t.Fatal("ErrorCtx creating core.")
//	}
//
//	var buf bytes.Buffer
//	logger.SetWriter(NewDefaultWriter(&buf, &TextFormatter{}))
//
//	logger.DebugCtx("debug message", nil)
//	if !bytes.Contains(buf.Bytes(), []byte("DEBUG")) {
//		t.Errorf("Expected 'DEBUG' log entry, got '%s'", buf.String())
//	}
//
//	logger.InfoCtx("info message", nil)
//	if !bytes.Contains(buf.Bytes(), []byte("INFO")) {
//		t.Errorf("Expected 'INFO' log entry, got '%s'", buf.String())
//	}
//
//	logger.WarnCtx("warn message", nil)
//	if !bytes.Contains(buf.Bytes(), []byte("WARN")) {
//		t.Errorf("Expected 'WARN' log entry, got '%s'", buf.String())
//	}
//
//	logger.ErrorCtx("error message", nil)
//	if !bytes.Contains(buf.Bytes(), []byte("ERROR")) {
//		t.Errorf("Expected 'ERROR' log entry, got '%s'", buf.String())
//	}
//}
//
//func TestIntegrationConcurrentAccess(t *testing.T) {
//	configManager := NewConfigManager()
//	if configManager == nil {
//		t.Fatal("ErrorCtx initializing ConfigManager.")
//	}
//	cfgMgr := *configManager
//	config, err := cfgMgr.LoadConfig()
//	if err != nil {
//		t.Fatalf("ErrorCtx loading configuration: %v", err)
//	}
//	logger := NewLogger(config)
//	if logger == nil {
//		t.Fatal("ErrorCtx creating core.")
//	}
//
//	var wg sync.WaitGroup
//	for i := 0; i < 100; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			logger.SetMetadata("key", i)
//			logger.DebugCtx("concurrent message", nil)
//		}(i)
//	}
//	wg.Wait()
//}
//
//func TestIntegrationLogRotation(t *testing.T) {
//	configManager := NewConfigManager()
//	if configManager == nil {
//		t.Fatal("ErrorCtx initializing ConfigManager.")
//	}
//	cfgMgr := *configManager
//	config, err := cfgMgr.LoadConfig()
//	if err != nil {
//		t.Fatalf("ErrorCtx loading configuration: %v", err)
//	}
//	logger := NewLogger(config)
//	if logger == nil {
//		t.Fatal("ErrorCtx creating core.")
//	}
//
//	logFile := config.Output()
//	defer os.Remove(logFile)
//
//	for i := 0; i < 1000; i++ {
//		logger.InfoCtx("log rotation test message", nil)
//	}
//
//	err = CheckLogSize(config)
//	if err != nil {
//		t.Fatalf("ErrorCtx checking log size: %v", err)
//	}
//
//	_, err = os.Stat(logFile + ".tar.gz")
//	if os.IsNotExist(err) {
//		t.Errorf("Expected log rotation to create archive, but it did not.")
//	}
//}
//
//func TestIntegrationDynamicConfigChanges(t *testing.T) {
//	configManager := NewConfigManager()
//	if configManager == nil {
//		t.Fatal("ErrorCtx initializing ConfigManager.")
//	}
//	cfgMgr := *configManager
//	config, err := cfgMgr.LoadConfig()
//	if err != nil {
//		t.Fatalf("ErrorCtx loading configuration: %v", err)
//	}
//	logger := NewLogger(config)
//	if logger == nil {
//		t.Fatal("ErrorCtx creating core.")
//	}
//
//	viper.Set("logLevel", "DEBUG")
//	time.Sleep(1 * time.Second)
//
//	if logger.GetLevel() != DEBUG {
//		t.Errorf("Expected log VLevel to be 'DEBUG', got '%v'", logger.GetLevel())
//	}
//}
//
//func TestIntegrationNotifierManager(t *testing.T) {
//	configManager := NewConfigManager()
//	if configManager == nil {
//		t.Fatal("ErrorCtx initializing ConfigManager.")
//	}
//	cfgMgr := *configManager
//	config, err := cfgMgr.LoadConfig()
//	if err != nil {
//		t.Fatalf("ErrorCtx loading configuration: %v", err)
//	}
//	logger := NewLogger(config)
//	if logger == nil {
//		t.Fatal("ErrorCtx creating core.")
//	}
//
//	notifierManager := config.NotifierManager()
//	if notifierManager == nil {
//		t.Fatal("ErrorCtx initializing NotifierManager.")
//	}
//
//	notifierManager.AddNotifier("testNotifier", NewHTTPNotifier("http://example.com", ""))
//	if _, ok := notifierManager.GetNotifier("testNotifier"); !ok {
//		t.Errorf("Expected 'testNotifier' to be added, but it was not.")
//	}
//
//	notifierManager.RemoveNotifier("testNotifier")
//	if _, ok := notifierManager.GetNotifier("testNotifier"); ok {
//		t.Errorf("Expected 'testNotifier' to be removed, but it was not.")
//	}
//}
//
//func TestIntegrationService(t *testing.T) {
//	go func() {
//		if err := Run(); err != nil {
//			t.Fatalf("ErrorCtx running service: %v", err)
//		}
//	}()
//	time.Sleep(1 * time.Second)
//
//	resp, err := http.Get("http://localhost:9999/health")
//	if err != nil {
//		t.Fatalf("ErrorCtx making HTTP request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
//	}
//
//	if err := Stop(); err != nil {
//		t.Fatalf("ErrorCtx stopping service: %v", err)
//	}
//}
