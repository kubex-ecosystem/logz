package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gorilla/websocket"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// NotifierImpl is the implementation of the Notifier interface.
type NotifierImpl struct {
	NotifierManager NotifierManager // Manager for notifier instances.
	EnabledFlag     bool            // Flag indicating if the notifier is enabled.
	WebhookURL      string          // URL for webhook notifications.
	HTTPMethod      string          // HTTP method for webhook notifications.
	AuthToken       string          // Authentication token for notifications.
	LogLevel        string          // Log VLevel for filtering notifications.
	WsEndpoint      string          // WebSocket endpoint for notifications.
	Whitelist       []string        // Whitelist of sources for notifications.
}

// NewNotifier creates a new NotifierImpl instance.
func NewNotifier(manager NotifierManager, enabled bool, webhookURL, HTTPMethod, authToken, logLevel, wsEndpoint string, whitelist []string) Notifier {
	if whitelist == nil {
		whitelist = []string{}
	}
	return &NotifierImpl{
		NotifierManager: manager,
		EnabledFlag:     enabled,
		WebhookURL:      webhookURL,
		HTTPMethod:      HTTPMethod,
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
	if n.HTTPMethod == "POST" {
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
		return fmt.Errorf("unsupported HTTP method: %s", n.HTTPMethod)
	}
	return nil
}

// wsNotify sends a WebSocket notification.
func (n *NotifierImpl) wsNotify(entry LogzEntry) error {
	_ = n.AuthToken + "|" + entry.GetMessage()
	ws, _, err := websocket.DefaultDialer.Dial(n.WsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("WebSocket connection error: %w", err)
	}
	defer ws.Close()
	if err := ws.WriteJSON(entry); err != nil {
		return fmt.Errorf("WebSocket write error: %w", err)
	}

	// Optionally, you can read a message from the WebSocket server to confirm the connection.
	// msg, _, err := ws.ReadMessage()
	if _, _, err := ws.ReadMessage(); err != nil {
		return fmt.Errorf("WebSocket read error: %w", err)
	}

	// Optionally, you can handle the response from the WebSocket server here.
	// For now, we just read a message to confirm the connection is working.
	if _, _, err := ws.ReadMessage(); err != nil {
		return fmt.Errorf("WebSocket read error: %w", err)
	}

	// Print the response message for debugging purposes.
	// You can uncomment or remove this line in production code.
	// This is just an example; you might want to handle the response differently.
	// fmt.Printf("WebSocket response: %s\n", msg)
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

// Websocket returns the Gorilla WebSocket connection instance.
func (n *NotifierImpl) Websocket() *websocket.Conn { return n.NotifierManager.Websocket() }

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
	notifier := &HTTPNotifier{
		NotifierImpl: NotifierImpl{
			EnabledFlag: true, // Habilitado por padrão
			WebhookURL:  webhookURL,
			AuthToken:   authToken,
			HTTPMethod:  "POST",
		},
	}
	// Inicializar NotifierManager padrão se necessário
	if notifier.NotifierManager == nil {
		notifier.NotifierManager = kbx.NewNotifierManagep(nil)
	}
	return notifier
}

// Notify sends an HTTP notification.
func (n *HTTPNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}

	// Serializar a entrada como JSON
	entryData := map[string]interface{}{
		"timestamp": entry.GetTimestamp(),
		"level":     entry.GetLevel(),
		"message":   entry.GetMessage(),
		"metadata":  entry.GetMetadata(),
		"source":    entry.GetSource(),
	}

	jsonData, err := json.Marshal(entryData)
	if err != nil {
		return fmt.Errorf("HTTPNotifier JSON marshal error: %w", err)
	}

	req, err := http.NewRequest(n.HTTPMethod, n.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTPNotifier request creation error: %w", err)
	}

	if n.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.AuthToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTPNotifier request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTPNotifier request failed with status: %s", resp.Status)
	}
	return nil
}

// WebSocketNotifier is a notifier that sends WebSocket notifications.
type WebSocketNotifier struct {
	NotifierImpl
}

// NewWebSocketNotifierWithConfig creates a new WebSocketNotifier instance.
func NewWebSocketNotifierWithConfig(config *notifiers.NotifierWebSocketConfig) *WebSocketNotifier {
	notifier := &WebSocketNotifier{
		NotifierImpl: NotifierImpl{
			EnabledFlag: true, // Habilitado por padrão
			WsEndpoint:  config.Endpoint,
		},
	}
	// Inicializar NotifierManager padrão se necessário
	if notifier.NotifierManager == nil {
		notifier.NotifierManager = notifiers.NewNotifierManager(nil)
	}
	return notifier
}

