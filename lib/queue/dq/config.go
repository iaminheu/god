package dq

import "god/lib/store/redis"

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
