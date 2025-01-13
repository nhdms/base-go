package app

import (
	"github.com/goccy/go-json"
	"github.com/nhdms/base-go/pkg/logger"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	client2 "go-micro.dev/v5/client"
	"go-micro.dev/v5/web"
)

type API interface {
	GetRoutes() transhttp.Routes
	GetBasePath() string
	SetGRPCClient(client client2.Client)
}

func NewAPI(api API) web.Service {
	if len(GlobalServiceConfig.ServiceName) == 0 {
		logger.DefaultLogger.Fatal("missing service_name")
	}

	port := viper.GetInt("api.port")
	if port == 0 {
		port = GlobalServiceConfig.Port
	}

	options := []web.Option{
		web.Name(GetAPIName(GlobalServiceConfig.ServiceName)),
		web.Registry(GetRegistry()),
		web.Address(":" + cast.ToString(port)),
	}

	var routes transhttp.Routes
	if api != nil {
		routes = api.GetRoutes()
		options = append(options, web.Metadata(generateRoutesMetadata(api.GetBasePath(), routes)))
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

func generateRoutesMetadata(path string, routes transhttp.Routes) map[string]string {
	bytes, _ := json.Marshal(routes)
	return map[string]string{
		"endpoints": string(bytes),
		"base_path": path,
	}
}
