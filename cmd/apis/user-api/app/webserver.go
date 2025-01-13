package app

import (
	"github.com/justinas/alice"
	handlers3 "github.com/nhdms/base-go/cmd/apis/user-api/app/handlers"
	"github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/pkg/logger"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"go-micro.dev/v5/client"
	"net/http"
)

type Server struct {
	Name   string
	client client.Client
}

func (s *Server) SetGRPCClient(client client.Client) {
	s.client = client
}

func (s *Server) GetBasePath() string {
	return "/users"
}

func (s *Server) GetName() string {
	return s.Name
}

func (s *Server) GetRoutes() transhttp.Routes {
	cl := internal.CreateNewUserServiceClient(nil)
	return []transhttp.Route{
		{
			Name:    "Hello world",
			Method:  http.MethodPut,
			Pattern: "/test/{id}/report",
			Handler: &handlers3.GetUserByIdHandler{
				UserClient: cl,
			},
			Middlewares: []alice.Constructor{
				loggingMiddleware,
			},
			AuthInfo: transhttp.AuthInfo{
				Enable: false,
			},
			Timeout: 10000, // 10 seconds
		},
		{
			Name:    "Hello world",
			Method:  http.MethodGet,
			Pattern: "/",
			Handler: &handlers3.GetUserByIdHandler{
				UserClient: cl,
			},
			Middlewares: []alice.Constructor{
				loggingMiddleware,
			},
			AuthInfo: transhttp.AuthInfo{},
			Timeout:  10000, // 10 seconds
		},
		{
			Name:    "Hello world put",
			Method:  http.MethodPost,
			Pattern: "/put",
			Handler: &handlers3.GetUserByIdHandler{
				UserClient: cl,
			},
			Middlewares: []alice.Constructor{
				loggingMiddleware,
			},
			AuthInfo: transhttp.AuthInfo{},
			Timeout:  10000, // 10 seconds
		},
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.DefaultLogger.Debugf("Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
