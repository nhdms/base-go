package templates

const ConsumerMainTemplate = `package main

import (
    "github.com/nhdms/base-go/cmd/consumers/{{.Name}}/handlers"
    "github.com/nhdms/base-go/pkg/app"
    "github.com/nhdms/base-go/pkg/logger"
)

func main() {
    {{.ServiceName}}Handler := &handlers.{{.Handler}}Handler{
        Name: "{{.Name}}",
    }

    err := app.StartNewConsumer({{.ServiceName}}Handler)
    if err != nil {
        logger.DefaultLogger.Fatal("Failed to start consumer: ", err)
    }
}`

const ConsumerHandlerTemplate = `package handlers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nhdms/base-go/internal"
    "github.com/nhdms/base-go/pkg/app"
    "github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/nhdms/base-go/pkg/logger"
)

type {{.Handler}}Handler struct {
    Publisher                   app.PublisherInterface
    {{.Handler}}Client         services.{{.Handler}}Service
    Name                       string
}

func (h *{{.Handler}}Handler) GetName() string {
    return h.Name
}

func (h *{{.Handler}}Handler) Init() error {
	h.{{.Handler}}Client  = internal.Create{{.Handler}}Client(nil)
    return nil
}

func (h *{{.Handler}}Handler) HandleMessage(msg *message.Message) error {
	logger.DefaultLogger.Debugw("Received ", "message", string(msg.Payload))
    return nil
}

func (h *{{.Handler}}Handler) SetPublisher(p app.PublisherInterface) {
	h.Publisher = p
}

func (h *{{.Handler}}Handler) Close() {}
`

const ConsumerCICDTemplate = `{
    "app_type": "consumer",
    "cmd_bin_dir": "cmd/consumers/{{.Name}}",
    "service_name": "{{.ServiceName}}",
    "port": 0,
    "config_remote_keys": [
        "database/redis.toml",
        "database/rabbitmq.toml",
        "consumers/{{.Name}}.toml"
    ]
}`

const ConsumerSampleConfig = `[consumers]
[consumers.{{.Name}}]
exchange = "{{.Name}}"
queue = "{{.Name}}"
routing_key = ""
type = "direct"
auto_delete = false
durable = true
exclusive = false
#disable=true
qos = 50
worker_count = 1
#additional_bindings =["ex1:routing_key_1", "ex2:routing_key_2", "ex3"]

[logger]
level="debug"
`
