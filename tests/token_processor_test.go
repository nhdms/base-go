package tests

import (
	"context"
	grpcClient "github.com/micro/plugins/v5/client/grpc"
	"github.com/micro/plugins/v5/registry/consul"
	"github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/internal/permissions"
	"github.com/nhdms/base-go/internal/token"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"go-micro.dev/v5/client"
	"log"
	"testing"
)

var processor token.Processor
var usersvc services.UserService
var ctx = context.Background()

var tokenString = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiQ0lRREVpai9mLy8vL3dFUGZ3OEJIeDhmSC8vLy93SC8vdzhmQXdkL0J6OC9BejhIQndjSEJ3TVBCd2NER2dVSTN3SVFWQm9GQ044Q0VGc2FCUWkwQVJCVUdnVUk2QUlRUEJvRUNDRVFQaG9FQ0NBUVBob0VDQjhRUGhvRUNCNFFQaG9FQ0J3UVBob0VDQnNRUGhvRUNCb1FQaG9FQ0FNUVBob0VDQ0FRUWhvRUNDQVFWQm9FQ0NBUVd4b0ZDT01DRUQ4YUJRamZBaEEvR2dVSXhRRVFQeG9GQ0s4QkVEOGFCUWkwQVJBL0dnVUltd0VRUHhvRUNIOFFQeG9FQ0NFUVB4b0VDQ0FRUHhvRUNCOFFQeG9FQ0I0UVB4b0VDQndRUHhvRUNCc1FQeG9FQ0JvUVB4b0VDQU1RUHhvRUNIOFFQaG9GQ0xRQkVENGFCUWpmQWhBK0dnVUlyd0VRV3hvRkNLOEJFRlFhQlFpdkFSQkNHZ1VJcndFUVBob0ZDT01DRUVJYUJBaC9FRUlhQkFnYkVFSWFCQWdERUVJYUJnamZBaERZQmhvRUNCc1FXeG9FQ0FNUVd4b0VDQnNRVkJvRUNDRVFWQm9GQ09jQ0VEd2FCQWdERUZRYUJRaWJBUkJDR2dVSW13RVFWQ0pZbGdqcEI5QUlBUVhYQjlzSDJBajBBcU1Ea2dUdUFxc0lvaEtsRXZ3RTd3THdBdklDOHdMMUF2WUM5d0w1QXZvQzZBZldCOU1DM1FLTUJaOEY0UUxnQXI0Q3dBSzhBcjhDeFFMRUFzSUN4Z0wwQWNnQ3lnTExBaW9nVlc1T2VVWm9XWGxOUkdWelQxWmhZa1J5YTBWRE9FSlZjV2xhVTJjeVdFaz0iLCJNYXBDbGFpbXMiOnsiZXhwIjoxNzMyNTEwNzQzfX0.dCbxgLjuWhQPN8Q60oKwVm5me9ec2TH1RzO9raWFNQ8`

//var tokenString = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiQ05FU0VpSUEvLzkvRDM4QUFBQUFBQUQvUHdBQUFBQi9Cd01IQUFBQUFBQUFBQUFBQUFBQUdnVUl4UUVRVkJvRkNOOENFRlFhQlFqakFoQlVHZ1VJNHdJUVZCb0ZDTjhDRUZRYUJRakZBUkJVR2dVSXRBRVFWQm9GQ0s4QkVGUWFCUWliQVJCVUdnUUlmeEJVR2dRSUlSQlVHZ1FJSUJCVUdnUUlIeEJVR2dRSUhoQlVHZ1FJSEJCVUdnUUlHeEJVR2dRSUdoQlVHZ1FJQXhCVUtpQlFUVEp0V0RCVGRrWldNMjFJZVdsRk0zSlpaRlpZVDFSRE5XczBiVGN6U2c9PSIsIk1hcENsYWltcyI6eyJleHAiOjE3MzI0Mzg3OTB9fQ._cAIdblLry5qKOPBSyo9rq2QEDw4s-Xux58Jb9nRSwI`

func init() {
	err := LoadTestConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	rd, err := dbtool.CreateRedisConnection(nil)
	if err != nil {
		log.Fatal("Failed to create Redis connection", err)
	}
	c := grpcClient.NewClient(
		client.Registry(consul.NewRegistry()),
		//client.Selector(se),
	)

	usersvc = internal.CreateNewUserServiceClient(c)
	processor = token.NewTokenProcessor(rd, usersvc)
}

func TestGenerateToken(t *testing.T) {
	user, err := usersvc.GetUserByID(ctx, &services.UserRequest{
		UserId: 388,
	})
	if err != nil {
		t.Fatal("Failed to get user", err)
	}

	token, err := processor.GenerateToken(ctx, user.User)
	if err != nil {
		t.Fatal("Failed to generate token", err)
	}
	t.Log("Generated token: ", token)
}

func TestCheckTokenPermission(t *testing.T) {
	tk, err := processor.GetToken(ctx, tokenString)
	if err != nil {
		t.Fatal("Failed to generate token", err)
	}

	claim := token.GetJWTClaimFromToken(tk)
	t.Log("user", claim.UserId)
	t.Log("permission", claim.Permissions)
}

func TestCheckPermissions(t *testing.T) {
	requirePermissions := map[int64]int64{
		int64(permissions.Dashboard): permissions.DashboardCarePageOverview | permissions.DashboardMarketingTelesales,
	}

	token, err := processor.GetToken(ctx, tokenString)
	if err != nil {
		t.Fatal("Failed to get token", err)
	}

	if !token.Valid {
		t.Fatal("invalid token")
	}

	ok := processor.CheckPermissions(token, requirePermissions)
	if !ok {
		t.Fatal("user has no permission")
	}
	t.Log("access granted")
}
