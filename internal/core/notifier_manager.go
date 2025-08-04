package core

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type NotifierWebSocketConfig struct {
	TLSClientConfig  *tls.Config
	HandshakeTimeout time.Duration
	Jar              http.CookieJar

	WriteBufferPool                 websocket.BufferPool
	ReadBufferSize, WriteBufferSize int
	Subprotocols                    []string

	EnableCompression bool

	NetDial           func(network, addr string) (net.Conn, error)
	NetDialContext    func(ctx context.Context, network, addr string) (net.Conn, error)
	NetDialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)
	Proxy             func(*http.Request) (*url.URL, error)
}

// NotifierManager defines the interface for managing notifiers.
type NotifierManager interface {
	// WebServer returns the HTTP server instance.
	WebServer() *http.Server

	// Websocket returns the Gorilla WebSocket connection instance.
	Websocket() *websocket.Conn

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
	webServer  *http.Server
	websocket  *websocket.Conn
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
		return fmt.Errorf("failed to parse notifiers VConfig: %w", err)
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
		case "websocket":
			endpoint, _ := conf["endpoint"].(string)
			notifier := NewWebSocketNotifier(endpoint)
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

// Websocket returns the Gorilla WebSocket connection instance.
// If the WebSocket connection is not initialized, it will create a new one.
// If the connection fails, it will return nil and log the error.
// This is useful for lazy initialization of the WebSocket connection.
func (nm *NotifierManagerImpl) Websocket() *websocket.Conn {
	if nm.websocket == nil {
		var err error
		dialer := websocket.Dialer{
			// Configure the WebSocket dialer if needed
		}
		nm.websocket, _, err = dialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			fmt.Printf("Failed to initialize WebSocket: %v\n", err)
			return nil
		}
	}
	return nm.websocket
}

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
