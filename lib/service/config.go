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

type ServiceConf struct {
	Name       string
	LogConf    logx.LogConf
	Mode       string              `json:",default=pro,options=dev|test|pre|pro"`
	MetricsUrl string              `json:",optional"`
	PromConf   prometheus.PromConf `json:",optional"`
}

func (c ServiceConf) MustSetup() {
	if err := c.Setup(); err != nil {
		log.Fatal(err)
	}
}

// 设置并初始化服务配置（初始化启动模式、普罗米修斯代理、统计输出器等）
func (c ServiceConf) Setup() error {
	if len(c.LogConf.ServiceName) == 0 {
		c.LogConf.ServiceName = c.Name
	}
	if err := logx.Setup(c.LogConf); err != nil {
		return err
	}

	// 非生产模式禁用负载均衡和日志汇报
	c.initMode()

	// 启动普罗米修斯http服务端口
	prometheus.StartAgent(c.PromConf)

	// 设置统计报告书写器（写入普罗米修斯）
	if len(c.MetricsUrl) > 0 {
		stat.SetReportWriter(stat.NewRemoteWriter(c.MetricsUrl))
	}

	return nil
}

func (c ServiceConf) initMode() {
	switch c.Mode {
	case DevMode, TestMode, PreMode:
		// 开发、测试、预发布模式，不启用负载均衡和统计上报
		load.Disable()
		stat.SetReporter(nil)
	}
}
