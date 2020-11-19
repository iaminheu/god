package internal

import (
	"context"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

// EtcdClient 用于服务发现的 Etcd 客户端
type EtcdClient interface {
	ActiveConnection() *grpc.ClientConn
	Close() error
	Ctx() context.Context
	Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error)
	KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error)
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error)
	Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan
}
