package core

import (
	"github.com/godbus/dbus/v5"

	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	// Websocket returns the WebSocket instance.
	//Websocket() *zmq4.Socket

	// WebClient returns the HTTP client instance.
	WebClient() *http.Client
	// DBusClient returns the DBus connection instance.
	DBusClient() *dbus.Conn
}

// NotifierImpl is the implementation of the Notifier interface.
type NotifierImpl struct {
	NotifierManager NotifierManager // Manager for notifier instances.
	EnabledFlag     bool            // Flag indicating if the notifier is enabled.
	WebhookURL      string          // URL for webhook notifications.
	HttpMethod      string          // HTTP method for webhook notifications.
	AuthToken       string          // Authentication token for notifications.
	LogLevel        string          // Log VLevel for filtering notifications.
	WsEndpoint      string          // WebSocket endpoint for notifications.
	Whitelist       []string        // Whitelist of sources for notifications.
}

// NewNotifier creates a new NotifierImpl instance.
func NewNotifier(manager NotifierManager, enabled bool, webhookURL, httpMethod, authToken, logLevel, wsEndpoint string, whitelist []string) Notifier {
	if whitelist == nil {
		whitelist = []string{}
	}
	return &NotifierImpl{
		NotifierManager: manager,
		EnabledFlag:     enabled,
		WebhookURL:      webhookURL,
		HttpMethod:      httpMethod,
		AuthToken:       authToken,
		LogLevel:        logLevel,
		WsEndpoint:      wsEndpoint,
		Whitelist:       whitelist,
	}
}

// Notify sends a log entry notification based on the configured settings.
func (n *NotifierImpl) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}

	// Validate log VLevel
	if n.LogLevel != "" && n.LogLevel != string(entry.GetLevel()) {
		return nil
	}

	// Validate Whitelist
	if len(n.Whitelist) > 0 && !contains(n.Whitelist, entry.GetSource()) {
		return nil
	}

	// HTTP Notification
	if n.WebhookURL != "" {
		if err := n.httpNotify(entry); err != nil {
			return err
		}
	}

	// WebSocket Notification
	if n.WsEndpoint != "" {
		if err := n.wsNotify(entry); err != nil {
			return err
		}
	}

	// DBus Notification
	if n.DBusClient() != nil {
		if err := n.dbusNotify(entry); err != nil {
			return err
		}
	}

	return nil
}

// httpNotify sends an HTTP notification.
func (n *NotifierImpl) httpNotify(entry LogzEntry) error {
	if n.HttpMethod == "POST" {
		req, err := http.NewRequest("POST", n.WebhookURL, strings.NewReader(entry.GetMessage()))
		if err != nil {
			return fmt.Errorf("HTTP request creation error: %w", err)
		}
		if n.AuthToken != "" {
			req.Header.Set("Authorization", "Bearer "+n.AuthToken)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := n.WebClient().Do(req)
		if err != nil {
			return fmt.Errorf("HTTP request error: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP request failed: %s", resp.Status)
		}
	} else {
		return fmt.Errorf("unsupported HTTP method: %s", n.HttpMethod)
	}
	return nil
}

// wsNotify sends a WebSocket notification.
func (n *NotifierImpl) wsNotify(entry LogzEntry) error {
	_ = n.AuthToken + "|" + entry.GetMessage()
	// message := n.AuthToken + "|" + entry.GetMessage()
	// Temporarily disabled due to external dependency on zmq4
	// Uncomment and ensure the required libraries are installed if needed in the future
	//if _, err := n.Websocket().Send(message, 0); err != nil {
	//	return fmt.Errorf("WebSocket error: %w", err)
	//}
	return nil
}

// dbusNotify sends a DBus notification.
func (n *NotifierImpl) dbusNotify(entry LogzEntry) error {
	output := n.AuthToken + "|" + entry.GetMessage()
	dbusObj := n.DBusClient().Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if call := dbusObj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0), "", output, []string{}, map[string]dbus.Variant{}, int32(5000)); call.Err != nil {
		return fmt.Errorf("DBus call error: %w", call.Err)
	}
	return nil
}

