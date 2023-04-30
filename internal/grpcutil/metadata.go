package grpcutil

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func HeaderFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	value := md.Get(key)
	if len(value) == 0 {
		return ""
	}

	return value[0]
}
