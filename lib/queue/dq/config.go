package dq

import "git.zc0901.com/go/god/lib/store/redis"

type (
	Beanstalk struct {
		Endpoint string
		Tube     string
	}

	Conf struct {
		Beanstalks []Beanstalk
		Redis      redis.Conf
	}
)
