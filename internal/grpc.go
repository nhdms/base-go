package internal

import (
	"context"
	"github.com/micro/plugins/v5/client/grpc"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/common"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/metadata"
	"time"
)

// All grpc clients should be created here.

// ClientWrapper represents a client middleware
type ClientWrapper func(client.Client) client.Client

// NewClientWrapper creates a new client wrapper with metadata
func NewClientWrapper() ClientWrapper {
	return func(c client.Client) client.Client {
		return &customClientWrapper{c}
	}
}

type customClientWrapper struct {
	client.Client
}

func (c *customClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// Add custom metadata to context
	md := make(map[string]string)
	md["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// fetch all key value from opts[i].Context
	callOpts := client.CallOptions{
		Context: ctx,
	}
	for _, opt := range opts {
		opt(&callOpts)
	}
	ctx = metadata.MergeContext(callOpts.Context, md, false)
	return c.Client.Call(ctx, req, rsp, opts...)
}

func createGRPCClient() client.Client {
	return grpc.NewClient(
		client.Wrap(client.Wrapper(NewClientWrapper())),
		client.Registry(app.GetRegistry()),
	)
}

// CreateNewUserServiceClient creates a new UserService client.
func CreateNewUserServiceClient(conn client.Client) services.UserService {
	if conn == nil {
		conn = createGRPCClient()
	}
	return services.NewUserService(app.GetGRPCServiceName(common.ServiceNameUser), conn)
}

func CreateWebhookClient(conn client.Client) services.WebhookService {
	if conn == nil {
		conn = createGRPCClient()
	}
	return services.NewWebhookService(app.GetGRPCServiceName(common.ServiceNameWebhook), conn)
}

func CreateEventClient(conn client.Client) services.EventService {
	if conn == nil {
		conn = createGRPCClient()
	}
	return services.NewEventService(app.GetGRPCServiceName(common.ServiceNameEvent), conn)
}
