package dbtool

type Table struct {
	Name           string
	AIColumns      []string
	ColumnMapper   map[string]string
	IgnoreColumns  []string
	DefaultAlias   string
	NotNullColumns map[string]interface{}
}
