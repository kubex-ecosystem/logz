package interfaces

import (
	"encoding/json"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type DynamicFormatter[T kbx.Entry] struct {
	compactThreshold kbx.Level
	prettyThreshold  kbx.Level
	jsonThreshold    kbx.Level

	enrichers []func(*T)
	filters   []func(*T) bool
}

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
	case e.GetLevel() >= f.jsonThreshold:
		return json.Marshal(e)

	case e.GetLevel() >= f.prettyThreshold:
		return f.Format(e)

	case e.GetLevel() >= f.compactThreshold:
		return f.Format(e)

	default:
		return f.Format(e)
	}
}
