package tables

import "gitlab.com/a7923/athena-go/pkg/dbtool"

func GetEventTable() *dbtool.Table {
	return &dbtool.Table{
		Name:      "event_events",
		AIColumns: []string{"id"},
		ColumnMapper: map[string]string{
		},
		IgnoreColumns: []string{},
		DefaultAlias:  "u",
	}
}
