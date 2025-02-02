package handler

import (
	"git.zc0901.com/go/god/api/internal/security"
	"git.zc0901.com/go/god/lib/prometheus"
	"git.zc0901.com/go/god/lib/prometheus/metric"
	"git.zc0901.com/go/god/lib/timex"
	"net/http"
	"strconv"
	"time"
)

const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "API服务器请求耗时（毫秒）。",
		Labels:    []string{"path"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "API服务器请求响应码计数器。",
		Labels:    []string{"path", "code"},
	})
)

// PrometheusHandler API 监控中间件
func PrometheusHandler(path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !prometheus.Enabled() {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			cw := &security.WithCodeResponseWriter{Writer: w}
			defer func() {
				metricServerReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), path)
				metricServerReqCodeTotal.Inc(path, strconv.Itoa(cw.Code))
			}()

			next.ServeHTTP(cw, r)
		})
	}
}
