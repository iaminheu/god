package main

import (
	"flag"
	"git.zc0901.com/go/god/api"
	"git.zc0901.com/go/god/api/httpx"
	"git.zc0901.com/go/god/lib/conf"
	"net/http"
)

var configFile = flag.String("f", "config.yaml", "API 配置文件")

func main() {
	// 读取配置文件
	flag.Parse()
	var c api.ApiConf
	conf.MustLoad(*configFile, &c)

	// 新建 API 服务器
	server := api.MustNewServer(c, api.WithNotAllowedHandler(api.CorsHandler()))
	defer server.Stop()

	// 增加路由
	server.AddRoute(api.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJson(w, "hello, world!")
		},
	})
	server.AddRoute(api.Route{
		Method: http.MethodGet,
		Path:   "/ping",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJson(w, "pong")
		},
	})

	// 启动 API 服务器
	server.Start()
}
