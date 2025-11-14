// Package interfaces defines interfaces for log notifiers.
package interfaces

import (
	"net/http"

	"github.com/godbus/dbus/v5"
	"github.com/gorilla/websocket"
)

// Notifier defines the interface for a log notifier.
type Notifier interface {
	// Notify sends a log entry notification.
	Notify(entry LogzEntry) error

	// Enable activates the notifier.
	Enable()

	// Disable deactivates the notifier.
	Disable()

	// Enabled checks if the notifier is active.
	Enabled() bool

	// WebServer returns the HTTP server instance.
	WebServer() *http.Server

	// Websocket returns the WebSocket instance.
	Websocket() *websocket.Conn

	// WebClient returns the HTTP client instance.
	WebClient() *http.Client
	// DBusClient returns the DBus connection instance.
	DBusClient() *dbus.Conn
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
