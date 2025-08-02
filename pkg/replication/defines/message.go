package replication

type EventType string
type Table struct {
	Name   string
	Schema string
}
type Message struct {
	Position string    `json:"position"`
	Type     EventType `json:"type"`
	Table    Table     `json:"table"`

	Data    map[string]interface{} `json:"data,omitempty"`
	OldData map[string]interface{} `json:"old_data,omitempty"`
}

type Messages []*Message
