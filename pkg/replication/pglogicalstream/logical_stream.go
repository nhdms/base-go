package pglogicalstream

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

var pluginArguments = []string{"\"pretty-print\" 'true'"}

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

	stream.startLr()
	go stream.streamMessagesAsync()

	return stream, err
}

func (s *Stream) startLr() {
	var err error
	err = pglogrepl.StartReplication(context.Background(), s.pgConn, s.slotName, s.lsnrestart, pglogrepl.StartReplicationOptions{PluginArgs: pluginArguments})
	if err != nil {
		logger.DefaultLogger.Fatalf("Starting replication slot failed: %s", err.Error())
	}
	logger.DefaultLogger.Infow("Started logical replication on slot", "slot-name", s.slotName)
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

func (s *Stream) streamMessagesAsync() {
	for {
		select {
		case <-s.streamCtx.Done():
			logger.DefaultLogger.Warn("Stream was cancelled...exiting...")
			return
		default:
			if time.Now().After(s.nextStandbyMessageDeadline) {
				var err error
				err = pglogrepl.SendStandbyStatusUpdate(context.Background(), s.pgConn, pglogrepl.StandbyStatusUpdate{
					WALWritePosition: s.clientXLogPos,
				})

				if err != nil {
					logger.DefaultLogger.Fatalf("SendStandbyStatusUpdate failed: %s", err.Error())
				}
				logger.DefaultLogger.Infof("Sent Standby status message at LSN#%s", s.clientXLogPos.String())
				s.nextStandbyMessageDeadline = time.Now().Add(s.standbyMessageTimeout)
			}

			ctx, cancel := context.WithDeadline(context.Background(), s.nextStandbyMessageDeadline)
			rawMsg, err := s.pgConn.ReceiveMessage(ctx)
			s.standbyCtxCancel = cancel

			if err != nil && (errors.Is(err, context.Canceled) || s.stopped) {
				logger.DefaultLogger.Warn("Service was interrpupted....stop reading from replication slot")
				return
			}

			if err != nil {
				if pgconn.Timeout(err) {
					continue
				}

				logger.DefaultLogger.Fatalf("Failed to receive messages from PostgreSQL %s", err.Error())
			}

			if errMsg, ok := rawMsg.(*pgproto3.ErrorResponse); ok {
				logger.DefaultLogger.Fatalf("Received broken Postgres WAL. Error: %+v", errMsg)
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
					logger.DefaultLogger.Fatalf("ParsePrimaryKeepaliveMessage failed: %s", err.Error())
				}

				if pkm.ReplyRequested {
					s.nextStandbyMessageDeadline = time.Time{}
				}

			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
				if err != nil {
					logger.DefaultLogger.Fatalf("ParseXLogData failed: %s", err.Error())
				}
				clientXLogPos := xld.WALStart + pglogrepl.LSN(len(xld.WALData))
				var changes Wal2JsonChanges
				bytesData := bytes.NewReader(xld.WALData)
				if err := json.NewDecoder(bytesData).Decode(&changes); err != nil {
					panic(fmt.Errorf("cant parse change from database to filter it %v", err))
				}

				if len(changes.Changes) == 0 {
					if s.autoAck {
						s.AckLSN(clientXLogPos.String())
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
	s.startLr()
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
