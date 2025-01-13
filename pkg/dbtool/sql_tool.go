package dbtool

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jmoiron/sqlx"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/spf13/viper"
	metadata2 "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"time"
)

const (
	KindSelect = "select"
	KindUpdate = "update"
	KindDelete = "delete"
	KindInsert = "insert"
)

type SQLTool struct {
	db          *sqlx.DB
	table       *Table
	columns     []string
	debug       bool
	ctx         context.Context
	canSetCache bool
	kind        string
	column2kind map[string]reflect.Kind
	column2type map[string]reflect.Type
	column2name map[string]string
	tx          *sqlx.Tx
}

func New(ctx context.Context, db *sqlx.DB, table *Table, model interface{}, kind string) *SQLTool {
	// parse columns
	s := &SQLTool{
		ctx:   ctx,
		db:    db,
		table: table,
		kind:  kind,
		debug: viper.GetBool("sql.debug") || utils.IsTestMode(),
	}
	s.prepare(ctx, table, model, kind)
	return s
}

func NewSelect(ctx context.Context, db *sqlx.DB, table *Table, model interface{}) *SQLTool {
	return New(ctx, db, table, model, KindSelect)
}

func NewInsert(ctx context.Context, db *sqlx.DB, table *Table, model interface{}) *SQLTool {
	return New(ctx, db, table, model, KindInsert)
}

func NewUpdate(ctx context.Context, db *sqlx.DB, table *Table, model interface{}) *SQLTool {
	return New(ctx, db, table, model, KindUpdate)
}

func NewDelete(ctx context.Context, db *sqlx.DB, table *Table, model interface{}) *SQLTool {
	return New(ctx, db, table, model, KindDelete)
}

func NewTransaction(ctx context.Context, db *sqlx.DB) (*SQLTool, error) {
	s := &SQLTool{
		ctx:   ctx,
		db:    db,
		debug: viper.GetBool("sql.debug") || utils.IsTestMode(),
	}
	var err error
	s.tx, err = s.db.BeginTxx(ctx, &sql.TxOptions{})
	return s, err
}

func (s *SQLTool) RollbackTransactions() error {
	return s.tx.Rollback()
}

func (s *SQLTool) CommitTransactions() error {
	return s.tx.Commit()
}

func (s *SQLTool) GetTable(alias string) string {
	if len(alias) == 0 {
		return s.table.Name
	}

	return s.table.Name + " " + alias
}

func (s *SQLTool) GetQueryColumnList(alias string) []string {
	if len(alias) == 0 {
		return s.columns
	}

	columns := make([]string, len(s.columns))
	for i, column := range s.columns {
		columns[i] = alias + "." + column
	}
	return columns
}

