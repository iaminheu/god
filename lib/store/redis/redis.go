package redis

import (
	"git.zc0901.com/go/god/lib/breaker"
	"github.com/go-redis/redis"
	"time"
)

const (
	ClusterMode    = "cluster"
	StandaloneMode = "standalone"

	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
	slowThreshold   = 100 * time.Millisecond
)

type (
	Redis struct {
		Addr     string
		Mode     string
		Password string
		brk      breaker.Breaker
	}

	Client interface {
		redis.Cmdable
	}

	// GeoLocation is used with GeoAdd to add geospatial location.
	GeoLocation = redis.GeoLocation
	// GeoRadiusQuery is used with GeoRadius to query geospatial index.
	GeoRadiusQuery = redis.GeoRadiusQuery
	GeoPos         = redis.GeoPos

	Pipeliner = redis.Pipeliner

	Pair struct {
		Key   string
		Score int64
	}
)

func NewRedis(addr, mode string, password ...string) *Redis {
	// 为了支持不提供 password 的情况
	var pwd string
	for _, v := range password {
		pwd = v
	}

	return &Redis{
		Addr:     addr,
		Mode:     mode,
		Password: pwd,
		brk:      breaker.NewBreaker(),
	}
}