// Enable activates the notifier.
func (n *NotifierImpl) Enable() { n.EnabledFlag = true }

// Disable deactivates the notifier.
func (n *NotifierImpl) Disable() { n.EnabledFlag = false }

// Enabled checks if the notifier is active.
func (n *NotifierImpl) Enabled() bool { return n.EnabledFlag }

// WebServer returns the HTTP server instance.
func (n *NotifierImpl) WebServer() *http.Server { return n.NotifierManager.WebServer() }

// Temporarily disabled due to external dependency on zmq4
// Uncomment and ensure the required libraries are installed if needed in the future
// Websocket returns the WebSocket instance.
//func (n *NotifierImpl) Websocket() *zmq4.Socket { return n.NotifierManager.Websocket() }

// WebClient returns the HTTP client instance.
func (n *NotifierImpl) WebClient() *http.Client { return n.NotifierManager.WebClient() }

// DBusClient returns the DBus connection instance.
func (n *NotifierImpl) DBusClient() *dbus.Conn { return n.NotifierManager.DBusClient() }

// contains checks if a slice contains a specific value.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// HTTPNotifier is a notifier that sends HTTP notifications.
type HTTPNotifier struct {
	NotifierImpl
}

// NewHTTPNotifier creates a new HTTPNotifier instance.
func NewHTTPNotifier(webhookURL, authToken string) *HTTPNotifier {
	return &HTTPNotifier{
		NotifierImpl: NotifierImpl{
			WebhookURL: webhookURL,
			AuthToken:  authToken,
			HttpMethod: "POST",
		},
	}
}

// Notify sends an HTTP notification.
func (n *HTTPNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}
	req, err := http.NewRequest(n.HttpMethod, n.WebhookURL, strings.NewReader(entry.GetMessage()))
	if err != nil {
		return fmt.Errorf("HTTPNotifier request creation error: %w", err)
	}
	if n.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.WebClient().Do(req)
	if err != nil {
		return fmt.Errorf("HTTPNotifier request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTPNotifier request failed: %s", resp.Status)
	}
	return nil
}

// ZMQNotifier is a notifier that sends WebSocket notifications.
type ZMQNotifier struct {
	NotifierImpl
}

// NewZMQNotifier creates a new ZMQNotifier instance.
func NewZMQNotifier(endpoint string) *ZMQNotifier {
	return &ZMQNotifier{
		NotifierImpl: NotifierImpl{
			WsEndpoint: endpoint,
		},
	}
}

// Temporarily disabled due to external dependency on zmq4
// Uncomment and ensure the required libraries are installed if needed in the future
// Notify sends a WebSocket notification.
func (n *ZMQNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}
	_ = n.AuthToken + "|" + entry.GetMessage()
	//message := n.AuthToken + "|" + entry.GetMessage()
	//if _, err := n.Websocket().Send(message, 0); err != nil {
	//	return fmt.Errorf("ZMQNotifier error: %w", err)
	//}
	return nil
}

// DBusNotifier is a notifier that sends DBus notifications.
type DBusNotifier struct {
	NotifierImpl
}

// NewDBusNotifier creates a new DBusNotifier instance.
func NewDBusNotifier() *DBusNotifier {
	return &DBusNotifier{
		NotifierImpl: NotifierImpl{},
	}
}

// Notify sends a DBus notification.
func (n *DBusNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}
	output := n.AuthToken + "|" + entry.GetMessage()
	dbusObj := n.DBusClient().Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if call := dbusObj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0), "", output, []string{}, map[string]dbus.Variant{}, int32(5000)); call.Err != nil {
		return fmt.Errorf("DBusNotifier error: %w", call.Err)
	}
	return nil
}

func GetLogPath() string {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home, homeErr = os.UserConfigDir()
		if homeErr != nil {
			home, homeErr = os.UserCacheDir()
			if homeErr != nil {
				home = "/tmp"
			}
		}
	}
	configPath := filepath.Join(home, ".kubex", "logz", "VConfig.json")
	if mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755); mkdirErr != nil && !os.IsExist(mkdirErr) {
		return ""
	}
	return configPath
}
