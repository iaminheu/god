package handler

import (
	"net/http"
	"time"
)

const reason = "API 请求超时"

// TimeoutHandler API 超时处理中间件
func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			return http.TimeoutHandler(next, duration, reason)
		}
		return next
	}
}
