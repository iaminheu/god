package server_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryTraceInterceptor 一元链路追踪拦截器
func UnaryTraceInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		// 提交上下文metadata中的负载参数信息
		payload, err := trace.Extract(trace.GrpcFormat, md)
		if err != nil {
			return handler(ctx, req)
		}

		ctx, span := trace.StartServerSpan(ctx, payload, serviceName, info.FullMethod)
		defer span.Finish()
		return handler(ctx, req)
	}
}
