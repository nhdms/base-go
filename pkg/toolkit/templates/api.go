package templates

const APICICDTemplate = `{
  "app_type": "api",
  "cmd_bin_dir": "cmd/apis/{{.Name}}",
  "service_name": "{{.ServiceName}}",
  "port": {{.Port}},
  "config_remote_keys": [
    "database/rabbitmq.toml",
    "apis/{{.ServiceName}}/config.toml"
  ]
}`

const APIMainTemplate = `package main

import (
    app2 "github.com/nhdms/base-go/cmd/apis/{{.Name}}/app"
    "github.com/nhdms/base-go/pkg/app"
    "github.com/nhdms/base-go/pkg/logger"
)

func main() {
    publisher, err := app.NewPublisher()
    if err != nil {
        logger.DefaultLogger.Fatal("Start publisher failed", err)
    }
    defer publisher.Close()

    s := app2.NewServer(publisher)
    api := app.NewAPI(s)
    err = api.Run()
    if err != nil {
        logger.DefaultLogger.Fatal("Start API failed", err)
    }
}`

const APIWebServerTemplate = `package app

import (
    "github.com/justinas/alice"
    "github.com/spf13/viper"
    "github.com/nhdms/base-go/cmd/apis/{{.Name}}/handlers"
    "github.com/nhdms/base-go/pkg/app"
    middleware "github.com/nhdms/base-go/pkg/middlewares"
    transhttp "github.com/nhdms/base-go/pkg/transport"
    "go-micro.dev/v5/client"
    "net/http"
)

type Server struct {
    Name     string
    client   client.Client
    producer app.PublisherInterface
}

func NewServer(publisher app.PublisherInterface) *Server {
    return &Server{
        Name:     "{{.ServiceName}}",
        producer: publisher,
    }
}

func (s *Server) SetGRPCClient(client client.Client) {
    s.client = client
}

func (s *Server) GetBasePath() string {
    return "/{{.ServiceName}}"
}

func (s *Server) GetName() string {
    return s.Name
}

func (s *Server) GetRoutes() transhttp.Routes {
    mdws := []alice.Constructor{}

    if viper.GetBool("logging.enable") {
        mdws = append(mdws, middleware.LoggingMiddleware)
    }

    return []transhttp.Route{
        {
            Name:    "Create {{.Handler}} Event",
            Method:  http.MethodPost,
            Pattern: "/events",
            Handler: &handlers.{{.Handler}}Handler{
                Producer: s.producer,
            },
            Middlewares: mdws,
            AuthInfo: transhttp.AuthInfo{
                Enable: false,
            },
            Timeout: 10000, // 10 seconds
        },
    }
}`

const APIHandlerTemplate = `package handlers

import (
    "github.com/nhdms/base-go/pkg/app"
    "github.com/nhdms/base-go/pkg/logger"
    transhttp "github.com/nhdms/base-go/pkg/transport"
    "io"
    "net/http"
    "time"
)

type {{.Handler}}Handler struct {
    Producer app.PublisherInterface
}

func (h *{{.Handler}}Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    defer func() {
        logger.DefaultLogger.Debugw("Processed request", "url", r.URL.Path, "took", time.Since(start).Milliseconds())
    }()

    requestBody, err := io.ReadAll(r.Body)
    if err != nil {
        transhttp.RespondError(w, http.StatusBadRequest, err.Error())
        return
    }

    logger.DefaultLogger.Debugw("published request", "url", r.URL.Path,"body", string(requestBody), "took", time.Since(start).Milliseconds())

    transhttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
        "success": time.Now().UnixNano(),
    })
}`
