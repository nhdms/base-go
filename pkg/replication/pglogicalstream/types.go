package pglogicalstream

type Wal2JsonChanges struct {
	Lsn     *string          `json:"lsn"`
	Changes []Wal2JsonChange `json:"change"`
}

type OldKeys struct {
	Keynames  []string      `json:"keynames"`
	Keytypes  []string      `json:"keytypes"`
	Keyvalues []interface{} `json:"keyvalues"`
}

type Wal2JsonChange struct {
	Kind         string        `json:"kind"`
	Schema       string        `json:"schema"`
	Table        string        `json:"table"`
	ColumnNames  []string      `json:"columnnames"`
	ColumnTypes  []string      `json:"columntypes"`
	ColumnValues []interface{} `json:"columnvalues"`
	OldData      OldKeys       `json:"oldkeys"`
}
type OnMessage = func(message Wal2JsonChanges)

func (c *Wal2JsonChange) GetValue(column string) (oldVal, newVal interface{}) {
	for i, name := range c.ColumnNames {
		if name == column {
			newVal = c.ColumnValues[i]
			break
		}
	}

	for i, name := range c.OldData.Keynames {
		if name == column {
			oldVal = c.OldData.Keyvalues[i]
			break
		}
	}
	return oldVal, newVal
}
