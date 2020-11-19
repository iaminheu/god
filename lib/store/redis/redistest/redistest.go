package redistest

import (
	"git.zc0901.com/go/god/lib/lang"
	"git.zc0901.com/go/god/lib/store/redis"
	"github.com/alicebob/miniredis/v2"
	"time"
)

func CreateRedis() (r *redis.Redis, clean func(), err error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	return redis.NewRedis(mr.Addr(), redis.StandaloneMode), func() {
		ch := make(chan lang.PlaceholderType)
		go func() {
			mr.Close()
			close(ch)
		}()

		select {
		case <-ch:
		case <-time.After(time.Second):
		}
	}, nil
}
