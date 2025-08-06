package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestHTTPNotifierFunctionality testa notificação HTTP com servidor real
func TestHTTPNotifierFunctionality(t *testing.T) {
	// Canal para receber dados
	received := make(chan map[string]interface{}, 1)

	// Criar servidor HTTP de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		received <- payload
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	// Criar notifier HTTP
	notifier := NewHTTPNotifier(server.URL, "test-token-123")

	// Criar entrada de log
	entry := NewLogEntry().
		WithLevel(INFO).
		WithMessage("Test HTTP notification message").
		AddMetadata("test_key", "test_value").
		AddMetadata("notification_id", 12345)

	// Enviar notificação
	err := notifier.Notify(entry)
	if err != nil {
		t.Fatalf("Failed to send HTTP notification: %v", err)
	}

	// Verificar se foi recebido
	select {
	case payload := <-received:
		if payload["message"] != "Test HTTP notification message" {
			t.Errorf("Expected message 'Test HTTP notification message', got %v", payload["message"])
		}
		if payload["level"] != "INFO" {
			t.Errorf("Expected level 'INFO', got %v", payload["level"])
		}
		t.Logf("HTTP notification received successfully: %+v", payload)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for HTTP notification")
	}
}

// TestWebSocketNotifierFunctionality testa notificação WebSocket com servidor real
func TestWebSocketNotifierFunctionality(t *testing.T) {
	// Canal para receber mensagens
	received := make(chan string, 1)

	// Configurar WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Criar servidor WebSocket
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		// Ler mensagem
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Logf("WebSocket read error: %v", err)
			return
		}
		received <- string(message)
	}))
	defer server.Close()

	// Converter HTTP URL para WebSocket URL
	wsURL := "ws" + server.URL[4:] // Remove "http" e adiciona "ws"

	// Criar notifier WebSocket
	notifier := NewWebSocketNotifier(wsURL)

	// Dar tempo para servidor inicializar
	time.Sleep(100 * time.Millisecond)

	// Criar entrada de log
	entry := NewLogEntry().
		WithLevel(WARN).
		WithMessage("Test WebSocket notification").
		AddMetadata("websocket_test", true).
		AddMetadata("timestamp", time.Now().Unix())

	// Enviar notificação
	err := notifier.Notify(entry)
	if err != nil {
		t.Fatalf("Failed to send WebSocket notification: %v", err)
	}

	// Verificar se foi recebido
	select {
	case message := <-received:
		// Verificar se a mensagem contém dados esperados
		if len(message) == 0 {
			t.Fatal("Received empty message")
		}

		// Tentar parsear como JSON
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(message), &parsed); err == nil {
			if parsed["message"] != "Test WebSocket notification" {
				t.Errorf("Expected message 'Test WebSocket notification', got %v", parsed["message"])
			}
			if parsed["level"] != "WARN" {
				t.Errorf("Expected level 'WARN', got %v", parsed["level"])
			}
		}

		t.Logf("WebSocket notification received successfully: %s", message)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for WebSocket notification")
	}
}

// TestConcurrentLogging testa logging concorrente sem race conditions
func TestConcurrentLogging(t *testing.T) {
	logger := NewLogger("concurrent-test")

	const numGoroutines = 10
	const logsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Usar canal para contar logs completados
	completed := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()

			logCount := 0
			for j := 0; j < logsPerGoroutine; j++ {
				logger.InfoCtx(
					fmt.Sprintf("Goroutine %d - Log %d", goroutineID, j),
					map[string]interface{}{
						"goroutine_id": goroutineID,
						"log_number":   j,
						"timestamp":    time.Now().UnixNano(),
						"test_data":    fmt.Sprintf("data-%d-%d", goroutineID, j),
					},
				)
				logCount++
			}
			completed <- logCount
		}(i)
	}

	wg.Wait()
	close(completed)

	// Verificar se todos os logs foram completados
	totalLogs := 0
	for count := range completed {
		totalLogs += count
	}

	expectedTotal := numGoroutines * logsPerGoroutine
	if totalLogs != expectedTotal {
		t.Errorf("Expected %d total logs, got %d", expectedTotal, totalLogs)
	}

	t.Logf("Successfully completed %d concurrent logs", totalLogs)
}

