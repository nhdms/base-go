package pglogicalstream

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
)

var pluginArguments = []string{"\"pretty-print\" 'true'"}

// Error types for replication slot issues
var (
	ErrSlotInvalid       = errors.New("replication slot is invalid or WAL has been removed")
	ErrSlotNotFound      = errors.New("replication slot not found")
	ErrConnectionLost    = errors.New("connection to PostgreSQL lost")
	ErrMaxRetriesReached = errors.New("maximum retries reached")
)

// SlotStatus represents the health status of a replication slot
type SlotStatus struct {
	Exists            bool
	Active            bool
	RestartLSN        string
	ConfirmedFlushLSN string
	WALStatus         string // "available", "extended", "unreserved", "lost"
}

type Stream struct {
	pgConn *pgconn.PgConn
	// extra copy of db config is required to establish a new db connection
	// which is required to take snapshot data
	dbConfig     pgconn.Config
	streamCtx    context.Context
	streamCancel context.CancelFunc

	standbyCtxCancel           context.CancelFunc
	clientXLogPos              pglogrepl.LSN
	standbyMessageTimeout      time.Duration
	nextStandbyMessageDeadline time.Time
	messages                   chan Wal2JsonChanges
	snapshotMessages           chan Wal2JsonChanges
	snapshotName               string
	changeFilter               ChangeFilter
	lsnrestart                 pglogrepl.LSN
	slotName                   string
	schema                     string
	tableNames                 []string
	separateChanges            bool
	snapshotBatchSize          int
	snapshotMemorySafetyFactor float64
	m                          sync.Mutex
	stopped                    bool
	autoAck                    bool
	conf                       *Config

	// Error recovery fields
	retryCount    int
	lastHeartbeat time.Time
	connMu        sync.RWMutex // Protects pgConn during reconnection
}

