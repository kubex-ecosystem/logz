package formatters

// TextFormatter formats log entries in plain text.
type TextFormatter = LogFormatter

func NewTextFormatter() TextFormatter {
	return &TextFormatterImpl{}
}
