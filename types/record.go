package interfaces

import "github.com/kubex-ecosystem/logz/internal/core"

// Record é o contrato mínimo para qualquer coisa logável no core.
// Entry implementa isso, mas nada impede de usar outros tipos no futuro
// (TraceEvent, MetricPoint, etc).
type Record = core.Entry