func NewPgStream(config *Config) (*Stream, error) {
	var (
		cfg *pgconn.Config
		err error
	)

	sslVerifyFull := ""
	if config.TlsVerify == TlsRequireVerify {
		sslVerifyFull = "&sslmode=verify-full"
	}

	err = config.InitDefaultAndValidate()
	if err != nil {
		return nil, err
	}

	if cfg, err = pgconn.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?replication=database%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
		sslVerifyFull,
	)); err != nil {
		return nil, err
	}

	if config.TlsVerify == TlsRequireVerify {
		cfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         config.DbHost,
		}
	} else {
		cfg.TLSConfig = nil
	}

	dbConn, err := pgconn.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	var tableNames []string
	for _, table := range config.DbTables {
		tableNames = append(tableNames, table)
	}

	stream := &Stream{
		pgConn:                     dbConn,
		dbConfig:                   *cfg,
		messages:                   make(chan Wal2JsonChanges),
		snapshotMessages:           make(chan Wal2JsonChanges, 100),
		slotName:                   config.ReplicationSlotName,
		schema:                     config.DbSchema,
		snapshotMemorySafetyFactor: config.SnapshotMemorySafetyFactor,
		separateChanges:            config.SeparateChanges,
		snapshotBatchSize:          config.BatchSize,
		tableNames:                 tableNames,
		changeFilter:               NewChangeFilter(tableNames, config.DbSchema),
		m:                          sync.Mutex{},
		stopped:                    false,
		autoAck:                    config.AutoAck,
		conf:                       config,
	}

	result := stream.pgConn.Exec(context.Background(), fmt.Sprintf("DROP PUBLICATION IF EXISTS pglog_stream_%s;", config.ReplicationSlotName))
	_, err = result.ReadAll()
	if err != nil {
		logger.DefaultLogger.Errorf("drop publication if exists error %s", err.Error())
		return nil, err
	}

	for i, table := range tableNames {
		tableNames[i] = fmt.Sprintf("%s.%s", config.DbSchema, table)
	}

	tablesSchemaFilter := fmt.Sprintf("FOR TABLE %s", strings.Join(tableNames, ","))
	logger.DefaultLogger.Infof("Create publication for table schemas with query %s", fmt.Sprintf("CREATE PUBLICATION pglog_stream_%s %s;", config.ReplicationSlotName, tablesSchemaFilter))
	result = stream.pgConn.Exec(context.Background(), fmt.Sprintf("CREATE PUBLICATION pglog_stream_%s %s;", config.ReplicationSlotName, tablesSchemaFilter))
	_, err = result.ReadAll()
	if err != nil {
		logger.DefaultLogger.Fatalf("create publication error %s", err.Error())
	}
	logger.DefaultLogger.Info("Created Postgresql publication", "publication_name", config.ReplicationSlotName)

	sysident, err := pglogrepl.IdentifySystem(context.Background(), stream.pgConn)
	if err != nil {
		logger.DefaultLogger.Fatalf("Failed to identify the system %s", err.Error())
	}

	logger.DefaultLogger.Info("System identification result", "SystemID:", sysident.SystemID, "Timeline:", sysident.Timeline, "XLogPos:", sysident.XLogPos, "DBName:", sysident.DBName)

	var freshlyCreatedSlot = false
	var confirmedLSNFromDB string
	// check is replication slot exist to get last restart SLN
	connExecResult := stream.pgConn.Exec(context.TODO(), fmt.Sprintf("SELECT confirmed_flush_lsn FROM pg_replication_slots WHERE slot_name = '%s'", config.ReplicationSlotName))
	if slotCheckResults, err := connExecResult.ReadAll(); err != nil {
		logger.DefaultLogger.Fatal(err)
	} else {
		if len(slotCheckResults) == 0 || len(slotCheckResults[0].Rows) == 0 {
			// here we create a new replication slot because there is no slot found
			var createSlotResult CreateReplicationSlotResult
			createSlotResult, err = CreateReplicationSlot(context.Background(), stream.pgConn, stream.slotName, "wal2json",
				CreateReplicationSlotOptions{Temporary: false,
					SnapshotAction: "export",
				})
			if err != nil {
				logger.DefaultLogger.Fatalf("Failed to create replication slot for the database: %s", err.Error())
			}
			stream.snapshotName = createSlotResult.SnapshotName
			freshlyCreatedSlot = true
		} else {
			slotCheckRow := slotCheckResults[0].Rows[0]
			confirmedLSNFromDB = string(slotCheckRow[0])
			logger.DefaultLogger.Infow("Replication slot restart LSN extracted from DB", "LSN", confirmedLSNFromDB)
		}
	}

	var lsnrestart pglogrepl.LSN
	if freshlyCreatedSlot {
		lsnrestart = sysident.XLogPos
	} else {
		lsnrestart, _ = pglogrepl.ParseLSN(confirmedLSNFromDB)
	}

	stream.lsnrestart = lsnrestart

	if freshlyCreatedSlot {
		stream.clientXLogPos = sysident.XLogPos
	} else {
		stream.clientXLogPos = lsnrestart
	}

	logger.DefaultLogger.Infof("starting from position %v %v", stream.lsnrestart.String(), stream.clientXLogPos.String())

	stream.standbyMessageTimeout = time.Second * 10
	stream.nextStandbyMessageDeadline = time.Now().Add(stream.standbyMessageTimeout)
	stream.streamCtx, stream.streamCancel = context.WithCancel(context.Background())

	if config.StreamOldData {
		go stream.processSnapshot()
		return stream, nil
	}

	if err := stream.startLr(); err != nil {
		return nil, fmt.Errorf("failed to start logical replication: %w", err)
	}
	go stream.streamMessagesAsync()

	return stream, err
}

// checkSlotStatus queries PostgreSQL to get the current status of the replication slot
func (s *Stream) checkSlotStatus() (*SlotStatus, error) {
	query := fmt.Sprintf(`
		SELECT
			active,
			restart_lsn,
			confirmed_flush_lsn,
			COALESCE(wal_status, 'unknown') as wal_status
		FROM pg_replication_slots
		WHERE slot_name = '%s'
	`, s.slotName)

	result := s.pgConn.Exec(context.Background(), query)
	results, err := result.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to check slot status: %w", err)
	}

	if len(results) == 0 || len(results[0].Rows) == 0 {
		return &SlotStatus{Exists: false}, nil
	}

	row := results[0].Rows[0]
	status := &SlotStatus{
		Exists:            true,
		Active:            string(row[0]) == "t",
		RestartLSN:        string(row[1]),
		ConfirmedFlushLSN: string(row[2]),
		WALStatus:         string(row[3]),
	}

	return status, nil
}

