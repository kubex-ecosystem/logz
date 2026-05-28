package observability

import (
	"github.com/kubex-ecosystem/logz/internal/events"
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Métricas Gerais (Logz)
	logzEntriesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "logz_entries_total", Help: "Total logz entries"},
		[]string{"level", "source", "context"},
	)

	// Métricas do Lab (Genkit & FSM)
	GnyxRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "gnyx_requests_total", Help: "Total requests processed"},
		[]string{"provider", "status"},
	)

	GnyxRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gnyx_request_latency_seconds",
			Help:    "Request latency TTFT",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
		},
		[]string{"provider"},
	)

	GnyxOutputTokens = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "gnyx_output_tokens_total", Help: "Total tokens"},
	)

	// FSM State Vector (A Magia do Zero Allocation)
	ethyrFsmTransitions = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "ethyr_fsm_transitions_total", Help: "FSM Atômica"},
		[]string{"state"},
	)

	// CACHE DE PONTEIROS (Alocado no startup, custo zero no runtime)
	FsmStatePreparing prometheus.Counter
	FsmStateExecuting prometheus.Counter
	FsmStateParsing   prometheus.Counter
)

func init() {
	prometheus.MustRegister(
		logzEntriesCounter, GnyxRequestsTotal, GnyxRequestLatency,
		GnyxOutputTokens, ethyrFsmTransitions,
	)

	// Pré-alocando os ponteiros! O Ethyr vai chamar só isso: `observability.FsmStateExecuting.Inc()`
	FsmStatePreparing = ethyrFsmTransitions.WithLabelValues("StateCogPreparing")
	FsmStateExecuting = ethyrFsmTransitions.WithLabelValues("StateCogExecuting")
	FsmStateParsing = ethyrFsmTransitions.WithLabelValues("StateCogParsing")
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
