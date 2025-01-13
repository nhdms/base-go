package pglogicalstream

import "fmt"

const (
	TlsNoVerify      = "none"
	TlsRequireVerify = "require"
)

type Config struct {
	DbHost                     string   `mapstructure:"host"`
	DbPassword                 string   `mapstructure:"password"`
	DbUser                     string   `mapstructure:"user"`
	DbPort                     int      `mapstructure:"port"`
	DbName                     string   `mapstructure:"database"`
	DbSchema                   string   `mapstructure:"schema"`
	DbTables                   []string `mapstructure:"tables"`
	ReplicationSlotName        string   `mapstructure:"slot_name"`
	TlsVerify                  string   `mapstructure:"tls_verify"`
	StreamOldData              bool     `mapstructure:"stream_old_data"`
	SeparateChanges            bool     `mapstructure:"separate_changes"`
	SnapshotMemorySafetyFactor float64  `mapstructure:"snapshot_memory_safety_factor"`
	BatchSize                  int      `mapstructure:"batch_size"`
	AutoAck                    bool     `mapstructure:"auto_ack"`
	LastLsn                    string   `mapstructure:"last_lsn"`
	SnapshotOffset             int64    `mapstructure:"snapshot_offset"`
}

func (c *Config) InitDefaultAndValidate() error {
	if c.DbHost == "" {
		return fmt.Errorf("db_host must be provided")
	}
	if c.DbPassword == "" {
		return fmt.Errorf("db_password must be provided")
	}
	if c.DbUser == "" {
		return fmt.Errorf("db_user must be provided")
	}
	if c.DbPort == 0 {
		return fmt.Errorf("db_port must be provided")
	}
	if c.DbName == "" {
		return fmt.Errorf("db_name must be provided")
	}
	if c.DbSchema == "" {
		return fmt.Errorf("db_schema must be provided")
	}
	if c.ReplicationSlotName == "" {
		return fmt.Errorf("replication_slot_name must be provided")
	}

	if c.BatchSize < 1 {
		c.BatchSize = 1000
	}
	if c.SnapshotMemorySafetyFactor == 0 {
		c.SnapshotMemorySafetyFactor = 0.7
	}
	if len(c.TlsVerify) == 0 {
		c.TlsVerify = TlsNoVerify
	}

	return nil
}
