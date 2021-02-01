package main

import (
	"git.zc0901.com/go/god/lib/conf"
	"git.zc0901.com/go/god/lib/logx"
	"time"
)

type TimeHolder struct {
	Date time.Time `json:"date"`
}

func main() {
	th := &TimeHolder{}
	err := conf.LoadConfig("./date.yml", th)
	if err != nil {
		logx.Error(err)
	}
	logx.Infof("%+v", th)
}
