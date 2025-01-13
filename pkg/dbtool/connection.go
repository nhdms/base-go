package dbtool

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/spf13/viper"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// DBType represents supported database types
type DBType string

const (
	DBTypeMySQL      DBType = "mysql"
	DBTypePostgreSQL DBType = "postgres"
)

// Config holds database configuration
type Config struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DatabaseName string        `mapstructure:"database"`
	MaxOpenConns int           `mapstructure:"max_open_conns"` // maximum number of open connections to the database
	MaxIdleConns int           `mapstructure:"max_idle_conns"` // maximum number of connections in the idle connection pool
	MaxLifetime  time.Duration `mapstructure:"max_life_time"`  // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`  // maximum amount of time a connection may be idle
}

// ConnectionManager todo implement mysql
// ConnectionManager manages database connection
type ConnectionManager struct {
	db       *sqlx.DB
	dbType   DBType
	mu       sync.RWMutex
	config   *Config
	isActive bool
	isDebug  bool
}

func (cm *ConnectionManager) SetIsDebug(isDebug bool) {
	cm.isDebug = isDebug
}

// NewConnectionManager creates a new instance of ConnectionManager for a specific database type
func NewConnectionManager(dbType DBType, config *Config) (*ConnectionManager, error) {
	if config == nil {
		config = &Config{}
		sub := viper.Sub(string(dbType))
		if sub == nil {
			return nil, fmt.Errorf("config not found for %s", dbType)
		}

		err := sub.Unmarshal(config)
		if err != nil {
			return nil, err
		}
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	cm := &ConnectionManager{
		dbType:  dbType,
		config:  config,
		isDebug: utils.IsTestMode(),
	}

	if err := cm.connect(); err != nil {
		return nil, err
	}
	logger.DefaultLogger.Infof("Connected to %s database: %s:%d/%s", dbType, config.Host, config.Port, config.DatabaseName)
	return cm, nil
}

// validateConfig checks if the configuration is valid
func validateConfig(config *Config) error {
	if config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if config.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}
	if config.User == "" {
		return fmt.Errorf("user cannot be empty")
	}
	if config.DatabaseName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Set default values if not provided
	if config.MaxOpenConns <= 0 {
		config.MaxOpenConns = 25
	}
	if config.MaxIdleConns <= 0 {
		config.MaxIdleConns = 5
	}
	if config.MaxLifetime <= 0 {
		config.MaxLifetime = time.Hour
	}
	if config.MaxIdleTime <= 0 {
		config.MaxIdleTime = 30 * time.Minute
	}

	return nil
}

// connect establishes a database connection
func (cm *ConnectionManager) connect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.isActive {
		return fmt.Errorf("connection already established")
	}

	dsn, err := cm.buildDSN()
	if err != nil {
		return err
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.Connect(string(cm.dbType), dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cm.config.MaxOpenConns)
	db.SetMaxIdleConns(cm.config.MaxIdleConns)
	db.SetConnMaxLifetime(cm.config.MaxLifetime)
	db.SetConnMaxIdleTime(cm.config.MaxIdleTime)
	db.Mapper = reflectx.NewMapperFunc("json", func(str string) string {
		return str
	})

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close() // Clean up if ping fails
		return fmt.Errorf("failed to ping database: %w", err)
	}

	cm.db = db
	cm.isActive = true
	return nil
}

// buildDSN constructs the data source name based on database type
func (cm *ConnectionManager) buildDSN() (string, error) {
	switch cm.dbType {
	case DBTypeMySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=30s&writeTimeout=30s&readTimeout=30s",
			cm.config.User,
			cm.config.Password,
			cm.config.Host,
			cm.config.Port,
			cm.config.DatabaseName,
		), nil
	case DBTypePostgreSQL:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
			cm.config.Host,
			cm.config.Port,
			cm.config.User,
			cm.config.Password,
			cm.config.DatabaseName,
		), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", cm.dbType)
	}
}

// GetConnection returns the database connection
func (cm *ConnectionManager) GetConnection() *sqlx.DB {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.db
}

// Close closes the database connection
func (cm *ConnectionManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.isActive {
		return nil
	}

	if err := cm.db.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	cm.isActive = false
	return nil
}

// Reconnect attempts to reconnect to the database
func (cm *ConnectionManager) Reconnect() error {
	if err := cm.Close(); err != nil {
		return fmt.Errorf("failed to close existing connection: %w", err)
	}
	return cm.connect()
}

// Status returns the current status of the connection
type Status struct {
	IsActive     bool
	OpenConns    int
	IdleConns    int
	InUseConns   int
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
	DatabaseType DBType
	DatabaseName string
}

// GetStatus returns the current status of the connection manager
func (cm *ConnectionManager) GetStatus() Status {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := cm.db.Stats()
	return Status{
		IsActive:     cm.isActive,
		OpenConns:    stats.OpenConnections,
		IdleConns:    stats.Idle,
		InUseConns:   stats.InUse,
		MaxOpenConns: cm.config.MaxOpenConns,
		MaxIdleConns: cm.config.MaxIdleConns,
		MaxLifetime:  cm.config.MaxLifetime,
		MaxIdleTime:  cm.config.MaxIdleTime,
		DatabaseType: cm.dbType,
		DatabaseName: cm.config.DatabaseName,
	}
}

// Health checks if the database connection is healthy
func (cm *ConnectionManager) Health(ctx context.Context) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.isActive {
		return fmt.Errorf("connection is not active")
	}

	return cm.db.PingContext(ctx)
}
