package formatter

import (
	"fmt"

	"github.com/kubex-ecosystem/logz/interfaces"
)

// MinimalFormatter: LEVEL message\n

type MinimalFormatter struct {
	pretty bool
}

func NewMinimalFormatter(pretty bool) interfaces.Formatter {
	return &MinimalFormatter{pretty: pretty}
}

func (f *MinimalFormatter) Format(e interfaces.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	line := fmt.Sprintf("%s %s\n", e.GetLevel(), e.String())
	return []byte(line), nil
}
