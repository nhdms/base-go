package tables

import "github.com/nhdms/base-go/pkg/dbtool"

func GetUserTable() *dbtool.Table {
	return &dbtool.Table{
		Name:      "users",
		AIColumns: []string{"id"},
		ColumnMapper: map[string]string{
			"saleChannel": `"saleChannel"`,
		},
		IgnoreColumns: []string{},
		DefaultAlias:  "u",
	}
}
