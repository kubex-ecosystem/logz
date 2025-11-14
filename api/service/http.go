// Package service provides HTTP-related services.
package service

import (
	"net/http"

	"github.com/godbus/dbus/v5"
	core "github.com/kubex-ecosystem/logz/internal/loggerz"
)

func Run() error { return core.Run() }

func Start(port string) error                      { return core.Start(port) }
func Stop() error                                  { return core.Stop() }
func Server() *http.Server                         { return core.Server() }
func Client() *http.Client                         { return core.Client() }
func DBus() *dbus.Conn                             { return core.DBus() }
func IsRunning() bool                              { return core.IsRunning() }
func GetServiceInfo() (int, string, string, error) { return core.GetServiceInfo() }