func (s *SQLTool) GetFilledValues(item interface{}) []interface{} {
	values := make([]interface{}, len(s.columns))
	val := reflect.ValueOf(item)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return values
	}

	for i, column := range s.columns {
		field, found := s.column2name[column]
		if !found {
			if defaultVal, hasDefault := s.table.NotNullColumns[column]; hasDefault {
				// Handle function type defaults (like time.Now)
				if fn, ok := defaultVal.(func() interface{}); ok {
					values[i] = fn()
				} else {
					values[i] = defaultVal
				}
			} else {
				values[i] = nil
			}
			continue
		}

		fieldValue := val.FieldByName(field)

		// Handle nil pointers for not null columns
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			if defaultVal, hasDefault := s.table.NotNullColumns[column]; hasDefault {
				if fn, ok := defaultVal.(func() interface{}); ok {
					values[i] = fn()
				} else {
					values[i] = defaultVal
				}
			} else {
				values[i] = nil
			}
			continue
		}

		// Get actual value based on the field's kind
		switch fieldValue.Kind() {
		case reflect.Ptr:
			fieldValue = fieldValue.Elem()

			// Special handling for timestamp
			if ts, ok := fieldValue.Interface().(timestamppb.Timestamp); ok {
				values[i] = time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
				continue
			}

			if vl, ok := fieldValue.Interface().(structpb.Struct); ok {
				str, _ := vl.MarshalJSON()
				values[i] = string(str)
				continue
			}

			values[i] = fieldValue.Interface()

		case reflect.Struct:
			// Special handling for timestamp
			if ts, ok := fieldValue.Interface().(*timestamp.Timestamp); ok {
				values[i] = time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
				continue
			}

			// Check for zero value in not null column
			if defaultVal, hasDefault := s.table.NotNullColumns[column]; hasDefault {
				if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
					if fn, ok := defaultVal.(func() interface{}); ok {
						values[i] = fn()
					} else {
						values[i] = defaultVal
					}
					continue
				}
			}

			values[i] = fieldValue.Interface()

		default:
			// Handle basic types
			isZero := false
			switch fieldValue.Kind() {
			case reflect.Float32, reflect.Float64:
				isZero = fieldValue.Float() == 0
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				isZero = fieldValue.Int() == 0
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				isZero = fieldValue.Uint() == 0
			case reflect.Bool:
				isZero = !fieldValue.Bool()
			default:
				isZero = reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface())
			}

			if isZero {
				if defaultVal, hasDefault := s.table.NotNullColumns[column]; hasDefault {
					if fn, ok := defaultVal.(func() interface{}); ok {
						values[i] = fn()
					} else {
						values[i] = defaultVal
					}
					continue
				}
				values[i] = nil
				continue
			}

			values[i] = fieldValue.Interface()
		}
	}

	return values
}

func (s *SQLTool) Get(ctx context.Context, dest interface{}, qb squirrel.SelectBuilder) error {
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if s.debug {
		logger.DefaultLogger.Debugw("Executing query ", "query", query, "args", args)
	}

	var rows *sqlx.Row
	if s.tx != nil {
		rows = s.tx.QueryRowxContext(ctx, query, args...)
	} else {
		rows = s.db.QueryRowxContext(ctx, query, args...)
	}

	return ScanRow(rows, dest)
}

func (s *SQLTool) Select(ctx context.Context, dest interface{}, qb squirrel.SelectBuilder) error {
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if s.debug {
		logger.DefaultLogger.Debugw("Executing query ", "query", query, "args", args)
	}

	var rows *sqlx.Rows
	if s.tx != nil {
		rows, err = s.tx.QueryxContext(ctx, query, args...)
	} else {
		rows, err = s.db.QueryxContext(ctx, query, args...)
	}

	if err != nil {
		return err
	}

	return ScanAll(rows, dest)
}

func (s *SQLTool) parseColumns(model interface{}) {
	column2kind := make(map[string]reflect.Kind)
	column2type := make(map[string]reflect.Type)
	column2name := make(map[string]string)
	columns := make([]string, 0)
	ignoreColumns := s.getIgnoreColumns()

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields and protobuf internal fields
		if !field.IsExported() || field.Name == "state" || field.Name == "sizeCache" || field.Name == "unknownFields" {
			continue
		}

		// Get column name from json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		columnName := strings.Split(jsonTag, ",")[0]

		if ignoreColumns[columnName] {
			continue
		}
		if v, ok := s.table.ColumnMapper[columnName]; ok && len(v) > 0 {
			columnName = v
		}
		// Add to mappings
		columns = append(columns, columnName)
		column2kind[columnName] = field.Type.Kind()
		column2type[columnName] = field.Type
		column2name[columnName] = field.Name
	}

	md, ok := metadata2.FromIncomingContext(s.ctx)
	if !ok {
		md = make(metadata2.MD)
	}

	s.columns = columns
	cols := md.Get(MetadataKeySelectedFields)
	if len(cols) > 0 {
		s.columns = cols
	}

	s.columns = s.filterColumns(s.columns)
	s.column2kind = column2kind
	s.column2type = column2type
	s.column2name = column2name
}

func (s *SQLTool) CanSetCache() bool {
	return s.canSetCache
}

