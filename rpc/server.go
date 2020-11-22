package rpc

import (
	"git.zc0901.com/go/god/lib/load"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/netx"
	"git.zc0901.com/go/god/lib/stat"
	"git.zc0901.com/go/god/rpc/internal"
	"git.zc0901.com/go/god/rpc/internal/auth"
	"git.zc0901.com/go/god/rpc/internal/server_interceptors"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
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

func NewServer(sc ServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error

	// 验证服务端配置
	if err := sc.Validate(); err != nil {
		return nil, err
	}

	// 新建监控指标，以监听端口作为监听名称
	metrics := stat.NewMetrics(sc.ListenOn)

	// 新建内部RPC服务器
	var server internal.Server
	if sc.HasEtcd() {
		listenOn := figureOutListenOn(sc.ListenOn)
		server, err = internal.NewPubServer(sc.Etcd.Hosts, sc.Etcd.Key, listenOn, internal.WithMetrics(metrics))
		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(sc.ListenOn, internal.WithMetrics(metrics))
	}

	server.SetName(sc.Name)
	if err = setupInterceptors(server, sc, metrics); err != nil {
		return nil, err
	}

	// 新建对外RPC服务器
	rpcServer := &RpcServer{server: server, register: register}
	if err = sc.Setup(); err != nil {
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

func (rs *RpcServer) Stop() {
	logx.Close()
}

func setupInterceptors(server internal.Server, sc ServerConf, metrics *stat.Metrics) error {
	// 自动降载（负载卸流拦截器）
	if sc.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(sc.CpuThreshold))
		server.AddUnaryInterceptors(server_interceptors.UnaryShedderInterceptor(shedder, metrics))
	}

	// 超时控制（超时拦截器）
	if sc.Timeout > 0 {
		server.AddUnaryInterceptors(server_interceptors.UnaryTimeoutInterceptor(
			time.Duration(sc.Timeout) * time.Millisecond))
	}

	// 调用鉴权（鉴权拦截器）
	if sc.Auth {
		authenticator, err := auth.NewAuthenticator(sc.Redis.NewRedis(), sc.Redis.Key, sc.StrictControl)
		if err != nil {
			return err
		}

		server.AddUnaryInterceptors(server_interceptors.UnaryAuthorizeInterceptor(authenticator))
		server.AddStreamInterceptors(server_interceptors.StreamAuthorizeInterceptor(authenticator))
	}

	return nil
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")

	// 未传监听地址
	if len(fields) == 0 {
		return listenOn
	}

	// 传host且不是0.0.0.0
	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	} else {
		return strings.Join(append([]string{ip}, fields[1:]...), ":")
	}
}
