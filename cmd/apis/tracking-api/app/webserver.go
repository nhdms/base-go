package app

import (
    "github.com/justinas/alice"
    "github.com/spf13/viper"
    "gitlab.com/a7923/athena-go/cmd/apis/tracking-api/handlers"
    "gitlab.com/a7923/athena-go/pkg/app"
    middleware "gitlab.com/a7923/athena-go/pkg/middlewares"
    transhttp "gitlab.com/a7923/athena-go/pkg/transport"
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
        Name:     "tracking",
        producer: publisher,
    }
}

func (s *Server) SetGRPCClient(client client.Client) {
    s.client = client
}

func (s *Server) GetBasePath() string {
    return "/tracking"
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
            Name:    "Create Tracking Event",
            Method:  http.MethodPost,
            Pattern: "/events",
            Handler: &handlers.TrackingHandler{
                Producer: s.producer,
            },
            Middlewares: mdws,
            AuthInfo: transhttp.AuthInfo{
                Enable: false,
            },
            Timeout: 10000, // 10 seconds
        },
    }
}