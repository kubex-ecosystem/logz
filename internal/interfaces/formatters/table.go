package formatters

// TableFormatter formats data as a table.
type TableFormatter = LogFormatter

func NewTableFormatter() TableFormatter {
	return &TableFormatterImpl{}
}

func NewTableFormatterImpl() *TableFormatterImpl {
	return &TableFormatterImpl{}
}
