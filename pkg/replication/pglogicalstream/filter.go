package pglogicalstream

type ChangeFilter struct {
	tablesWhiteList map[string]bool
	schemaWhiteList string
}

type Filtered func(change Wal2JsonChanges)

func NewChangeFilter(tableSchemas []string, schema string) ChangeFilter {
	tablesMap := map[string]bool{}
	for _, table := range tableSchemas {
		tablesMap[table] = true
	}

	return ChangeFilter{
		tablesWhiteList: tablesMap,
		schemaWhiteList: schema,
	}
}

func (c ChangeFilter) FilterChange(lsn string, changes Wal2JsonChanges, OnFiltered Filtered) {
	if len(changes.Changes) == 0 {
		return
	}

	for _, ch := range changes.Changes {
		var filteredChanges = Wal2JsonChanges{
			Lsn:     &lsn,
			Changes: []Wal2JsonChange{},
		}
		if ch.Schema != c.schemaWhiteList {
			continue
		}

		var (
			tableExist bool
		)

		if _, tableExist = c.tablesWhiteList[ch.Table]; !tableExist {
			continue
		}

		if ch.Kind == "delete" {
			ch.ColumnValues = make([]interface{}, len(ch.OldData.Keyvalues))
			for i, changedValue := range ch.OldData.Keyvalues {
				if len(ch.ColumnValues) == 0 {
					break
				}
				ch.ColumnValues[i] = changedValue
			}
		}

		filteredChanges.Changes = append(filteredChanges.Changes, Wal2JsonChange{
			Kind:         ch.Kind,
			Schema:       ch.Schema,
			Table:        ch.Table,
			ColumnNames:  ch.ColumnNames,
			ColumnTypes:  ch.ColumnTypes,
			ColumnValues: ch.ColumnValues,
			OldData:      ch.OldData,
		})

		OnFiltered(filteredChanges)
	}
}
