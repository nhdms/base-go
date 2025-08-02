package app

import (
	"fmt"
	"github.com/nhdms/base-go/pkg/micro/grpc"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/spf13/cast"
	"go-micro.dev/v5"
	"go-micro.dev/v5/server"
	"os"
)

const (
	DefaultGRPCPort = 6688
)

func NewGRPCService(serviceName string) micro.Service {
	if len(serviceName) == 0 {
		serviceName = GlobalServiceConfig.ServiceName
	}

	name := GetGRPCServiceName(serviceName)
	port := cast.ToInt(os.Getenv("PORT"))
	if port < 1 {
		port = GlobalServiceConfig.Port
	}

	if !utils.IsLocalDevelopment() {
		port = DefaultGRPCPort
	}

	grpcServer := grpc.NewServer(
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
