package internal

import (
	"git.zc0901.com/go/god/lib/discovery"
	"git.zc0901.com/go/god/lib/netx"
	"os"
	"strings"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

func NewPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		listenOn := figureOutListenOn(listenOn)
		pubClient := discovery.NewPublisher(etcdEndpoints, etcdKey, listenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}

type keepAliveServer struct {
	registerEtcd func() error
	Server
}

func (s keepAliveServer) Start(register RegisterFn) error {
	if err := s.registerEtcd(); err != nil {
		return err
	}

	return s.Server.Start(register)
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
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
