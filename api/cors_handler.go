package api

import "net/http"

const (
	allowOrigin  = "Access-Control-Allow-Origin"
	allOrigins   = "*"
	allowMethods = "Access-Control-Allow-Methods"
	allowHeaders = "Access-Control-Allow-Headers"
	headers      = "Content-Type, Content-Length, Origin"
	methods      = "GET, HEAD, POST, PATCH, PUT, DELETE"
)

// API 跨资源共享中间件，默认全部域名均允许
func CorsHandler(origins ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(origins) > 0 {
			w.Header().Set(allowOrigin, origins[0])
		} else {
			w.Header().Set(allowOrigin, allOrigins)
		}
		w.Header().Set(allowMethods, methods)
		w.Header().Set(allowHeaders, headers)
		w.WriteHeader(http.StatusNoContent)
	})
}
