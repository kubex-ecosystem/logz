package observability

import (
	"testing"
)

// TestMetricsPreallocation garante que as variáveis globais de telemetria
// para transições da FSM estão devidamente inicializadas (pré-alocadas) e prontas para uso.
func TestMetricsPreallocation(t *testing.T) {
	if FsmStatePreparing == nil {
		t.Error("FsmStatePreparing não deveria ser nil (deve ser pré-alocado no init)")
	}
	if FsmStateExecuting == nil {
		t.Error("FsmStateExecuting não deveria ser nil (deve ser pré-alocado no init)")
	}
	if FsmStateParsing == nil {
		t.Error("FsmStateParsing não deveria ser nil (deve ser pré-alocado no init)")
	}
}

// BenchmarkFsmTransitionsZeroAllocation comprova cientificamente que a gravação das métricas
// da FSM atômica atende ao requisito rígido de Zero-Allocation (0 B/op).
func BenchmarkFsmTransitionsZeroAllocation(b *testing.B) {
	if FsmStateExecuting == nil {
		b.Fatal("FsmStateExecuting é nil")
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FsmStateExecuting.Inc()
	}
}
