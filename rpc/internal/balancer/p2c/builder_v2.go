package p2c

import (
	"git.zc0901.com/go/god/lib/syncx"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
	"time"
)

type pickerBuilderV2 struct{}

func (b *pickerBuilderV2) Build(info base.PickerBuildInfo) balancer.V2Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	var conns []*subConn
	for conn, info := range info.ReadySCs {
		conns = append(conns, &subConn{addr: info.Address, conn: conn, success: initSuccess})
	}

	return &v2P2cPicker{
		conns: conns,
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		stamp: syncx.NewAtomicDuration(),
	}
}
