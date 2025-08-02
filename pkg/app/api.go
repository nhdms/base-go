package app

import (
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/micro/web"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	client2 "go-micro.dev/v5/client"
)

const (
	DefaultAPIPort = 8080
)

type API interface {
	GetRoutes() web.Routes
	GetBasePath() string
	SetGRPCClient(client client2.Client)
}

func NewAPI(api API) web.Service {
	if len(GlobalServiceConfig.ServiceName) == 0 {
		logger.DefaultLogger.Fatal("missing service_name")
	}

	port := viper.GetInt("api.port")

	// we will deploy to k8s, to implement health check for all services we need to expose a fixed port 8080
	if !utils.IsLocalDevelopment() {
		port = DefaultAPIPort
	}

	if port == 0 {
		port = GlobalServiceConfig.Port
	}

	var routes web.Routes
	var basePath string
	if api != nil {
		routes = api.GetRoutes()
		basePath = api.GetBasePath()
	}

	serviceName := GlobalServiceConfig.ServiceName
	if len(basePath) > 0 {
		serviceName = basePath
	}

	options := []web.Option{
		web.Name(GetAPIName(serviceName)),
		web.Registry(GetRegistry()),
		web.Address(":" + cast.ToString(port)),
		web.WithRoutes(routes),
		web.BasePath(basePath),
	}

	svc := web.NewService(
		options...,
	)

	if routes != nil {
		transhttp.InitRoutes(svc, routes, api.GetBasePath())
	}

	err := svc.Init()
	if err != nil {
		logger.DefaultLogger.Fatal("Init API failed", err)
	}

	return svc
}
