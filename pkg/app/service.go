package app

import (
	"fmt"
	svr "github.com/nhdms/base-go/grpc"
	"github.com/spf13/cast"
	"go-micro.dev/v5"
	"go-micro.dev/v5/server"
	"os"
)

func NewGRPCService() micro.Service {
	name := GetGRPCServiceName(GlobalServiceConfig.ServiceName)
	port := cast.ToInt(os.Getenv("PORT"))
	if port < 1 {
		port = GlobalServiceConfig.Port
	}
	grpcServer := svr.NewServer(
		server.Name(name),
		server.Address(fmt.Sprintf(":%d", port)),
		server.Registry(GetRegistry()),
	)

	// Create new service
	svc := micro.NewService(
		micro.Server(grpcServer),
		micro.Name(name),
	)

	svc.Init()
	return svc
}
