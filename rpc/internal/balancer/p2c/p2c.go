package p2c

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"time"
)

const (
	Name = "p2c_ewma"

	initSuccess     = 1000
	throttleSuccess = initSuccess / 2      // 健康检测阈值
	penalty         = int64(math.MaxInt32) // 最大惩罚值

	forcePickTime = int64(time.Second) // 强制选举时间
	pickTimes     = 3
	decayTime     = int64(time.Second * 10) // 默认衰减时间
	logInterval   = time.Minute
)

func init() {
	balancer.Register(newBuilder())
	//balancer.Register(newBuilderV2())
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(pickerBuilder))
}

func newBuilderV2() balancer.Builder {
	return base.NewBalancerBuilderV2(Name, new(pickerBuilderV2), base.Config{
		HealthCheck: true,
	})
}
