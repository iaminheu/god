package service

import (
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/prometheus"
	"git.zc0901.com/go/god/lib/stat"
	"github.com/prometheus/common/log"
)

const (
	DevMode  = "dev"  // 开发模式
	TestMode = "test" // 测试环境
	PreMode  = "pre"  // 预发布模式
	ProMode  = "pro"  // 生产模式
)

type ServiceConf struct {
	Name       string
	LogConf    logx.LogConf
	Mode       string              `json:",default=pro,options=dev|test|pre|pro"`
	MetricsUrl string              `json:",optional"`
	PromConf   prometheus.PromConf `json:",optional"`
}

func (sc ServiceConf) MustSetup() {
	if err := sc.Setup(); err != nil {
		log.Fatal(err)
	}
}

func (sc ServiceConf) Setup() error {
	if len(sc.LogConf.ServiceName) == 0 {
		sc.LogConf.ServiceName = sc.Name
	}
	if err := logx.Setup(sc.LogConf); err != nil {
		return err
	}

	sc.initMode()
	prometheus.StartAgent(sc.PromConf)
	if len(sc.MetricsUrl) > 0 {
		// TODO
	}
}

func (sc ServiceConf) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, PreMode:
		load.Disable()
		stat.SetReporter(nil)
	}
}
