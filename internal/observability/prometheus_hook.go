// Package observability exports observability hooks for Logz.
package observability

import (
	"github.com/kubex-ecosystem/logz/internal/events"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/prometheus/client_golang/prometheus"
)

// logzEntriesCounter increments for each logz entry processed.
// It is registered with the default Prometheus registry.
//
// *********************WARNING*********************:
// Do not modify this variable directly.
// Do not use this variable in production code.
// It is initialized in the init() function.
// It is used internally by the Logz package.
var logzEntriesCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "logz_entries_total",
		Help: "Total number of logz entries processed",
	},
	[]string{"level", "source", "context"}, // Dimensions
)

func init() {
	// Registers the metric at system boot
	prometheus.MustRegister(logzEntriesCounter)
}

// NewPrometheusHook creates a new hook that collects Prometheus metrics.
func NewPrometheusHook() events.HookFunc {
	return func(entry kbx.Entry) error {
		lvl := string(entry.GetLevel())
		src := entry.GetCaller()
		ctx := entry.GetContext()

		if src == "" {
			src = "unknown"
		}
		if ctx == "" {
			ctx = "global"
		}

		// Increments atomically in memory (~15 nanoseconds)
		logzEntriesCounter.WithLabelValues(lvl, src, ctx).Inc()

		// Releases for the next hook or for the Formatter
		return nil
	}
}
