package handler

import (
	"git.zc0901.com/go/god/api/httpx"
	"git.zc0901.com/go/god/api/internal/security"
	"git.zc0901.com/go/god/lib/load"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stat"
	"net/http"
	"sync"
)

const serviceType = "api"

var (
	shedderStat *load.ShedderStat
	lock        sync.Mutex
)

// API 负载均衡中间件
func ShedderHandler(shedder load.Shedder, metrics *stat.Metrics) func(http.Handler) http.Handler {
	if shedder == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	ensureShedderStat()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shedderStat.IncrTotal()
			promise, err := shedder.Allow()
			if err != nil {
				metrics.AddDrop()
				shedderStat.IncrDrop()
				logx.Errorf("[http] 请求被丢弃，%s - %s - %s",
					r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			cw := &security.WithCodeResponseWriter{Writer: w}
			defer func() {
				if cw.Code == http.StatusServiceUnavailable {
					promise.Fail()
				} else {
					shedderStat.IncrPass()
					promise.Pass()
				}
			}()

			next.ServeHTTP(cw, r)
		})
	}
}

func ensureShedderStat() {
	lock.Lock()
	if shedderStat == nil {
		shedderStat = load.NewShedderStat(serviceType)
	}
	lock.Unlock()
}
