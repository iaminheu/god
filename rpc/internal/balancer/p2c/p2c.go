package p2c

import (
	"math"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	// Name p2c平衡器的名字
	Name = "p2c_ewma"

	decayTime       = int64(time.Second * 10) // 默认衰减时间
	forcePickTime   = int64(time.Second)      // 强制选举时间
	initSuccess     = 1000
	throttleSuccess = initSuccess / 2      // 健康检测阈值
	penalty         = int64(math.MaxInt32) // 最大惩罚值
	pickTimes       = 3
	logInterval     = time.Minute
)

var emptyPickResult balancer.PickResult

func init() {
	balancer.Register(newBuilder())
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(pickerBuilder), base.Config{HealthCheck: true})
}
