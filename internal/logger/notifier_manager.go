package logger

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

// NotifierManager defines the interface for managing notifiers.
type NotifierManager interface {
	// WebServer returns the HTTP server instance.
	WebServer() *http.Server
	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	// Websocket returns the ZMQ socket instance.
	// Websocket() *zmq4.Socket

	// WebClient returns the HTTP client instance.
	WebClient() *http.Client
	// DBusClient returns the DBus connection instance.
	DBusClient() *dbus.Conn

	// AddNotifier adds or updates a notifier with the given name.
	AddNotifier(name string, notifier Notifier)
	// RemoveNotifier removes the notifier with the given name.
	RemoveNotifier(name string)
	// GetNotifier retrieves the notifier with the given name.
	GetNotifier(name string) (Notifier, bool)
	// ListNotifiers lists all registered notifier names.
	ListNotifiers() []string

	// UpdateFromConfig updates notifiers dynamically based on the provided configuration.
	UpdateFromConfig() error
}

// NotifierManagerImpl is the implementation of the NotifierManager interface.
type NotifierManagerImpl struct {
	webServer *http.Server
	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	// websocket  *zmq4.Socket
	webClient  *http.Client
	dbusClient *dbus.Conn
	notifiers  map[string]Notifier
	mu         sync.RWMutex
}

// NewNotifierManager creates a new instance of NotifierManagerImpl.
func NewNotifierManager(notifiers map[string]Notifier) NotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]Notifier)
	}
	return &NotifierManagerImpl{
		notifiers: notifiers,
	}
}

// AddNotifier adds or updates a notifier with the given name.
func (nm *NotifierManagerImpl) AddNotifier(name string, notifier Notifier) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.notifiers[name] = notifier
	fmt.Printf("Notifier '%s' added/updated.\n", name)
}

// RemoveNotifier removes the notifier with the given name.
func (nm *NotifierManagerImpl) RemoveNotifier(name string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	delete(nm.notifiers, name)
	fmt.Printf("Notifier '%s' removed.\n", name)
}

// GetNotifier retrieves the notifier with the given name.
func (nm *NotifierManagerImpl) GetNotifier(name string) (Notifier, bool) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	notifier, ok := nm.notifiers[name]
	return notifier, ok
}

// ListNotifiers lists all registered notifier names.
func (nm *NotifierManagerImpl) ListNotifiers() []string {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	keys := make([]string, 0, len(nm.notifiers))
	for name := range nm.notifiers {
		keys = append(keys, name)
	}
	return keys
}

// UpdateFromConfig updates notifiers dynamically based on the provided configuration.
func (nm *NotifierManagerImpl) UpdateFromConfig() error {
	var configNotifiers map[string]map[string]interface{}
	if err := viper.UnmarshalKey("notifiers", &configNotifiers); err != nil {
		return fmt.Errorf("failed to parse notifiers config: %w", err)
	}

	// Update or recreate notifiers dynamically
	for name, conf := range configNotifiers {
		typ, ok := conf["type"].(string)
		if !ok {
			fmt.Printf("Notifier '%s' does not specify a type and will be ignored.\n", name)
			continue
		}

		switch typ {
		case "http":
			webhookURL, _ := conf["webhookURL"].(string)
			authToken, _ := conf["authToken"].(string)
			notifier := NewHTTPNotifier(webhookURL, authToken)
			nm.AddNotifier(name, notifier)
		case "zmq":
			endpoint, _ := conf["endpoint"].(string)
			notifier := NewZMQNotifier(endpoint)
			nm.AddNotifier(name, notifier)
		case "dbus":
			notifier := NewDBusNotifier()
			nm.AddNotifier(name, notifier)
		default:
			fmt.Printf("Unknown notifier type '%s' for notifier '%s'.\n", typ, name)
		}
	}
	return nil
}

// WebServer returns the HTTP server instance.
func (nm *NotifierManagerImpl) WebServer() *http.Server {
	if nm.webServer == nil {
		nm.webServer = Server()
	}
	return nm.webServer
}

// Temporarily disabled due to external dependency on zmq4
// Uncomment and ensure the required libraries are installed if needed in the future
// Websocket returns the ZMQ socket instance.
//func (nm *NotifierManagerImpl) Websocket() *zmq4.Socket {
//	if nm.websocket == nil {
//		nm.websocket = Socket()
//	}
//	return nm.websocket
//}

// WebClient returns the HTTP client instance.
func (nm *NotifierManagerImpl) WebClient() *http.Client {
	if nm.webClient == nil {
		nm.webClient = Client()
	}
	return nm.webClient
}

// DBusClient returns the DBus connection instance.
func (nm *NotifierManagerImpl) DBusClient() *dbus.Conn {
	if nm.dbusClient == nil {
		nm.dbusClient = DBus()
	}
	return nm.dbusClient
}
