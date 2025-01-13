package app

import (
	"github.com/justinas/alice"
	"github.com/nhdms/base-go/cmd/apis/webhook-api/app/handlers"
	"github.com/nhdms/base-go/pkg/app"
	middleware "github.com/nhdms/base-go/pkg/middlewares"
	transhttp "github.com/nhdms/base-go/pkg/transport"
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
		Name:     "webhook-api",
		producer: publisher,
	}
}

func (s *Server) SetGRPCClient(client client.Client) {
	s.client = client
}

func (s *Server) GetBasePath() string {
	return "/webhooks"
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
			Name:    "Handle Facebook Webhook",
			Method:  http.MethodPost,
			Pattern: "/facebook",
			Handler: &handlers.FacebookWebhookHandler{
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
