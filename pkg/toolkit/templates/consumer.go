package templates

const ConsumerMainTemplate = `package main

import (
    "gitlab.com/a7923/athena-go/cmd/consumers/{{.Name}}/handlers"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/logger"
    "gitlab.com/a7923/athena-go/proto/exmsg/services"
    "go-micro.dev/v5"
)

func main() {
    svc := micro.NewService()
    svc.Init()

    // Create gRPC client
    client := services.New{{.Handler}}Service("{{.ServiceName}}", svc.Client())
    
    handler := handlers.New{{.Handler}}Handler(client)
    err := app.StartDataReplica(handler)
    if err != nil {
        logger.AthenaLogger.Fatal("Failed to start consumer", err)
    }
}`

const ConsumerHandlerTemplate = `package handlers

import (
    "context"
    "encoding/json"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/proto/exmsg/models"
    "gitlab.com/a7923/athena-go/proto/exmsg/services"
)

type {{.Handler}}Handler struct {
    client services.{{.Handler}}Service
}

func New{{.Handler}}Handler(client services.{{.Handler}}Service) *{{.Handler}}Handler {
    return &{{.Handler}}Handler{
        client: client,
    }
}

func (h *{{.Handler}}Handler) Process(msg []byte) error {
    var event models.{{.Handler}}Event
    err := json.Unmarshal(msg, &event)
    if err != nil {
        return err
    }

    // Call gRPC service to save data
    _, err = h.client.Create{{.Handler}}(context.Background(), &event)
    return err
}

func (h *{{.Handler}}Handler) GetName() string {
    return "{{.Name}}"
}

func (h *{{.Handler}}Handler) Init() error {
    return nil
}

func (h *{{.Handler}}Handler) Close() error {
    return nil
}`
