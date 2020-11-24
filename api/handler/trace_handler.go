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
		payload, err := trace.Extract(trace.HttpFormat, r.Header)
		if err != nil && err != trace.ErrInvalidPayload {
			logx.Error(err)
		}

		ctx, span := trace.StartServerSpan(r.Context(), payload, sysx.Hostname(), r.RequestURI)
		defer span.Finish()
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