// isSlotInvalidError checks if the error is a SQLSTATE 55000 (invalid replication slot)
func isSlotInvalidError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// SQLSTATE 55000: object_not_in_prerequisite_state
	return strings.Contains(errStr, "55000") ||
		strings.Contains(errStr, "can no longer get changes from replication slot") ||
		strings.Contains(errStr, "replication slot") && strings.Contains(errStr, "does not exist")
}

// isConnectionError checks if the error indicates a lost connection
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Connection-related error codes
		return pgErr.Code == "57P01" || // admin_shutdown
			pgErr.Code == "57P02" || // crash_shutdown
			pgErr.Code == "57P03" || // cannot_connect_now
			pgErr.Code == "08000" || // connection_exception
			pgErr.Code == "08003" || // connection_does_not_exist
			pgErr.Code == "08006" // connection_failure
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "reset by peer")
}

// recreateSlot drops the existing slot and creates a new one
func (s *Stream) recreateSlot() error {
	logger.DefaultLogger.Warnw("Recreating replication slot", "slot_name", s.slotName)

	// Drop the existing slot
	err := DropReplicationSlot(context.Background(), s.pgConn, s.slotName, DropReplicationSlotOptions{Wait: true})
	if err != nil {
		// Ignore errors if slot doesn't exist
		if !strings.Contains(err.Error(), "does not exist") {
			logger.DefaultLogger.Warnw("Failed to drop replication slot (may not exist)", "error", err.Error())
		}
	}

	// Create new slot
	createResult, err := CreateReplicationSlot(context.Background(), s.pgConn, s.slotName, "wal2json",
		CreateReplicationSlotOptions{
			Temporary:      false,
			SnapshotAction: "", // No snapshot needed for recreation
		})
	if err != nil {
		return fmt.Errorf("failed to recreate replication slot: %w", err)
	}

	// Update LSN to current position
	sysident, err := pglogrepl.IdentifySystem(context.Background(), s.pgConn)
	if err != nil {
		return fmt.Errorf("failed to identify system after slot recreation: %w", err)
	}

	s.lsnrestart = sysident.XLogPos
	s.clientXLogPos = sysident.XLogPos

	logger.DefaultLogger.Infow("Replication slot recreated successfully",
		"slot_name", s.slotName,
		"new_lsn", s.lsnrestart.String(),
		"consistent_point", createResult.ConsistentPoint)

	return nil
}

// startLr starts logical replication with error recovery
func (s *Stream) startLr() error {
	return s.startLrWithRetry(s.conf.MaxRetries)
}

// startLrWithRetry attempts to start logical replication with retry logic
func (s *Stream) startLrWithRetry(maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			logger.DefaultLogger.Infow("Retrying to start replication",
				"attempt", attempt,
				"max_retries", maxRetries,
				"delay", s.conf.RetryDelay)
			time.Sleep(s.conf.RetryDelay)
		}

		// Check slot status before attempting to start
		status, err := s.checkSlotStatus()
		if err != nil {
			logger.DefaultLogger.Warnw("Failed to check slot status", "error", err.Error())
		} else if status != nil {
			logger.DefaultLogger.Infow("Replication slot status",
				"exists", status.Exists,
				"active", status.Active,
				"wal_status", status.WALStatus,
				"restart_lsn", status.RestartLSN)

			// If WAL is lost or unreserved, we need to recreate the slot
			if status.WALStatus == "lost" || status.WALStatus == "unreserved" {
				logger.DefaultLogger.Warnw("WAL segments for replication slot are no longer available",
					"wal_status", status.WALStatus)
				if s.conf.SlotInvalidAutoRecreate {
					if err := s.recreateSlot(); err != nil {
						lastErr = err
						continue
					}
				} else {
					return fmt.Errorf("%w: WAL status is %s", ErrSlotInvalid, status.WALStatus)
				}
			}
		}

		// Attempt to start replication
		err = pglogrepl.StartReplication(
			context.Background(),
			s.pgConn,
			s.slotName,
			s.lsnrestart,
			pglogrepl.StartReplicationOptions{PluginArgs: pluginArguments},
		)

		if err == nil {
			logger.DefaultLogger.Infow("Started logical replication on slot",
				"slot_name", s.slotName,
				"lsn", s.lsnrestart.String())
			s.retryCount = 0
			s.lastHeartbeat = time.Now()
			return nil
		}

		lastErr = err
		logger.DefaultLogger.Errorw("Failed to start replication",
			"error", err.Error(),
			"attempt", attempt+1,
			"max_retries", maxRetries)

		// Handle SQLSTATE 55000 - slot invalid
		if isSlotInvalidError(err) {
			logger.DefaultLogger.Warnw("Replication slot is invalid (SQLSTATE 55000)",
				"slot_name", s.slotName,
				"auto_recreate", s.conf.SlotInvalidAutoRecreate)

			if s.conf.SlotInvalidAutoRecreate {
				if recreateErr := s.recreateSlot(); recreateErr != nil {
					logger.DefaultLogger.Errorw("Failed to recreate slot", "error", recreateErr.Error())
					lastErr = recreateErr
					continue
				}
				// Slot recreated, retry starting replication
				continue
			}
			return fmt.Errorf("%w: %v", ErrSlotInvalid, err)
		}

		// Handle connection errors
		if isConnectionError(err) {
			logger.DefaultLogger.Warnw("Connection error detected, attempting to reconnect")
			if reconnectErr := s.reconnect(); reconnectErr != nil {
				logger.DefaultLogger.Errorw("Failed to reconnect", "error", reconnectErr.Error())
				lastErr = reconnectErr
				continue
			}
			continue
		}
	}

	return fmt.Errorf("%w after %d attempts: %v", ErrMaxRetriesReached, maxRetries+1, lastErr)
}

