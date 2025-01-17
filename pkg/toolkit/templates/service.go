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
    "gitlab.com/a7923/athena-go/proto/exmsg/models"
)

type {{.Handler}}Handler struct {
    db *dbtool.ConnectionManager
}

func New{{.Handler}}Handler(db *dbtool.ConnectionManager) *{{.Handler}}Handler {
    return &{{.Handler}}Handler{db: db}
}

func (h *{{.Handler}}Handler) Create{{.Handler}}(ctx context.Context, event *models.{{.Handler}}Event) (*models.SQLResult, error) {
    sqlTool := dbtool.NewInsert(ctx, h.db.GetConnection(), tables.Get{{.Handler}}Table())
    
    result, err := sqlTool.Insert("{{.TableName}}").
        Columns("event_type", "user_id", "payload", "created_at").
        Values(event.EventType, event.UserId, event.Payload, event.CreatedAt).
        Execute()
    
    if err != nil {
        return nil, err
    }

    return result, nil
}` 