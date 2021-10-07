package svc

import (
	"git.zc0901.com/go/god/example/test/cmd/api/internal/config"
	"git.zc0901.com/go/god/example/test/cmd/rpc/testclient"
	"git.zc0901.com/go/god/rpc"
)

type ServiceContext struct {
	Config  config.Config
	TestRPC testclient.Test
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		TestRPC: testclient.NewTest(rpc.MustNewClient(c.TestRPC)),
	}
}
