package p2c

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/syncx"
	"git.zc0901.com/go/god/lib/timex"
	"git.zc0901.com/go/god/rpc/internal/codes"
	"google.golang.org/grpc/balancer"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
)

type p2cPicker struct {
	conns []*subConn
	r     *rand.Rand
	stamp *syncx.AtomicDuration
	lock  sync.Mutex
}

func (p *p2cPicker) Pick(ctx context.Context, info balancer.PickInfo) (
	conn balancer.SubConn, done func(doneInfo balancer.DoneInfo), err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var chosen *subConn
	switch len(p.conns) {
	case 0: // 没有连接
		return nil, nil, balancer.ErrNoSubConnAvailable
	case 1: // 一个连接
		chosen = p.choose(p.conns[0], nil)
	case 2: // 2个连接
		chosen = p.choose(p.conns[0], p.conns[1])
	default: // 2个以上连接，随机选举两个健康的并在其中选举1个
		var c1, c2 *subConn
		for i := 0; i < pickTimes; i++ {
			a := p.r.Intn(len(p.conns))
			b := p.r.Intn(len(p.conns) - 1)
			if b >= a {
				b++
			}
			c1 = p.conns[a]
			c2 = p.conns[b]
			if c1.healthy() && c2.healthy() {
				break
			}
		}
		chosen = p.choose(c1, c2)
	}

	atomic.AddInt64(&chosen.inflight, 1) // 飞行中+1
	atomic.AddInt64(&chosen.requests, 1) // 请求数+1
	return chosen.conn, p.buildDoneFunc(chosen), nil
}

// choose 从两个连接中选举一个用于使用
func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	startTime := int64(timex.Now())
	if c2 == nil {
		atomic.StoreInt64(&c1.pickTime, startTime)
		return c1
	}

	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	pickTime := atomic.LoadInt64(&c2.pickTime)
	if startTime-pickTime > forcePickTime && atomic.CompareAndSwapInt64(&c2.pickTime, pickTime, startTime) {
		return c2
	} else {
		atomic.StoreInt64(&c1.pickTime, startTime)
		return c1
	}
}

// buildDoneFunc 构建完成连接函数
func (p *p2cPicker) buildDoneFunc(conn *subConn) func(doneInfo balancer.DoneInfo) {
	startTime := int64(timex.Now())
	return func(doneInfo balancer.DoneInfo) {
		atomic.AddInt64(&conn.inflight, -1)
		now := timex.Now()
		lastTime := atomic.SwapInt64(&conn.lastTime, int64(now))
		duration := int64(now) - lastTime
		if duration < 0 {
			duration = 0
		}
		w := math.Exp(float64(-duration) / float64(decayTime))
		lag := int64(now) - startTime
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadUint64(&conn.lag)
		if oldLag == 0 {
			w = 0
		}
		atomic.StoreUint64(&conn.lag, uint64(float64(oldLag)*w+float64(lag)*(1-w)))
		success := initSuccess
		if doneInfo.Err != nil && !codes.Acceptable(doneInfo.Err) {
			success = 0
		}
		oldSuccess := atomic.LoadUint64(&conn.success)
		atomic.StoreUint64(&conn.success, uint64(float64(oldSuccess)*w+float64(success)*(1-w)))

		stamp := p.stamp.Load()
		if now-stamp >= logInterval {
			if p.stamp.CompareAndSwap(stamp, now) {
				p.logStats()
			}
		}
	}
}

func (p *p2cPicker) logStats() {
	var stats []string

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, conn := range p.conns {
		stats = append(stats, fmt.Sprintf("conn: %s, load: %d, reqs: %d",
			conn.addr.Addr, conn.load(), atomic.SwapInt64(&conn.requests, 0)))
	}

	logx.Statf("p2c - %s", strings.Join(stats, "; "))
}
