package security

import (
	"bufio"
	"net"
	"net/http"
)

// 带有响应码的输出器
type WithCodeResponseWriter struct {
	Writer http.ResponseWriter
	Code   int
}

func (w *WithCodeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Writer.(http.Hijacker).Hijack()
}
func (w *WithCodeResponseWriter) Header() http.Header {
	return w.Writer.Header()
}

func (w *WithCodeResponseWriter) Write(bytes []byte) (int, error) {
	return w.Writer.Write(bytes)
}

func (w *WithCodeResponseWriter) WriteHeader(code int) {
	w.Writer.WriteHeader(code)
	w.Code = code
}