// reconnect attempts to re-establish the PostgreSQL connection
func (s *Stream) reconnect() error {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	// Close existing connection
	if s.pgConn != nil {
		_ = s.pgConn.Close(context.Background())
	}

	// Create new connection
	newConn, err := pgconn.ConnectConfig(context.Background(), &s.dbConfig)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	s.pgConn = newConn
	logger.DefaultLogger.Infow("Successfully reconnected to PostgreSQL")
	return nil
}

func (s *Stream) AckLSN(lsn string) error {
	var err error
	s.clientXLogPos, err = pglogrepl.ParseLSN(lsn)
	if err != nil {
		logger.DefaultLogger.Errorf("Failed to parse LSN for Acknowledge %s", err.Error())
		return err
	}

	err = pglogrepl.SendStandbyStatusUpdate(context.Background(), s.pgConn, pglogrepl.StandbyStatusUpdate{
		WALApplyPosition: s.clientXLogPos,
		WALWritePosition: s.clientXLogPos,
		ReplyRequested:   true,
	})

	if err != nil {
		logger.DefaultLogger.Errorf("SendStandbyStatusUpdate failed: %s", err.Error())
		return err
	}
	logger.DefaultLogger.Debugf("Sent Standby status message at LSN#%s", s.clientXLogPos.String())
	s.nextStandbyMessageDeadline = time.Now().Add(s.standbyMessageTimeout)
	return nil
}

// sendHeartbeat sends a standby status update to PostgreSQL
func (s *Stream) sendHeartbeat() error {
	s.connMu.RLock()
	defer s.connMu.RUnlock()

	err := pglogrepl.SendStandbyStatusUpdate(context.Background(), s.pgConn, pglogrepl.StandbyStatusUpdate{
		WALWritePosition: s.clientXLogPos,
	})

	if err != nil {
		return err
	}

	s.lastHeartbeat = time.Now()
	logger.DefaultLogger.Debugf("Sent Standby status message at LSN#%s", s.clientXLogPos.String())
	return nil
}

