// Package telemetry exports telemetry for Logz.
package telemetry

import (
	"context"
	"net/http"
	"time"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Observer is an alias for interfaces.Observer
type Observer = interfaces.Observer

// CollectorImpl is the implementation of TelemetryCollector.
type CollectorImpl[O Observer, C interfaces.Collector[O]] struct {
	context context.Context

	*kbx.TelemetryConfig
	*kbx.LogzObservabilityOptions

	metricsCache map[kbx.Level]int

	collector O

	error error
}

// NewCollector creates a new TelemetryCollector with an observer
func NewCollector[O Observer, C interfaces.Collector[O]](observer O) interfaces.Collector[O] {
	return &CollectorImpl[O, C]{
		collector: observer,
		TelemetryConfig: &kbx.TelemetryConfig{
			Timeout:  10 * time.Second,
			Interval: 1 * time.Minute,
			Enabled:  true,
		},
	}
}

func (t *CollectorImpl[O, C]) Initialize(opts *kbx.LogzObservabilityOptions) error {
	t.TelemetryConfig = opts.TelemetryConfig
	t.LogzObservabilityOptions = opts
	return nil
}

func (t *CollectorImpl[O, C]) Observer() O {
	return t.collector
}

func (t *CollectorImpl[O, C]) Close() error {
	return t.collector.Close()
}

func (t *CollectorImpl[O, C]) Error() error {
	return t.error
}

// Handler devolve a interface HTTP para o roteador do GNyx plugar na rota /metrics
func (t *CollectorImpl[O, C]) Handler() http.Handler {
	return promhttp.Handler()
}

func (t *CollectorImpl[O, C]) Run(ctx context.Context) error {
	t.context = ctx

	// O Run() do logz pode ser usado apenas para
	// iniciar rotinas de limpeza do metricsCache interno se você usar, ou dar flush.
	// Por enquanto, ele apenas sinaliza que a telemetria está pronta.

	return nil
}

// EmptyCollector is the Observer interface constraints ensure
type EmptyCollector[O Observer, C interfaces.Collector[O]] struct {
	*CollectorImpl[O, C]
}

// NewEmptyCollector creates a new empty TelemetryCollector
// This is only to keep backward compatibility and for internal use only. do not use it in your code.
func NewEmptyCollector() interfaces.Collector[Observer] {
	return &EmptyCollector[Observer, interfaces.Collector[Observer]]{
		CollectorImpl: &CollectorImpl[Observer, interfaces.Collector[Observer]]{
			collector: nil,
		},
	}
}
