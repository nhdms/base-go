package handlers

import (
    "context"
    "testing"
    "gitlab.com/a7923/athena-go/pkg/dbtool"
    "gitlab.com/a7923/athena-go/proto/exmsg/models"
    "gitlab.com/a7923/athena-go/tests"
    "google.golang.org/protobuf/types/known/structpb"
    "google.golang.org/protobuf/types/known/timestamppb"
)

var eventHandler *EventHandler
var ctx = context.Background()

func init() {
    err := tests.LoadTestConfig()
    if err != nil {
        t.Fatal("Failed to load config", err)
    }

    psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
    if err != nil {
        t.Fatal("Failed to connect to database:", err)
    }

    eventHandler = NewEventHandler(psql)
}

func map2jsonb(m map[string]interface{}) *structpb.Struct {
    v, _ := structpb.NewStruct(m)
    return v
}

func TestEventHandler_CreateEvent(t *testing.T) {
    event := &models.EventEvent{
        EventType: "test",
        UserId: 1,
        Payload: map2jsonb(map[string]interface{}{
            "test": "data",
        }),
        CreatedAt: timestamppb.Now(),
    }

    resp, err := eventHandler.CreateEvent(ctx, event)
    if err != nil {
        t.Fatal("Failed to create event:", err)
    }

    if !resp.Success {
        t.Fatal("Expected success response")
    }
}