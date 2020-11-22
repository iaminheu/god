package rpc

import (
	"git.zc0901.com/go/god/lib/discovery"
	"git.zc0901.com/go/god/rpc/internal"
	"git.zc0901.com/go/god/rpc/internal/auth"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	WithDialOption             = internal.WithDialOption
	WithTimeout                = internal.WithTimeout
	WithUnaryClientInterceptor = internal.WithUnaryClientInterceptor
)

type (
	ClientOption = internal.ClientOption

	Client interface {
		Conn() *grpc.ClientConn
	}

	RpcClient struct {
		client Client
	}
)

func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}

// MustNewClient 根据配置文件新建rpc客户端，出错直接退出。
func MustNewClient(conf ClientConf, options ...ClientOption) Client {
	cli, err := NewClient(conf, options...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

// MustNewClient 根据配置文件新建rpc客户端，出错直接退出。
func NewClient(c ClientConf, options ...ClientOption) (Client, error) {
	var opts []ClientOption
	if c.HasCredential() {
		opts = append(opts, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   c.App,
			Token: c.Token,
		})))
	}
	if c.Timeout > 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	opts = append(opts, options...)

	var client Client
	var err error
	if len(c.Endpoints) > 0 {
		client, err = internal.NewClient(internal.BuildDirectTarget(c.Endpoints), opts...)
	} else {
		client, err = internal.NewClient(internal.BuildDiscoveryTarget(c.Etcd.Hosts, c.Etcd.Key), opts...)
	}
	if err != nil {
		return nil, err
	}

	return &RpcClient{client: client}, nil
}

// NewClientNoAuth 新建无需鉴权的、基于etcd的rpc客户端
func NewClientNoAuth(c discovery.EtcdConf, opts ...ClientOption) (Client, error) {
	client, err := internal.NewClient(internal.BuildDiscoveryTarget(c.Hosts, c.Key), opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{client: client}, nil
}

// NewClientWithTarget 新建指向目标地址的代理客户端
func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	return internal.NewClient(target, opts...)
}
