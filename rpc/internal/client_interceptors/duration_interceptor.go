package client_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/timex"
	"google.golang.org/grpc"
	"path"
	"time"
)

const slowThreshold = time.Millisecond * 500

// DurationInterceptor rpc调用时长拦截器
func DurationInterceptor(ctx context.Context, method string, req, replay interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serviceName := path.Join(cc.Target(), method)
	startTime := timex.Now()
	err := invoker(ctx, method, req, replay, cc, opts...)
	if err != nil {
		logx.WithContext(ctx).WithDuration(timex.Since(startTime)).Infof("失败 - %s - %v - %s",
			serviceName, req, err.Error())
	} else {
		elapsed := timex.Since(startTime)
		if elapsed > slowThreshold {
			logx.WithContext(ctx).WithDuration(elapsed).Slowf("[RPC] ok - 慢查询 - %s -%v - %v",
				serviceName, req, replay)
		}
	}

	return err
}
