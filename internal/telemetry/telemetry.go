// Package telemetry exports telemetry for Logz.
package telemetry

import (
	"context"
	"time"

	"github.com/kubex-ecosystem/logz/interfaces"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
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

func (t *CollectorImpl[O, C]) Run(ctx context.Context) error {
	t.context = ctx

	// TODO: Implement!

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
