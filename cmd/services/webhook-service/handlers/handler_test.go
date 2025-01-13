package handlers

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/tests"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"testing"
)

var webhookHandler *WebhookHandler
var ctx = context.Background()

func init2() {
	err := tests.LoadTestConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	//redis, err := dbtool.CreateRedisConnection(nil)

	webhookHandler = NewWebhookHandler(psql)
}

func map2jsonb(m map[string]interface{}) *structpb.Struct {
	v, _ := structpb.NewStruct(m)
	return v
}

func TestWebhookHandler_InsertLogs(t *testing.T) {
	// Test GetUserByID
	events := &models.WebhookEvents{
		Events: []*models.WebhookEvent{
			{
				Kind: "facebook",
				Payload: map2jsonb(map[string]interface{}{
					"a": 1,
					"b": 2,
				}),
				Meta: map2jsonb(map[string]interface{}{
					"signature": "123456",
				}),
				MessageUuid: "042a6d5b-361f-49b4-92a0-2690b14e563f",
			},
			{
				Kind: "facebook2",
				Payload: map2jsonb(map[string]interface{}{
					"a": 4,
					"b": 25,
				}),
				Meta: map2jsonb(map[string]interface{}{
					"signature": "1234564",
				}),
				MessageUuid: watermill.NewUUID(),
			},
		},
	}
	resp := models.SQLResult{}
	err := webhookHandler.InsertLogs(ctx, events, &resp)
	if err != nil {
		t.Fatal("Failed to insert logs", err)
	}
	t.Log("logs: ", resp.LastInsertIds)
}

func TestName(t *testing.T) {
	a :=
		fmt.Println(a)
}
