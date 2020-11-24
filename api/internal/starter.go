package internal

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/proc"
	"net/http"
)

func StartHttp(host string, port int, handler http.Handler) error {
	return start(host, port, handler, func(server *http.Server) error {
		return server.ListenAndServe()
	})
}

func StartHttps(host string, port int, certFile, keyFile string, handler http.Handler) error {
	return start(host, port, handler, func(server *http.Server) error {
		// 证书文件和秘钥文件在 buildHttpsServer 中设置
		return server.ListenAndServeTLS(certFile, keyFile)
	})
}

func start(host string, port int, handler http.Handler, run func(server *http.Server) error) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}
	waitForCalled := proc.AddWrapUpListener(func() {
		server.Shutdown(context.Background())
	})
	defer waitForCalled()

	return run(server)
}
