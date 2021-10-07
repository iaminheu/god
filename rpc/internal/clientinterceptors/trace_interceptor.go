package clientinterceptors

import (
	"context"

	"git.zc0901.com/go/god/lib/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryTraceInterceptor rpc客户端链路追踪拦截器
func UnaryTraceInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 开启客户端跟踪操作
	ctx, span := trace.StartClientSpan(ctx, cc.Target(), method)
	defer span.Finish()

	var pairs []string
	span.Visit(func(key, value string) bool {
		pairs = append(pairs, key, value)
		return true
	})
	// **3** 将 pair 中的data以map的形式加入 ctx，传递到下一个中间件，流至下游
	ctx = metadata.AppendToOutgoingContext(ctx, pairs...)

	return invoker(ctx, method, req, reply, cc, opts...)
}
