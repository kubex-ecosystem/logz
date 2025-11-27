package formatter

import (
	"fmt"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

// MinimalFormatter: LEVEL message\n

type MinimalFormatter struct {
	pretty bool
}

func NewMinimalFormatter(pretty bool) Formatter {
	return &MinimalFormatter{pretty: pretty}
}

func (f *MinimalFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	line := fmt.Sprintf("%s %s\n", e.GetLevel(), e.String())
	return []byte(line), nil
}
