package templates

const ServiceMainTemplate = `package main

import (
    "gitlab.com/a7923/athena-go/cmd/services/{{.Name}}/handlers"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/dbtool"
    "gitlab.com/a7923/athena-go/pkg/logger"
    "gitlab.com/a7923/athena-go/proto/exmsg/services"
)

func main() {
    svc := app.NewGRPCService()
    psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
    if err != nil {
        logger.AthenaLogger.Fatal("Failed to connect to database: ", err)
    }

    grpcSvc := handlers.New{{.Handler}}Handler(psql)
    err = services.Register{{.Handler}}ServiceHandler(svc.Server(), grpcSvc)
    if err != nil {
        logger.AthenaLogger.Fatal(err)
    }

    err = svc.Run()
    if err != nil {
        logger.AthenaLogger.Fatal(err)
    }
}`

const ServiceHandlerTemplate = `package handlers

import (
	"context"
	"gitlab.com/a7923/athena-go/pkg/dbtool"
	"gitlab.com/a7923/athena-go/proto/exmsg/services"
)

type EventHandler struct {
	db *dbtool.ConnectionManager
}

func (e EventHandler) GetEvents(ctx context.Context, request *services.EventRequest, response *services.EventResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewEventHandler(db *dbtool.ConnectionManager) *EventHandler {
	return &EventHandler{db: db}
}
`

const ServiceTableTemplate = `package tables

import "gitlab.com/a7923/athena-go/pkg/dbtool"

func Get{{.Handler}}Table() *dbtool.Table {
	return &dbtool.Table{
		Name:      "{{.TableName}}",
		AIColumns: []string{"id"},
		ColumnMapper: map[string]string{
		},
		IgnoreColumns: []string{},
		DefaultAlias:  "u",
	}
}
`

const ServiceHandlerTestTemplate = `package handlers

import (
    "context"
    "testing"
    "gitlab.com/a7923/athena-go/pkg/dbtool"
    "gitlab.com/a7923/athena-go/proto/exmsg/models"
    "gitlab.com/a7923/athena-go/tests"
    "google.golang.org/protobuf/types/known/structpb"
    "google.golang.org/protobuf/types/known/timestamppb"
)

var {{.ServiceName}}Handler *{{.Handler}}Handler
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

    {{.ServiceName}}Handler = New{{.Handler}}Handler(psql)
}

func map2jsonb(m map[string]interface{}) *structpb.Struct {
    v, _ := structpb.NewStruct(m)
    return v
}

func Test{{.Handler}}Handler_Create{{.Handler}}(t *testing.T) {
    event := &models.{{.Handler}}Event{
        EventType: "test",
        UserId: 1,
        Payload: map2jsonb(map[string]interface{}{
            "test": "data",
        }),
        CreatedAt: timestamppb.Now(),
    }

    resp, err := {{.ServiceName}}Handler.Create{{.Handler}}(ctx, event)
    if err != nil {
        t.Fatal("Failed to create event:", err)
    }

    if !resp.Success {
        t.Fatal("Expected success response")
    }
}`

const ServiceCICDTemplate = `{
    "app_type": "service",
    "cmd_bin_dir": "cmd/services/{{.Name}}",
    "service_name": "{{.ServiceName}}",
    "port": {{.Port}},
    "config_remote_keys": [
        "database/postgres.toml",
        "services/{{.ServiceName}}.toml"
    ]
}`