func (s *SQLTool) Insert(ctx context.Context, qb squirrel.InsertBuilder) (*models.SQLResult, error) {
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	qb = qb.Suffix("RETURNING id")
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	if s.debug {
		logger.DefaultLogger.Debugw("Executing query ", "query", query, "args", args)
	}

	v := &models.SQLResult{
		LastInsertIds: make([]int64, 0),
		RowsAffected:  0,
	}

	if s.tx != nil {
		err = s.tx.SelectContext(ctx, &v.LastInsertIds, query, args...)
	} else {
		err = s.db.SelectContext(ctx, &v.LastInsertIds, query, args...)
	}

	if err != nil {
		return nil, err
	}

	return v, err
}

func (s *SQLTool) Update(ctx context.Context, qb squirrel.UpdateBuilder) (*models.SQLResult, error) {
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	return s.execContext(ctx, qb)
}

func (s *SQLTool) Delete(ctx context.Context, qb squirrel.DeleteBuilder) (*models.SQLResult, error) {
	qb = qb.PlaceholderFormat(squirrel.Dollar)
	return s.execContext(ctx, qb)
}

func (s *SQLTool) execContext(ctx context.Context, qb squirrel.Sqlizer) (*models.SQLResult, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	if s.debug {
		logger.DefaultLogger.Debugw("Executing query ", "query", query, "args", args)
	}

	v := &models.SQLResult{
		RowsAffected: 0,
	}

	var res sql.Result
	if s.tx != nil {
		res, err = s.tx.ExecContext(ctx, query, args...)
	} else {
		res, err = s.db.ExecContext(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	v.RowsAffected, _ = res.RowsAffected()
	return v, nil
}

func (s *SQLTool) getIgnoreColumns() map[string]bool {
	r := make(map[string]bool)
	if s.kind != KindInsert {
		return r
	}

	for _, col := range s.table.AIColumns {
		r[col] = true
	}
	return r
}

func (s *SQLTool) PrepareInsert(ctx context.Context, table *Table, model interface{}) {
	s.prepare(ctx, table, model, KindInsert)
}

func (s *SQLTool) PrepareUpdate(ctx context.Context, table *Table, model interface{}) {
	s.prepare(ctx, table, model, KindUpdate)
}

func (s *SQLTool) PrepareDelete(ctx context.Context, table *Table, model interface{}) {
	s.prepare(ctx, table, model, KindDelete)
}

func (s *SQLTool) PrepareSelect(ctx context.Context, table *Table, model interface{}) {
	s.prepare(ctx, table, model, KindSelect)
}

func (s *SQLTool) prepare(ctx context.Context, table *Table, model interface{}, kind string) {
	s.kind = kind
	s.table = table
	s.defineDefaultValues()
	s.parseColumns(model)
}

func (s *SQLTool) GetUpdateMap(dest interface{}, updateFields ...string) map[string]interface{} {
	m := make(map[string]interface{})
	val := s.GetFilledValues(dest)

	cols := s.columns
	if len(updateFields) > 0 {
		cols = updateFields
	}
	for i, c := range cols {
		m[c] = val[i]
	}

	return m
}

func (s *SQLTool) defineDefaultValues() {
	if s.table.NotNullColumns == nil {
		s.table.NotNullColumns = make(map[string]interface{})
	}

	defaultDateTimeColumns := []string{
		"created_at",
		"updated_at",
	}

	for _, c := range defaultDateTimeColumns {
		if _, ok := s.table.NotNullColumns[c]; !ok {
			s.table.NotNullColumns[c] = func() interface{} { return time.Now() }
		}
	}
}

func (s *SQLTool) filterColumns(columns []string) []string {
	ignore := make(map[string]bool)
	for _, c := range s.table.IgnoreColumns {
		ignore[c] = true
	}
	cols := make([]string, 0)
	for _, c := range columns {
		if _, ok := ignore[c]; !ok {
			cols = append(cols, c)
		}
	}
	return cols
}
