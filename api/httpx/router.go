package httpx

import "net/http"

// 继承于 http.Handler 的路由器
type Router interface {
	http.Handler
	Handle(method, path string, handler http.Handler) error
	SetNotFoundHandler(handler http.Handler)
	SetNotAllowedHandler(handler http.Handler)
}
