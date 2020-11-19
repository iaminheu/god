package server_interceptors

import (
	"context"
	"git.zc0901.com/go/god/lib/prometheus/metric"
	"git.zc0901.com/go/god/lib/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const serverNamespace = "rpc_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc 服务器请求时长(毫秒)",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc 服务器请求代码计数器",
		Labels:    []string{"method", "code"},
	})
)

// 统计rpc服务端请求时长和状态代码
func UnaryPrometheusInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := timex.Now()
		resp, err := handler(ctx, req)
		metricServerReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), info.FullMethod)
		metricServerReqCodeTotal.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
		return resp, err
	}
}
