package tables

import "github.com/nhdms/base-go/pkg/dbtool"

func GetWebhookEventsTable() *dbtool.Table {
	return &dbtool.Table{
		Name:          "webhook_events",
		AIColumns:     []string{"id"},
		ColumnMapper:  map[string]string{},
		IgnoreColumns: []string{"is_retry"},
		DefaultAlias:  "we",
	}
}
