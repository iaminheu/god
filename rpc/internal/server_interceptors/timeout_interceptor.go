package server_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/contextx"
	"google.golang.org/grpc"
	"time"
)

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := contextx.ShrinkDeadline(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}
