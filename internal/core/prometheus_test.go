package core

import (
	"testing"
)

// TestPrometheusBasic testa funcionalidades básicas do Prometheus sem HTTP server
func TestPrometheusBasic(t *testing.T) {
	// Criar instância sem inicializar singleton global (evitar problemas)
	pm := &PrometheusManager{
		enabled:         false,
		Metrics:         make(map[string]Metric),
		exportWhitelist: make(map[string]bool),
	}

	// Testar adição de métricas
	pm.AddMetric("test_counter", 42.0, map[string]string{"type": "test"})
	pm.AddMetric("test_gauge", 3.14159, map[string]string{"unit": "float"})

	// Testar incremento
	pm.IncrementMetric("test_counter", 8.0)

	// Obter métricas
	metrics := pm.GetMetrics()

	// Verificar valores
	if metrics["test_counter"] != 50.0 {
		t.Errorf("Expected test_counter to be 50.0, got %f", metrics["test_counter"])
	}

	if metrics["test_gauge"] != 3.14159 {
		t.Errorf("Expected test_gauge to be 3.14159, got %f", metrics["test_gauge"])
	}

	// Testar remoção
	pm.RemoveMetric("test_gauge")
	updatedMetrics := pm.GetMetrics()

	if _, exists := updatedMetrics["test_gauge"]; exists {
		t.Error("test_gauge should have been removed")
	}

	t.Logf("Basic Prometheus metrics test completed successfully")
}

// TestPrometheusValidation testa validação de nomes de métricas
func TestPrometheusValidation(t *testing.T) {
	pm := &PrometheusManager{
		enabled: false,
		Metrics: make(map[string]Metric),
	}

	// Testar nomes válidos
	validNames := []string{
		"valid_metric",
		"http_requests_total",
		"cpu_usage_percent",
		"_private_metric",
		"metric_with_numbers_123",
	}

	for _, name := range validNames {
		pm.AddMetric(name, 1.0, nil)
		if _, exists := pm.Metrics[name]; !exists {
			t.Errorf("Valid metric name '%s' was rejected", name)
		}
	}

	// Testar nomes inválidos (devem ser rejeitados)
	invalidNames := []string{
		"123invalid",
		"invalid-metric",
		"invalid.metric",
		"invalid metric",
	}

	initialCount := len(pm.Metrics)
	for _, name := range invalidNames {
		pm.AddMetric(name, 1.0, nil)
	}

	// Nenhuma métrica inválida deve ter sido adicionada
	if len(pm.Metrics) != initialCount {
		t.Errorf("Invalid metric names were accepted")
	}

	t.Logf("Prometheus validation test completed successfully")
}

// TestPrometheusWhitelist testa funcionalidade de whitelist
func TestPrometheusWhitelist(t *testing.T) {
	pm := &PrometheusManager{
		enabled:         false,
		Metrics:         make(map[string]Metric),
		exportWhitelist: make(map[string]bool),
	}

	// Adicionar várias métricas
	pm.AddMetric("metric_a", 1.0, nil)
	pm.AddMetric("metric_b", 2.0, nil)
	pm.AddMetric("metric_c", 3.0, nil)

	// Definir whitelist para apenas duas métricas
	pm.SetExportWhitelist([]string{"metric_a", "metric_c"})

	// Obter métricas (deve respeitar whitelist)
	metrics := pm.GetMetrics()

	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics in whitelist, got %d", len(metrics))
	}

	if _, exists := metrics["metric_a"]; !exists {
		t.Error("metric_a should be in whitelist")
	}

	if _, exists := metrics["metric_c"]; !exists {
		t.Error("metric_c should be in whitelist")
	}

	if _, exists := metrics["metric_b"]; exists {
		t.Error("metric_b should not be in whitelist")
	}

	t.Logf("Prometheus whitelist test completed successfully")
}
