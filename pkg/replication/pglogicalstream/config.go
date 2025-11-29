package pglogicalstream

import (
	"fmt"
	"time"
)

const (
	TlsNoVerify      = "none"
	TlsRequireVerify = "require"

	// Default values for recovery and heartbeat
	DefaultHeartbeatInterval       = 10 * time.Second
	DefaultMaxRetries              = 3
	DefaultRetryDelay              = 5 * time.Second
	DefaultSlotInvalidAutoRecreate = true
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

	// New configuration options for error recovery
	HeartbeatInterval       time.Duration `mapstructure:"heartbeat_interval"`         // Interval for sending standby status updates
	MaxRetries              int           `mapstructure:"max_retries"`                // Max retries for recoverable errors
	RetryDelay              time.Duration `mapstructure:"retry_delay"`                // Delay between retries
	SlotInvalidAutoRecreate bool          `mapstructure:"slot_invalid_auto_recreate"` // Auto-recreate slot on SQLSTATE 55000
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

	// Set defaults for recovery options
	if c.HeartbeatInterval == 0 {
		c.HeartbeatInterval = DefaultHeartbeatInterval
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = DefaultMaxRetries
	}
	if c.RetryDelay == 0 {
		c.RetryDelay = DefaultRetryDelay
	}
	// Note: SlotInvalidAutoRecreate defaults to false (zero value)
	// but we want it to default to true for safety
	// This is handled in NewPgStream

	return nil
}
