package handler

import (
	"git.zc0901.com/go/god/lib/stat"
	"git.zc0901.com/go/god/lib/timex"
	"net/http"
)

// API 耗时监控中间件
func MetricHandler(metrics *stat.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			defer func() {
				metrics.Add(stat.Task{
					Duration: timex.Since(startTime),
				})
			}()

			next.ServeHTTP(w, r)
		})
	}
}
