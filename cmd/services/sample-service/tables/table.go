package tables

import "github.com/nhdms/base-go/pkg/dbtool"

func GetSampleTable() *dbtool.Table {
	return &dbtool.Table{
		Name:          "sample_events",
		AIColumns:     []string{"id"},
		ColumnMapper:  map[string]string{},
		IgnoreColumns: []string{},
		DefaultAlias:  "u",
	}
}
