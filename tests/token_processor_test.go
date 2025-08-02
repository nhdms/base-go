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

var tokenString = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEyMywiZnVsbG5hbWUiOiJExaluZyBIVCAgLSBNS1QgLSBMZWFkZXIgSEFERVMiLCJjb21wYW55SWQiOjIsInR5cGUiOiJlbXBsb3llZSIsIndhcmVob3VzZXMiOm51bGwsImRpc3BsYXlJZCI6IkVPMTI0Iiwic2lkIjoidUtMejhlWWN5bEc5UlI1dVpSS1dTc1JNQmdrZkRQMFEiLCJzZXJ2aWNlIjoic2FsZSIsInAiOlsiMCIsIjY1NDM5IiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIyMDkxMDA3IiwiMCIsIjciLCIxIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCIsIjAiLCIwIiwiMCJdLCJkaWRzIjpbMTIzLDkyMyw0ODc5LDQ5NzUsNTAzMF0sImlhdCI6MTc0NzM2MzY4OSwiZXhwIjoxNzU1MTM5Njg5fQ.LiOxOP2v4j-FfjzsbbppsZWOeb88NdHoknkbixgvXCE`

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
	processor = token.NewAGSaleTokenProcessor(rd, usersvc)
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
	t.Log("user", claim.Sub)
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
