package internal

import "git.zc0901.com/go/god/lib/discovery"

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

func NewPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubClient := discovery.NewPublisher(etcdEndpoints, etcdKey, listenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}
