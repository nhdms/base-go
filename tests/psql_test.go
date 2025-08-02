package tests

import (
	"fmt"
	"github.com/nhdms/base-go/pkg/replication/pglogicalstream"
	"github.com/nhdms/base-go/pkg/utils/pg_converter"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"log"
	"testing"
	"time"
)

func TestParseWalMessage(t *testing.T) {
	// Example WAL message
	walMsg := pglogicalstream.Wal2JsonChanges{
		Changes: []pglogicalstream.Wal2JsonChange{
			{
				Kind:         "insert",
				Schema:       "public",
				Table:        "users",
				ColumnNames:  []string{"id", "customer_name", "sale_id", "created_at", "source_detail"},
				ColumnTypes:  []string{"bigint", "text", "integer", "timestamp with time zone", "jsonb"},
				ColumnValues: []interface{}{int64(1), "John Doe", int32(100), "2024-11-18 11:29:46.781285+00", `{"a": 1, "b": 2}`},
			},
		},
	}

	// Create a new user message
	user := &models.Order{}
	// Create converter with UTC timezone
	protoConverter := pg_converter.NewProtoConverter(user, time.UTC)

	// Convert WAL message to protobuf
	for _, change := range walMsg.Changes {
		err := protoConverter.ConvertToStruct(
			change.ColumnNames,
			change.ColumnTypes,
			change.ColumnValues,
			user,
		)
		if err != nil {
			log.Fatalf("Error converting to proto: %v", err)
		}
	}

	fmt.Printf("Converted user: %+v\n", user)
}
