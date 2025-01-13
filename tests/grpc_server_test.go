package tests

import (
	"context"
	grpcServer "github.com/nhdms/base-go/grpc"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"go-micro.dev/v5"
	"log"
	"testing"
)

type sample struct {
}

func (s sample) GetProjects(ctx2 context.Context, request *services.ProjectRequest, response *services.ProjectResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s sample) GetProfiles(ctx context.Context, request *services.ProfileRequest, response *services.ProfileResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s sample) GetDepartments(ctx context.Context, request *services.DepartmentRequest, response *services.DepartmentResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s sample) GetRoles(ctx context.Context, request *services.RoleRequest, response *services.RoleResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s sample) GetDataSetScopes(ctx context.Context, request *services.DataSetScopeRequest, response *services.DataSetScopeResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s sample) GetUserByID(ctx context.Context, request *services.UserRequest, response *services.UserResponse) error {
	//TODO implement me
	panic("implement me")
}

func TestCreateServer(t *testing.T) {
	service := micro.NewService(
		micro.Name("your.service.name"),
		micro.Server(grpcServer.NewServer()),
	)
	service.Init()

	// Register your service handlers
	if err := services.RegisterUserServiceHandler(service.Server(), new(sample)); err != nil {
		log.Fatal(err)
	}

	// Type assert the Micro server to the gRPC server
	x := service.Server()
	_ = x
	// Run the micro service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
