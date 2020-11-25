package main

import (
	"fmt"
	"git.zc0901.com/go/god/lib/queue/dq"
	"os"
	"strconv"
	"time"
)

func main() {
	producer := dq.NewProducer([]dq.Beanstalk{
		{
			Endpoint: "localhost:11300",
			Tube:     "tube",
		},
		{
			Endpoint: "localhost:11301",
			Tube:     "tube",
		},
	})

	for i := 0; i < 100000; i++ {
		id, err := producer.Delay([]byte(strconv.Itoa(i)), time.Second*5)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("job id", id)
	}
}
