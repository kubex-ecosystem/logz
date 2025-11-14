// Package notifiers provides the interface for notifier implementations.
package notifiers

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/kubex-ecosystem/logz/internal/config"
	interfaces "github.com/kubex-ecosystem/logz/internal/interfaces"

	apiC "github.com/kubex-ecosystem/logz/api/config"
)

type Config = interfaces.Config
type ConfigImpl = interfaces.ConfigManager

type ConfigManager = interfaces.ConfigManager

// type ConfigManagerImpl = interfaces.ConfigManagerImpl

func NewLogzConfigManager() ConfigManager {
	return interfaces.NewConfigManager()
}

type NotifierManager = interfaces.NotifierManager
type NotifierManagerImpl = interfaces.NotifierManagerImpl

// NewLogzNotifierManager creates a new instance of LogzNotifierManager.
func NewLogzNotifierManager(notifiers map[string]Notifier) NotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]Notifier)
	}
	return NewNotifierManager(notifiers)
}

type Notifier = apiC.Notifier
type NotifierImpl = apiC.NotifierImpl

// NewNotifier creates a new Notifier service instance.
func NewNotifier(
	manager apiC.NotifierManager,
	enabled bool,
	webhookURL string,
	HTTPMethod string,
	authToken string,
	logLevel string,
	wsEndpoint string,
	whitelist []string,
) Notifier {
	return apiC.NewNotifier(manager, enabled, webhookURL, HTTPMethod, authToken, logLevel, wsEndpoint, whitelist)
}

type HTTPNotifier = apiC.HTTPNotifier

func NewHTTPNotifier(webhookURL string, authToken string) HTTPNotifier {
	return *apiC.NewHTTPNotifier(webhookURL, authToken)
}

type WebSocketNotifier = apiC.WebSocketNotifier

func NewWebSocketNotifierWithConfig(config *NotifierWebSocketConfig) WebSocketNotifier {
	return *apiC.NewWebSocketNotifierWithConfig(config)
}

func NewWebSocketNotifier(endpoint string, authToken string) WebSocketNotifier {
	return *apiC.NewWebSocketNotifier(endpoint, authToken)
}

type NotifierWebSocketConfig = apiC.NotifierWebSocketConfig

func NewNotifierWebSocketConfig(
	TLSClientConfig *tls.Config,
	HandshakeTimeout time.Duration,
	Jar http.CookieJar,
	WriteBufferPool websocket.BufferPool,
	ReadBufferSize, WriteBufferSize int,
	Subprotocols []string,
	EnableCompression bool,
	NetDial func(network, addr string) (net.Conn, error),
	NetDialContext func(ctx context.Context, network, addr string) (net.Conn, error),
	NetDialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error),
	Proxy func(*http.Request) (*url.URL, error),
) *NotifierWebSocketConfig {
	return &apiC.NotifierWebSocketConfig{
		TLSClientConfig:   TLSClientConfig,
		HandshakeTimeout:  HandshakeTimeout,
		Jar:               Jar,
		WriteBufferPool:   WriteBufferPool,
		ReadBufferSize:    ReadBufferSize,
		WriteBufferSize:   WriteBufferSize,
		Subprotocols:      Subprotocols,
		EnableCompression: EnableCompression,
		NetDial:           NetDial,
		NetDialContext:    NetDialContext,
		NetDialTLSContext: NetDialTLSContext,
		Proxy:             Proxy,
	}
}

type DBusNotifier = apiC.DBusNotifier

func NewDBusNotifier() DBusNotifier {
	return *apiC.NewDBusNotifier()
}