// handleStreamError handles errors during message streaming with recovery logic
func (s *Stream) handleStreamError(err error) bool {
	if err == nil {
		return false
	}

	// Check if we're stopped or context is cancelled
	if errors.Is(err, context.Canceled) || s.stopped {
		logger.DefaultLogger.Warn("Service was interrupted...stop reading from replication slot")
		return true // Signal to exit the loop
	}

	// Timeout is normal, continue
	if pgconn.Timeout(err) {
		return false
	}

	// Handle slot invalid error (SQLSTATE 55000)
	if isSlotInvalidError(err) {
		logger.DefaultLogger.Errorw("Replication slot became invalid during streaming",
			"error", err.Error(),
			"slot_name", s.slotName)

		if s.conf.SlotInvalidAutoRecreate {
			s.retryCount++
			if s.retryCount > s.conf.MaxRetries {
				logger.DefaultLogger.Fatalw("Max retries reached for slot recovery",
					"retries", s.retryCount,
					"max_retries", s.conf.MaxRetries)
				return true
			}

			logger.DefaultLogger.Infow("Attempting to recover replication slot",
				"attempt", s.retryCount,
				"max_retries", s.conf.MaxRetries)

			// Reconnect and recreate slot
			if reconnectErr := s.reconnect(); reconnectErr != nil {
				logger.DefaultLogger.Errorw("Failed to reconnect", "error", reconnectErr.Error())
				time.Sleep(s.conf.RetryDelay)
				return false // Retry
			}

			if recreateErr := s.recreateSlot(); recreateErr != nil {
				logger.DefaultLogger.Errorw("Failed to recreate slot", "error", recreateErr.Error())
				time.Sleep(s.conf.RetryDelay)
				return false // Retry
			}

			// Restart replication
			if startErr := s.startLrWithRetry(1); startErr != nil {
				logger.DefaultLogger.Errorw("Failed to restart replication", "error", startErr.Error())
				time.Sleep(s.conf.RetryDelay)
				return false // Retry
			}

			s.retryCount = 0
			return false // Continue streaming
		}

		logger.DefaultLogger.Fatalw("Replication slot invalid and auto-recreate is disabled",
			"error", err.Error())
		return true
	}

	// Handle connection errors
	if isConnectionError(err) {
		logger.DefaultLogger.Warnw("Connection error during streaming, attempting recovery",
			"error", err.Error())

		s.retryCount++
		if s.retryCount > s.conf.MaxRetries {
			logger.DefaultLogger.Fatalw("Max retries reached for connection recovery",
				"retries", s.retryCount)
			return true
		}

		time.Sleep(s.conf.RetryDelay)

		if reconnectErr := s.reconnect(); reconnectErr != nil {
			logger.DefaultLogger.Errorw("Failed to reconnect", "error", reconnectErr.Error())
			return false // Retry
		}

		if startErr := s.startLrWithRetry(1); startErr != nil {
			logger.DefaultLogger.Errorw("Failed to restart replication after reconnect",
				"error", startErr.Error())
			return false // Retry
		}

		s.retryCount = 0
		return false // Continue streaming
	}

	// Unknown error - log and continue or fatal based on severity
	logger.DefaultLogger.Errorw("Error receiving message from PostgreSQL",
		"error", err.Error())

	// For unknown errors, we'll try to continue but increment retry count
	s.retryCount++
	if s.retryCount > s.conf.MaxRetries {
		logger.DefaultLogger.Fatalw("Max retries reached for unknown errors",
			"retries", s.retryCount,
			"last_error", err.Error())
		return true
	}

	return false
}

