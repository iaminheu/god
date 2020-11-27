package svc

import "git.zc0901.com/go/god/example/graceful/dns/rpc/internal/config"

type ServiceContext struct {
	c config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		c: c,
	}
}
