package interfaces

import (
	"encoding/json"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// DynamicFormatter é um formatador genérico que pode ser usado para formatar logs de diferentes níveis.
type DynamicFormatter[T kbx.Entry] struct {
	compactThreshold kbx.Level
	prettyThreshold  kbx.Level
	jsonThreshold    kbx.Level
	hudThreshold     kbx.Level
	pretty           bool

	enrichers []func(*T)
	filters   []func(*T) bool
}

// NewDynamicFormatter creates a new DynamicFormatter.
func NewDynamicFormatter[T kbx.Entry](compactThreshold kbx.Level, prettyThreshold kbx.Level, jsonThreshold kbx.Level, hudThreshold kbx.Level) *DynamicFormatter[T] {
	return &DynamicFormatter[T]{compactThreshold: compactThreshold, prettyThreshold: prettyThreshold, jsonThreshold: jsonThreshold, hudThreshold: hudThreshold}
}

// Format implementa a interface do logz.
func (f *DynamicFormatter[T]) Format(e T) ([]byte, error) {
	// 1. Enrichment
	for _, enrich := range f.enrichers {
		enrich(&e)
	}

	// 2. Filtering
	for _, filter := range f.filters {
		if !filter(&e) {
			return nil, nil // drop
		}
	}

	// 3. Heurística de formato
	switch {
	case kbx.ParseLevel(e.GetLevel()) >= f.jsonThreshold:
		return json.Marshal(e)

	case kbx.ParseLevel(e.GetLevel()) >= f.prettyThreshold:
		return f.Format(e)

	case kbx.ParseLevel(e.GetLevel()) >= f.compactThreshold:
		return f.Format(e)

	case kbx.ParseLevel(e.GetLevel()) == f.hudThreshold:
		return f.Format(e)
	default:
		return f.Format(e)
	}
}