func (s *Stream) streamMessagesAsync() {
	for {
		select {
		case <-s.streamCtx.Done():
			logger.DefaultLogger.Warn("Stream was cancelled...exiting...")
			return
		default:
			// Send heartbeat if deadline has passed
			if time.Now().After(s.nextStandbyMessageDeadline) {
				if err := s.sendHeartbeat(); err != nil {
					logger.DefaultLogger.Warnw("Failed to send heartbeat",
						"error", err.Error(),
						"will_retry", true)

					// Handle heartbeat error with recovery
					if isConnectionError(err) || isSlotInvalidError(err) {
						if shouldExit := s.handleStreamError(err); shouldExit {
							return
						}
						continue
					}
				}
				s.nextStandbyMessageDeadline = time.Now().Add(s.standbyMessageTimeout)
			}

			ctx, cancel := context.WithDeadline(context.Background(), s.nextStandbyMessageDeadline)
			s.connMu.RLock()
			rawMsg, err := s.pgConn.ReceiveMessage(ctx)
			s.connMu.RUnlock()
			s.standbyCtxCancel = cancel

			// Handle errors with recovery logic
			if err != nil {
				cancel()
				if shouldExit := s.handleStreamError(err); shouldExit {
					return
				}
				continue
			}

			// Reset retry count on successful message
			s.retryCount = 0

			// Handle ErrorResponse from PostgreSQL
			if errMsg, ok := rawMsg.(*pgproto3.ErrorResponse); ok {
				logger.DefaultLogger.Errorw("Received error from PostgreSQL WAL",
					"severity", errMsg.Severity,
					"code", errMsg.Code,
					"message", errMsg.Message,
					"detail", errMsg.Detail)

				// Check if it's a slot invalid error
				if errMsg.Code == "55000" {
					syntheticErr := fmt.Errorf("SQLSTATE %s: %s", errMsg.Code, errMsg.Message)
					if shouldExit := s.handleStreamError(syntheticErr); shouldExit {
						return
					}
					continue
				}

				logger.DefaultLogger.Fatalw("Fatal PostgreSQL error", "error", errMsg)
				return
			}

			msg, ok := rawMsg.(*pgproto3.CopyData)
			if !ok {
				logger.DefaultLogger.Warnf("Received unexpected message: %T\n", rawMsg)
				continue
			}

			switch msg.Data[0] {
			case pglogrepl.PrimaryKeepaliveMessageByteID:
				pkm, err := pglogrepl.ParsePrimaryKeepaliveMessage(msg.Data[1:])
				if err != nil {
					logger.DefaultLogger.Errorw("ParsePrimaryKeepaliveMessage failed",
						"error", err.Error())
					continue
				}

				if pkm.ReplyRequested {
					s.nextStandbyMessageDeadline = time.Time{}
				}

			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
				if err != nil {
					logger.DefaultLogger.Errorw("ParseXLogData failed", "error", err.Error())
					continue
				}
				clientXLogPos := xld.WALStart + pglogrepl.LSN(len(xld.WALData))
				var changes Wal2JsonChanges
				bytesData := bytes.NewReader(xld.WALData)
				if err := json.NewDecoder(bytesData).Decode(&changes); err != nil {
					logger.DefaultLogger.Errorw("Failed to decode WAL data",
						"error", err.Error(),
						"lsn", clientXLogPos.String())
					continue
				}

				if len(changes.Changes) == 0 {
					if s.autoAck {
						_ = s.AckLSN(clientXLogPos.String())
					}
				} else {
					s.changeFilter.FilterChange(clientXLogPos.String(), changes, func(change Wal2JsonChanges) {
						s.messages <- change
					})
				}
			}
		}
	}
}
func (s *Stream) processSnapshot() {
	snapshotter, err := NewSnapshotter(s.dbConfig, s.snapshotName)
	if err != nil {
		logger.DefaultLogger.Errorf("Failed to create database snapshot: %v", err.Error())
		s.cleanUpOnFailure()
		os.Exit(1)
	}
	if err = snapshotter.Prepare(); err != nil {
		logger.DefaultLogger.Errorf("Failed to prepare database snapshot: %v", err.Error())
		s.cleanUpOnFailure()
		os.Exit(1)
	}
	defer func() {
		snapshotter.ReleaseSnapshot()
		snapshotter.CloseConn()
	}()

	for _, table := range s.tableNames {
		logger.DefaultLogger.Info("Processing snapshot for table", "table", table)

		var (
			avgRowSizeBytes sql.NullInt64
			offset          = int64(0)
		)
		if s.conf.SnapshotOffset > 0 {
			offset = s.conf.SnapshotOffset
		}
		avgRowSizeBytes = snapshotter.FindAvgRowSize(table)

		memUsage := utils.GetAvailableMemory()
		batchSize := int64(snapshotter.CalculateBatchSize(memUsage, uint64(avgRowSizeBytes.Int64)))
		logger.DefaultLogger.Info("Querying snapshot", "batch_side", batchSize, "available_memory", memUsage, "avg_row_size", avgRowSizeBytes.Int64)

		tablePk, err := s.getPrimaryKeyColumn(table)
		if err != nil {
			panic(err)
		}

		for {
			var snapshotRows *sql.Rows
			if snapshotRows, err = snapshotter.QuerySnapshotData(table, tablePk, batchSize, offset); err != nil {
				log.Fatalf("Can't query snapshot data %v", err)
			}

			columnTypes, err := snapshotRows.ColumnTypes()
			var columnTypesString = make([]string, len(columnTypes))
			columnNames, err := snapshotRows.Columns()
			for i, _ := range columnNames {
				columnTypesString[i] = columnTypes[i].DatabaseTypeName()
			}

			if err != nil {
				panic(err)
			}

			count := len(columnTypes)
			var rowsCount = 0
			for snapshotRows.Next() {
				rowsCount += 1
				scanArgs := make([]interface{}, count)
				for i, v := range columnTypes {
					switch v.DatabaseTypeName() {
					case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
						scanArgs[i] = new(sql.NullString)
						break
					case "BOOL":
						scanArgs[i] = new(sql.NullBool)
						break
					case "INT4":
						scanArgs[i] = new(sql.NullInt64)
						break
					default:
						scanArgs[i] = new(sql.NullString)
					}
				}

				err := snapshotRows.Scan(scanArgs...)

				if err != nil {
					panic(err)
				}

				var columnValues = make([]interface{}, len(columnTypes))
				for i, _ := range columnTypes {
					if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
						columnValues[i] = z.Bool
						continue
					}
					if z, ok := (scanArgs[i]).(*sql.NullString); ok {
						columnValues[i] = z.String
						continue
					}
					if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
						columnValues[i] = z.Int64
						continue
					}
					if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
						columnValues[i] = z.Float64
						continue
					}
					if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
						columnValues[i] = z.Int32
						continue
					}

					columnValues[i] = scanArgs[i]
				}

				var snapshotChanges []Wal2JsonChange
				snapshotChanges = append(snapshotChanges, Wal2JsonChange{
					Kind:         "insert",
					Schema:       s.schema,
					Table:        table,
					ColumnNames:  columnNames,
					ColumnValues: columnValues,
				})
				var lsn *string
				snapshotChangePacket := Wal2JsonChanges{
					Lsn:     lsn,
					Changes: snapshotChanges,
				}

				s.snapshotMessages <- snapshotChangePacket
			}

			offset += batchSize

			if batchSize != int64(rowsCount) {
				break
			}
		}

	}
	if err := s.startLr(); err != nil {
		logger.DefaultLogger.Errorf("Failed to start logical replication after snapshot: %v", err)
		s.cleanUpOnFailure()
		os.Exit(1)
	}
	go s.streamMessagesAsync()
}

