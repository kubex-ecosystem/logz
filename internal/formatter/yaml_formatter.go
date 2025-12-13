package formatter

import (
	"github.com/kubex-ecosystem/logz/internal/module/kbx"
	"gopkg.in/yaml.v3"
)

type YamlFormatter struct {
	Pretty bool
}

func NewYamlFormatter(pretty bool) Formatter {
	return &YamlFormatter{Pretty: pretty}
}

func (f *YamlFormatter) Name() string {
	return "yaml"
}

type yamlOutput struct {
	Entries []kbx.Entry `yaml:"entries"`
}

func (f *YamlFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if f.Pretty {
		return yaml.Marshal(yamlOutput{Entries: []kbx.Entry{e}})
	}
	return yaml.Marshal(yamlOutput{Entries: []kbx.Entry{e}})
}
