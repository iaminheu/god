package handler

import (
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/sysx"
	"git.zc0901.com/go/god/lib/trace"
	"net/http"
)

// API 链路跟踪中间件
func TraceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// **1** 拦截请求，获取 header 中的traceId信息
		payload, err := trace.Extract(trace.HttpFormat, r.Header)
		if err != nil && err != trace.ErrInvalidPayload {
			logx.Error(err)
		}

		// **2** 在此中间件中开启一个新的 根span（即操作），并把「traceId，spanId」封装在context中
		ctx, span := trace.StartServerSpan(r.Context(), payload, sysx.Hostname(), r.RequestURI)
		defer span.Finish()

		// **5** 如果没有设置，则随机生成返回
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