// NewWebSocketNotifier creates a new WebSocketNotifier instance.
func NewWebSocketNotifier(endpoint string) *WebSocketNotifier {
	notifier := &WebSocketNotifier{
		NotifierImpl: NotifierImpl{
			EnabledFlag: true, // Habilitado por padrão
			WsEndpoint:  endpoint,
		},
	}
	// Inicializar NotifierManager padrão se necessário
	if notifier.NotifierManager == nil {
		notifier.NotifierManager = notifiers.NewNotifierManager(nil)
	}
	return notifier
}

// Notify sends a WebSocket notification.
func (n *WebSocketNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}
	if n.WsEndpoint == "" {
		return fmt.Errorf("WebSocket endpoint not configured")
	}

	// Conectar ao WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(n.WsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("WebSocketNotifier dial error: %w", err)
	}
	defer ws.Close()

	// Serializar a entrada como JSON
	entryData := map[string]interface{}{
		"timestamp": entry.GetTimestamp(),
		"level":     entry.GetLevel(),
		"message":   entry.GetMessage(),
		"metadata":  entry.GetMetadata(),
		"source":    entry.GetSource(),
	}

	jsonData, err := json.Marshal(entryData)
	if err != nil {
		return fmt.Errorf("WebSocketNotifier JSON marshal error: %w", err)
	}

	// Enviar mensagem
	err = ws.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return fmt.Errorf("WebSocketNotifier write error: %w", err)
	}

	return nil
}

// DBusNotifier is a notifier that sends DBus notifications.
type DBusNotifier struct {
	NotifierImpl
}

// NewDBusNotifier creates a new DBusNotifier instance.
func NewDBusNotifier() *DBusNotifier {
	notifier := &DBusNotifier{
		NotifierImpl: NotifierImpl{
			EnabledFlag: true, // Habilitado por padrão
		},
	}
	// Inicializar NotifierManager padrão se necessário
	if notifier.NotifierManager == nil {
		notifier.NotifierManager = notifiers.NewNotifierManager(nil)
	}
	return notifier
}

// Notify sends a DBus notification.
func (n *DBusNotifier) Notify(entry LogzEntry) error {
	if !n.EnabledFlag {
		return nil
	}

	// Tentar conectar ao DBus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("DBusNotifier connection error: %w", err)
	}
	defer conn.Close()

	// Preparar dados da notificação
	appName := "logz"
	iconName := ""
	summary := fmt.Sprintf("Log %s", entry.GetLevel())
	body := entry.GetMessage()
	actions := []string{}
	hints := map[string]dbus.Variant{}
	timeout := int32(5000)

	// Enviar notificação
	dbusObj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := dbusObj.Call("org.freedesktop.Notifications.Notify", 0,
		appName, uint32(0), iconName, summary, body, actions, hints, timeout)

	if call.Err != nil {
		return fmt.Errorf("DBusNotifier call error: %w", call.Err)
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

// ensureConfigExists checks if the configuration file exists, and creates it with default values if it does not.
func ensureConfigExists(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := map[string]any{
			"VlPort":            kbx.DefaultServerPort,
			"VlBindAddress":     kbx.DefaultServerHost,
			"VlAddress":         fmt.Sprintf("%s:%s", kbx.DefaultMode, kbx.DefaultServerPort),
			"VlPidFile":         "logz_srv.pid",
			"VlReadTimeout":     15 * time.Second,
			"VlWriteTimeout":    15 * time.Second,
			"VlIdleTimeout":     60 * time.Second,
			"VlOutput":          kbx.DefaultOutput,
			"VlNotifierManager": notifiers.NewNotifierManager(nil),
			"VlMode":            kbx.DefaultMode,
		}
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		if writeErr := os.WriteFile(configPath, data, 0644); writeErr != nil {
			return fmt.Errorf("failed to create default VConfig: %w", writeErr)
		}
	}
	return nil
}

func getConfigType(configPath string) string {
	configType := filepath.Ext(configPath)
	switch configType {
	case ".yaml":
		return "yaml"
	case ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	default:
		return "json"
	}

}

// getOrDefault returns the value if it is not empty, otherwise returns the default value.
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
