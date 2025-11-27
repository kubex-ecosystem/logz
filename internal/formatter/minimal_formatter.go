package formatter

import (
	"github.com/kubex-ecosystem/logz/interfaces"
)

// MinimalFormatter: LEVEL message\n
type MinimalFormatter struct{}

func NewMinimalFormatter() *MinimalFormatter { return &MinimalFormatter{} }

func (f *MinimalFormatter) Format(e interfaces.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	// line := fmt.Sprintf("%s %s\n", e.Level, e.Message)
	// return []byte(line), nil
	return nil, nil
}
