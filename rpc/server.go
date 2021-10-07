package rpc

import (
	"log"
	"time"

	"git.zc0901.com/go/god/lib/load"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stat"
	"git.zc0901.com/go/god/rpc/internal"
	"git.zc0901.com/go/god/rpc/internal/auth"
	"git.zc0901.com/go/god/rpc/internal/serverinterceptors"
	"google.golang.org/grpc"
)

type RpcServer struct {
	server   internal.Server
	register internal.RegisterFn
}

func MustNewServer(sc ServerConf, register internal.RegisterFn) *RpcServer {
	server, err := NewServer(sc, register)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func NewServer(c ServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error

	// 验证服务端配置
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// 新建监控指标，以监听端口作为监听名称
	metrics := stat.NewMetrics(c.ListenOn)

	// 新建内部RPC服务器
	var server internal.Server
	if c.HasEtcd() {
		server, err = internal.NewPubServer(c.Etcd.Hosts, c.Etcd.Key, c.ListenOn, internal.WithMetrics(metrics))
		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(c.ListenOn, internal.WithMetrics(metrics))
	}

	server.SetName(c.Name)
	if err = setupInterceptors(server, c, metrics); err != nil {
		return nil, err
	}

	// 新建对外RPC服务器
	rpcServer := &RpcServer{server: server, register: register}
	if err = c.Setup(); err != nil {
		return nil, err
	}

	return rpcServer, nil
}

func (rs *RpcServer) AddOptions(options ...grpc.ServerOption) {
	rs.server.AddOptions(options...)
}

func (rs *RpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	rs.server.AddUnaryInterceptors(interceptors...)
}

func (rs *RpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	rs.server.AddStreamInterceptors(interceptors...)
}

func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}

// Stop 关闭RPC服务器
func (rs *RpcServer) Stop() {
	logx.Close()
}

func setupInterceptors(server internal.Server, c ServerConf, metrics *stat.Metrics) error {
	// 自动降载（负载泄流拦截器）
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		server.AddUnaryInterceptors(serverinterceptors.UnaryShedderInterceptor(shedder, metrics))
	}

	// 超时控制（超时拦截器）
	if c.Timeout > 0 {
		server.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout) * time.Millisecond))
	}

	// 调用鉴权（鉴权拦截器）
	if c.Auth {
		authenticator, err := auth.NewAuthenticator(c.Redis.NewRedis(), c.Redis.Key, c.StrictControl)
		if err != nil {
			return err
		}

		server.AddStreamInterceptors(serverinterceptors.StreamAuthorizeInterceptor(authenticator))
		server.AddUnaryInterceptors(serverinterceptors.UnaryAuthorizeInterceptor(authenticator))
	}

	return nil
}
