package logger

import (
	iKbxCfg "github.com/faelmori/kubex-interfaces/config"

	"fmt"
	"sync"
)

type LogInterface interface {
	Initialize(config iKbxCfg.ConfigManager[iKbxCfg.Configurable]) error
	Log(message string)
	Stop()
}

var (
	loggerRegistry = make(map[string]*LoggerInstance)
	registryMutex  sync.Mutex
)

type LoggerInstance struct {
	LogChannel  chan string
	DoneChannel chan bool
	Config      *iKbxCfg.ConfigManager[LogzConfig]
}

func (li *LoggerInstance) Start() {
	go func() {
		for {
			select {
			case log := <-li.LogChannel:
				fmt.Println("Log:", log)
			case <-li.DoneChannel:
				fmt.Println("Logger finalizado.")
				return
			}
		}
	}()
}

func GetLoggerInstance(name string, config *iKbxCfg.ConfigManager[LogzConfig]) (*LoggerInstance, error) {
	// Validate the configuration
	if config == nil {
		return nil, fmt.Errorf("configuração inválida")
	}

	// Necessary to validate the configuration because of the pointer receiver
	cfg := config.GetConfig()

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuração inválida: %v", err)
	}

	registryMutex.Lock()
	defer registryMutex.Unlock()

	if logger, exists := loggerRegistry[name]; exists {
		return logger, nil
	}

	newLogger := &LoggerInstance{
		LogChannel:  make(chan string, 100),
		DoneChannel: make(chan bool),
		Config:      config,
	}
	newLogger.Start()
	loggerRegistry[name] = newLogger
	return newLogger, nil
}
