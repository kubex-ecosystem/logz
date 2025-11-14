// Package formatters provides functionality to format data as JSON.
package formatters

type JSONFormatter = LogFormatter

func NewJSONFormatter() JSONFormatter {
	return &JSONFormatterImpl{}
}

func NewJSONFormatterImpl() *JSONFormatterImpl {
	return &JSONFormatterImpl{}
}
