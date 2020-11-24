package rpc

import (
	"git.zc0901.com/go/god/lib/discovery"
	"git.zc0901.com/go/god/lib/service"
	"git.zc0901.com/go/god/lib/store/redis"
)

type (
	// ServerConfig Rpc服务端配置
	ServerConfig struct {
		service.ServiceConf
		ListenOn      string
		Etcd          discovery.EtcdConf `json:",optional"`
		Auth          bool               `json:",optional"` // rpc通信鉴权
		Redis         redis.KeyConf      `json:",optional"` // 鉴权redis存储及鉴权key
		StrictControl bool               `json:",optional"`

		// 不允许永久挂起，不要设置超时时间为0，如设为0，底层将自动重置为2秒。
		Timeout      int64 `json:",default=2000"`
		CpuThreshold int64 `json:",default=900,range=[0:1000]"`
	}

	// ServerConfig Rpc客户端配置
	ClientConf struct {
		Etcd      discovery.EtcdConf `json:",optional"`
		Endpoints []string           `json:",optional=!Etcd"`
		App       string             `json:",optional"`
		Token     string             `json:",optional"`
		Timeout   int64              `json:",optional"`
	}
)

// NewDirectClientConf 新建Rpc直联客户端配置
func NewDirectClientConf(endpoints []string, app, token string) ClientConf {
	return ClientConf{
		Endpoints: endpoints,
		App:       app,
		Token:     token,
	}
}

// NewEtcdClientConf 新建Rpc Etcd连接客户端配置
func NewEtcdClientConf(hosts []string, key, app, token string) ClientConf {
	return ClientConf{
		Etcd: discovery.EtcdConf{
			Hosts: hosts,
			Key:   key,
		},
		App:   app,
		Token: token,
	}
}

// HasEtcd 判断服务端配置是否传递Etcd主机和键
func (sc ServerConfig) HasEtcd() bool {
	return len(sc.Etcd.Hosts) > 0 && len(sc.Etcd.Key) > 0
}

// Validate 验证鉴权情况下redis的主机和部署模式是否已提供
func (sc ServerConfig) Validate() error {
	if sc.Auth {
		if err := sc.Redis.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// HasCredential 判断客户端配置是否存在App+token凭证
func (cc ClientConf) HasCredential() bool {
	return len(cc.App) > 0 && len(cc.Token) > 0
}
