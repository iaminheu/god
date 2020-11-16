package service

import (
	"git.zc0901.com/go/god/lib/load"
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

type Config struct {
	Name       string
	LogConf    logx.LogConf
	Mode       string              `json:",default=pro,options=dev|test|pre|pro"`
	MetricsUrl string              `json:",optional"`
	PromConf   prometheus.PromConf `json:",optional"`
}

func (sc Config) MustSetup() {
	if err := sc.Setup(); err != nil {
		log.Fatal(err)
	}
}

func (sc Config) Setup() error {
	if len(sc.LogConf.ServiceName) == 0 {
		sc.LogConf.ServiceName = sc.Name
	}
	if err := logx.Setup(sc.LogConf); err != nil {
		return err
	}

	// 非生产模式禁用负载均衡和日志汇报
	sc.initMode()

	// 启动普罗米修斯http服务端口
	prometheus.StartAgent(sc.PromConf)

	// 设置统计报告书写器（写入普罗米修斯）
	if len(sc.MetricsUrl) > 0 {
		stat.SetReportWriter(stat.NewRemoteWriter(sc.MetricsUrl))
	}

	return nil
}

func (sc Config) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, PreMode:
		// 开发、测试、预发布模式，不启用负载均衡和统计上报
		load.Disable()
		stat.SetReporter(nil)
	}
}
