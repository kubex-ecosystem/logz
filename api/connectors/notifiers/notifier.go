// Package notifiers provides the interface for notifier implementations.
package notifiers

import il "github.com/rafa-mori/logz/internal/core"

type LogzNotifierManager = il.NotifierManager
type LogzNotifier = il.Notifier

// type LogzNotifierType = il.NotifierImpl
// type LogzNotifierManagerImpl = il.NotifierManagerImpl

// NewLogzNotifierManager creates a new instance of LogzNotifierManagerImpl.
func NewLogzNotifierManager(notifiers map[string]LogzNotifier) LogzNotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]LogzNotifier)
	}
	return il.NewNotifierManager(notifiers)
}

//	type LogzNotifierManager interface {
//		// WebServer returns the HTTP server instance.
//		WebServer() *http.Server
//		// Websocket returns the Gorilla WebSocket connection instance.
//		Websocket() *websocket.Conn
//		// WebClient returns the HTTP client instance.
//		WebClient() *http.Client
//		// DBusClient returns the DBus connection instance.
//		DBusClient() *dbus.Conn
//		// AddNotifier adds or updates a notifier with the given name.
//		AddNotifier(name string, notifier Notifier)
//		// RemoveNotifier removes the notifier with the given name.
//		RemoveNotifier(name string)
//		// GetNotifier retrieves the notifier with the given name.
//		GetNotifier(name string) (Notifier, bool)
//		// ListNotifiers lists all registered notifier names.
//		ListNotifiers() []string
//		// UpdateFromConfig updates notifiers dynamically based on the provided configuration.
//		UpdateFromConfig() error
//	}