func (s *Stream) OnMessage(callback OnMessage) {
	for {
		select {
		case snapshotMessage := <-s.snapshotMessages:
			callback(snapshotMessage)
		case message := <-s.messages:
			callback(message)
		case <-s.streamCtx.Done():
			return
		}
	}
}

func (s *Stream) SnapshotMessageC() chan Wal2JsonChanges {
	return s.snapshotMessages
}

func (s *Stream) LrMessageC() chan Wal2JsonChanges {
	return s.messages
}

// cleanUpOnFailure drops replication slot and publication if database snapshotting was failed for any reason
func (s *Stream) cleanUpOnFailure() {
	logger.DefaultLogger.Warn("Cleaning up resources on accident.", "replication-slot", s.slotName)
	err := DropReplicationSlot(context.Background(), s.pgConn, s.slotName, DropReplicationSlotOptions{Wait: true})
	if err != nil {
		logger.DefaultLogger.Errorf("Failed to drop replication slot: %s", err.Error())
	}
	s.pgConn.Close(context.TODO())
}

func (s *Stream) getPrimaryKeyColumn(tableName string) (string, error) {
	q := fmt.Sprintf(`
		SELECT a.attname
		FROM   pg_index i
		JOIN   pg_attribute a ON a.attrelid = i.indrelid
							 AND a.attnum = ANY(i.indkey)
		WHERE  i.indrelid = '%s'::regclass
		AND    i.indisprimary;
	`, tableName)

	reader := s.pgConn.Exec(context.Background(), q)
	data, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	pkResultRow := data[0].Rows[0]
	pkColName := string(pkResultRow[0])
	return pkColName, nil
}

func (s *Stream) Stop() error {
	s.m.Lock()
	s.stopped = true
	s.m.Unlock()

	if s.pgConn != nil {
		if s.streamCtx != nil {
			s.streamCancel()
			s.standbyCtxCancel()
		}
		return s.pgConn.Close(context.Background())
	}

	return nil
}
