// Code generated by god. DO NOT EDIT!
// Source: test.proto

package main

import (
	"flag"
	"fmt"

	"git.zc0901.com/go/god/example/test/cmd/rpc/internal/config"
	"git.zc0901.com/go/god/example/test/cmd/rpc/internal/server"
	"git.zc0901.com/go/god/example/test/cmd/rpc/internal/svc"
	"git.zc0901.com/go/god/example/test/cmd/rpc/test"
	"git.zc0901.com/go/god/lib/conf"
	"git.zc0901.com/go/god/rpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/test.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewTestServer(ctx)

	s := rpc.MustNewServer(c.ServerConf, func(grpcServer *grpc.Server) {
		test.RegisterTestServer(grpcServer, srv)
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