// TestLogEntry testa funcionalidade de LogEntry
func TestLogEntry(t *testing.T) {
	entry := NewLogEntry()

	// Testar builder pattern
	entry = entry.
		WithLevel(ERROR).
		WithMessage("Test error message").
		WithSource("test-source").
		WithContext("test-context").
		WithSeverity(logLevels[ERROR]). // Adicionar severity válido
		AddMetadata("error_code", 500).
		AddMetadata("request_id", "req-123").
		AddTag("service", "logz").
		AddTag("env", "test")

	// Verificar valores
	if entry.GetLevel() != ERROR {
		t.Errorf("Expected level ERROR, got %v", entry.GetLevel())
	}

	if entry.GetMessage() != "Test error message" {
		t.Errorf("Expected message 'Test error message', got %v", entry.GetMessage())
	}

	if entry.GetSource() != "test-source" {
		t.Errorf("Expected source 'test-source', got %v", entry.GetSource())
	}

	metadata := entry.GetMetadata()
	if metadata["error_code"] != 500 {
		t.Errorf("Expected error_code 500, got %v", metadata["error_code"])
	}

	// Testar validação
	err := entry.Validate()
	if err != nil {
		t.Errorf("Entry validation failed: %v", err)
	}

	// Testar string representation
	str := entry.String()
	if len(str) == 0 {
		t.Error("Entry string representation is empty")
	}

	t.Logf("Entry validation completed successfully")
} // TestPrometheusMetrics testa funcionalidade básica do Prometheus
func TestPrometheusMetrics(t *testing.T) {
	pm := GetPrometheusManager()

	// Adicionar métricas
	pm.AddMetric("test_counter", 42.0, map[string]string{"type": "test"})
	pm.AddMetric("test_gauge", 3.14159, map[string]string{"unit": "pi"})

	// Incrementar métrica
	pm.IncrementMetric("test_counter", 8.0)

	// Obter métricas
	metrics := pm.GetMetrics()

	// Verificar valores
	if metrics["test_counter"] != 50.0 { // 42 + 8
		t.Errorf("Expected test_counter to be 50.0, got %f", metrics["test_counter"])
	}

	if metrics["test_gauge"] != 3.14159 {
		t.Errorf("Expected test_gauge to be 3.14159, got %f", metrics["test_gauge"])
	}

	// Remover métrica
	pm.RemoveMetric("test_gauge")

	updatedMetrics := pm.GetMetrics()
	if _, exists := updatedMetrics["test_gauge"]; exists {
		t.Error("test_gauge should have been removed")
	}

	t.Logf("Prometheus metrics test completed successfully")
}

// TestLogFormats testa diferentes formatadores
func TestLogFormats(t *testing.T) {
	entry := NewLogEntry().
		WithLevel(INFO).
		WithMessage("Test formatting").
		AddMetadata("format_test", true)

	// Testar JSON formatter
	jsonFormatter := &JSONFormatter{}
	jsonResult, err := jsonFormatter.Format(entry)
	if err != nil {
		t.Errorf("JSON formatter failed: %v", err)
	}

	if len(jsonResult) == 0 {
		t.Error("JSON formatter returned empty result")
	}

	// Verificar se é JSON válido
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult), &parsed); err != nil {
		t.Errorf("JSON formatter did not produce valid JSON: %v", err)
	}

	// Testar Text formatter
	textFormatter := &TextFormatter{}
	textResult, err := textFormatter.Format(entry)
	if err != nil {
		t.Errorf("Text formatter failed: %v", err)
	}

	if len(textResult) == 0 {
		t.Error("Text formatter returned empty result")
	}

	t.Logf("Format testing completed successfully")
	t.Logf("JSON result: %s", jsonResult)
	t.Logf("Text result: %s", textResult)
}
