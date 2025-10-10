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
	il "github.com/kubex-ecosystem/logz/internal/core"
)

type Config = il.Config
type ConfigImpl = il.ConfigImpl

type ConfigManager = il.ConfigManager
type ConfigManagerImpl = il.ConfigManagerImpl

func NewLogzConfigManager() ConfigManager {
	return *il.NewConfigManager()
}

type NotifierManager = il.NotifierManager
type NotifierManagerImpl = il.NotifierManagerImpl

// NewLogzNotifierManager creates a new instance of LogzNotifierManager.
func NewLogzNotifierManager(notifiers map[string]Notifier) NotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]Notifier)
	}
	return il.NewNotifierManager(notifiers)
}

type Notifier = il.Notifier
type NotifierImpl = il.NotifierImpl

// NewNotifier creates a new Notifier service instance.
func NewNotifier(
	manager il.NotifierManager,
	enabled bool,
	webhookURL string,
	HTTPMethod string,
	authToken string,
	logLevel string,
	wsEndpoint string,
	whitelist []string,
) Notifier {
	return il.NewNotifier(manager, enabled, webhookURL, HTTPMethod, authToken, logLevel, wsEndpoint, whitelist)
}

type HTTPNotifier = il.HTTPNotifier

func NewHTTPNotifier(webhookURL string, authToken string) HTTPNotifier {
	return *il.NewHTTPNotifier(webhookURL, authToken)
}

type WebSocketNotifier = il.WebSocketNotifier

func NewWebSocketNotifierWithConfig(config *NotifierWebSocketConfig) WebSocketNotifier {
	return *il.NewWebSocketNotifierWithConfig(config)
}

func NewWebSocketNotifier(endpoint string, authToken string) WebSocketNotifier {
	return *il.NewWebSocketNotifier(endpoint)
}

type NotifierWebSocketConfig = il.NotifierWebSocketConfig

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
	return &il.NotifierWebSocketConfig{
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

type DBusNotifier = il.DBusNotifier

func NewDBusNotifier() DBusNotifier {
	return *il.NewDBusNotifier()
}
