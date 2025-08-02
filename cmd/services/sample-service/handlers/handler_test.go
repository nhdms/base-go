package handlers

import (
	"context"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/tests"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

var sampleHandler *SampleHandler
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

	sampleHandler = NewSampleHandler(psql)
}

func map2jsonb(m map[string]interface{}) *structpb.Struct {
	v, _ := structpb.NewStruct(m)
	return v
}

func TestSampleHandler_CreateSample(t *testing.T) {
	event := &models.SampleEvent{
		EventType: "test",
		UserId:    1,
		Payload: map2jsonb(map[string]interface{}{
			"test": "data",
		}),
		CreatedAt: timestamppb.Now(),
	}

	resp, err := sampleHandler.CreateSample(ctx, event)
	if err != nil {
		t.Fatal("Failed to create event:", err)
	}

	if !resp.Success {
		t.Fatal("Expected success response")
	}
}
