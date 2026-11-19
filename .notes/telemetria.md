# Meta-prompt: Implementação de Telemetria Zero-Allocation no Logz

## 1. Papel e Contexto

Você atuará como um **Engenheiro Go especialista em Performance e Observabilidade**. Sua tarefa é concluir a integração do Prometheus no módulo `logz`.
Este módulo é o sistema nervoso central do ecossistema Kubex e será consumido pela máquina de estados (FSM atômica) do `ethyr`, que realiza milhares de transições por segundo.

**Restrição Absoluta:** O footprint de memória deve ser mínimo. Alocações dinâmicas no caminho crítico (como o uso de `WithLabelValues` em tempo de execução para a FSM) são **estritamente proibidas**, pois destroem a performance e sobrecarregam o Garbage Collector.

## 2. Tarefa 1: Refatoração do `prometheus_hook.go`

Você deve atualizar o arquivo `prometheus_hook.go` para registrar as métricas do laboratório e, crucialmente, **pré-alocar os ponteiros da FSM** no momento de `init()`.

Implemente exatamente a seguinte estrutura de métricas:

- `logz_entries_total` (CounterVec para logs gerais)
- `gnyx_requests_total` (CounterVec: provider, status)
- `gnyx_request_latency_seconds` (HistogramVec: provider. Buckets: 0.1, 0.5, 1.0, 2.0, 5.0)
- `gnyx_output_tokens_total` (Counter)
- `ethyr_fsm_transitions_total` (CounterVec: state)

**A Regra de Ouro (Zero-Allocation):**
Exporte variáveis globais pre-alocadas para os estados da FSM. O `init()` deve fazer o `WithLabelValues` e salvar o ponteiro na variável exportada.
Exemplo obrigatório:

```go
var (
    FsmStatePreparing prometheus.Counter
    FsmStateExecuting prometheus.Counter
    FsmStateParsing   prometheus.Counter
)

func init() {
    // ... MustRegister de todas as métricas ...
    FsmStatePreparing = ethyrFsmTransitions.WithLabelValues("StateCogPreparing")
    FsmStateExecuting = ethyrFsmTransitions.WithLabelValues("StateCogExecuting")
    FsmStateParsing   = ethyrFsmTransitions.WithLabelValues("StateCogParsing")
}
```

Você NÃO DEVE ALTERAR NADA NO LOGZ ALÉM DO PROMETHEUS E MÉTRICAS sem minha permissão.
