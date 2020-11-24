package handler

import (
	"git.zc0901.com/go/god/api/internal"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/syncx"
	"net/http"
)

// API 最大请求连接数中间件
func MaxConns(max int) func(http.Handler) http.Handler {
	if max <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		limit := syncx.NewLimit(max)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limit.TryBorrow() {
				defer func() {
					if err := limit.Return(); err != nil {
						logx.Error(err)
					}
				}()

				next.ServeHTTP(w, r)
			} else {
				internal.Errorf(r, "并发连接数超出 %d，错误码：%d",
					max, http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		})
	}
}
