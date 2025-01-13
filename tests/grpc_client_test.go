package tests

import (
	"context"
	transhttp "github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"go-micro.dev/v5/client"
	"log"
	"testing"
	"time"
)

func TestRequestGRPCServer(t *testing.T) {
	st := time.Now()
	cl := transhttp.CreateNewUserServiceClient(nil)
	d, e := cl.GetUserByID(context.Background(), &services.UserRequest{
		UserId: 388,
	}, client.WithRequestTimeout(51*time.Second))

	if e != nil {
		log.Fatal(e)
	}
	log.Println(time.Since(st), "Get user by ID: ", d)
}
