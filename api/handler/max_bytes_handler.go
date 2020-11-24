package handler

import (
	"git.zc0901.com/go/god/api/internal"
	"net/http"
)

// API 最大请求字节中间件
func MaxBytesHandler(max int64) func(http.Handler) http.Handler {
	if max <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > max {
				internal.Errorf(r, "请求实体过大，限制为：%d，接收到：%d，错误码：%d",
					max, r.ContentLength, http.StatusRequestEntityTooLarge)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
