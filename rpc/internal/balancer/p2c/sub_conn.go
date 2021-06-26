package p2c

import (
	"math"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type subConn struct {
	addr     resolver.Address
	conn     balancer.SubConn
	lag      uint64 // 延迟
	inflight int64  // 飞行中的数量
	requests int64  // 请求数
	success  uint64 // 请求成功数
	last     int64  // 该连接最后使用时间
	pick     int64  // 该连接被选举时间
}

func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

func (c *subConn) load() int64 {
	// +1 以防被零除
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		return penalty
	} else {
		return load
	}
}
