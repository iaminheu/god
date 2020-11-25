package main

import (
	"fmt"
	"git.zc0901.com/go/god/lib/queue/dq"
	"git.zc0901.com/go/god/lib/store/redis"
)

func main() {
	consumer := dq.NewConsumer(dq.Conf{
		Beanstalks: []dq.Beanstalk{
			{
				Endpoint: "localhost:11300",
				Tube:     "tube",
			},
			{
				Endpoint: "localhost:11301",
				Tube:     "tube",
			},
		},
		Redis: redis.Conf{
			Host: "192.168.0.17:6379",
			Mode: redis.StandaloneMode,
		},
	})

	consumer.Consume(func(body []byte) {
		fmt.Println(body)
	})
}
