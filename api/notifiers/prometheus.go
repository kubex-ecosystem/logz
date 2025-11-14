// Package integrations provides integration with various external services.
package notifiers

import (
	"github.com/kubex-ecosystem/logz/internal/core"
)

type Metric = core.Metric
type PrometheusManager = core.PrometheusManager

// GetPrometheusManager returns the singleton instance of PrometheusManager.
// If it doesn't exist, it initializes a new one.
func GetPrometheusManager() *PrometheusManager {
	return core.GetPrometheusManager()
}
