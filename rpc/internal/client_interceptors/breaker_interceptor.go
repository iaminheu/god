package client_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/breaker"
	"git.zc0901.com/go/god/rpc/internal/codes"
	"google.golang.org/grpc"
	"path"
)

func BreakerInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptable(breakerName, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, codes.Acceptable)
}
