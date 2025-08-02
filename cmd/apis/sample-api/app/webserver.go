package app

import (
	"github.com/justinas/alice"
	"github.com/nhdms/base-go/cmd/apis/sample-api/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/micro/web"
	middleware "github.com/nhdms/base-go/pkg/middlewares"
	"github.com/spf13/viper"
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
		Name:     "sample",
		producer: publisher,
	}
}

func (s *Server) SetGRPCClient(client client.Client) {
	s.client = client
}

func (s *Server) GetBasePath() string {
	return "/samples"
}

func (s *Server) GetName() string {
	return s.Name
}

func (s *Server) GetRoutes() web.Routes {
	mdws := []alice.Constructor{}

	if viper.GetBool("logging.enable") {
		mdws = append(mdws, middleware.LoggingMiddleware)
	}

	return []web.Route{
		{
			Name:    "Create Sample Event",
			Method:  http.MethodPost,
			Pattern: "/events",
			Handler: &handlers.SampleHandler{
				Producer: s.producer,
			},
			Middlewares: mdws,
			AuthInfo: web.AuthInfo{
				Enable: false,
			},
			Timeout: 10000, // 10 seconds
		},
	}
}
