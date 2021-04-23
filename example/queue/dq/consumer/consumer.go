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
				Endpoint: "dev:11300",
				Tube:     "dhome-sms-login",
			},
			{
				Endpoint: "dev:11301",
				Tube:     "dhome-sms-login",
			},
		},
		Redis: redis.Conf{
			Host: "192.168.0.17:6382",
			Mode: redis.StandaloneMode,
		},
	})

	consumer.Consume(func(body []byte) {
		//time.Sleep(1* time.Second)
		fmt.Println(body)
	})
}
