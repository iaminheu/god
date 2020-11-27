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
	conf.MustLoad("/Users/zs/git/god/example/config/load_from_yaml/date.yml", th)
	logx.Infof("%+v", th)
}
