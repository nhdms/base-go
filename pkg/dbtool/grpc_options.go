package dbtool

import (
	"context"
	"github.com/spf13/cast"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/metadata"
	metadata2 "google.golang.org/grpc/metadata"
	"strings"
)

const (
	MetadataKeySelectedFields = "selected-fields"
	MetadataKeyCacheEnable    = "cache-enable"
)

// WithFieldSelect creates a custom CallOption for field selection
func WithFieldSelect(fields ...string) client.CallOption {
	return func(o *client.CallOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = metadata.Set(o.Context, MetadataKeySelectedFields, strings.Join(fields, ","))
	}
}

func WithCacheEnable(enable bool) client.CallOption {
	return func(o *client.CallOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = metadata.Set(o.Context, MetadataKeyCacheEnable, cast.ToString(enable))
	}
}

// GetMetadataFromServer extracts metadata from incoming context
func GetMetadataFromServer(ctx context.Context, key string) (string, bool) {
	md, ok := metadata2.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}

// GetSelectedFields helper function to get selected fields
func GetSelectedFields(ctx context.Context) []string {
	value, ok := GetMetadataFromServer(ctx, MetadataKeySelectedFields)
	if !ok || value == "" {
		return nil
	}

	return strings.Split(value, ",")
}
