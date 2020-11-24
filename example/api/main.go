package main

import (
	"flag"
	"git.zc0901.com/go/god/api"
	"git.zc0901.com/go/god/api/httpx"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/service"
	"net/http"
)

func main() {
	flag.Parse()

	// 新建 API 服务器
	server := api.MustNewServer(api.ApiConf{
		ServiceConf: service.ServiceConf{
			LogConf: logx.LogConf{Mode: "console"},
		},
		Port:     8080,
		Timeout:  3000,
		MaxConns: 500,
	}, api.WithNotAllowedHandler(api.CorsHandler()))
	defer server.Stop()

	// 增加路由
	server.AddRoute(api.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJson(w, "hello, world!")
		},
	})

	// 启动 API 服务器
	server.Start()
}
