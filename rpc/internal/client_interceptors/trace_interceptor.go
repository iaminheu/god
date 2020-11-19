package client_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TraceInterceptor rpc客户端链路追踪拦截器
func TraceInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx, span := trace.StartClientSpan(ctx, cc.Target(), method)
	defer span.Finish()

	var pairs []string
	span.Visit(func(key, value string) bool {
		pairs = append(pairs, key, value)
		return true
	})
	ctx = metadata.AppendToOutgoingContext(ctx, pairs...)

	return invoker(ctx, method, req, reply, cc, opts...)
}
