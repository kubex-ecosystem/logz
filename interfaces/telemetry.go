package interfaces

import (
	"context"
)

// telemetryHook is any type that can observe telemetry data.
type telemetryHook interface {
	Run(ctx context.Context) error
	Close() error
}

// Observer is any type that can observe telemetry data.
type Observer interface {
	telemetryHook
	Listen(event <-chan any) error
	Notify(data any) error
	Error() error
}

// Collector is any type that can collect telemetry data.
type Collector[O Observer] interface {
}

// MetricsCollector is any type that can collect metrics data.
type MetricsCollector = Collector[Observer]

// TracingCollector is any type that can collect trace data.
type TracingCollector = Collector[Observer]

// ProfilerCollector is any type that can collect profiler data.
type ProfilerCollector = Collector[Observer]

// TelemetryManager is any type that can manage telemetry data.
type TelemetryManager interface {
	Run(ctx context.Context) error
	Close() error
	Error() error
}
