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
	Name       string            // 服务名称
	Log        logx.LogConf      // 日志配置
	Mode       string            `json:",default=pro,options=dev|test|pre|pro"` // 服务环境，dev-开发环境，test-测试环境，pre-预发环境，pro-正式环境
	MetricsUrl string            `json:",optional"`                             // 指标上报接口地址，该地址需要支持 post json 即可
	Prometheus prometheus.Config `json:",optional"`                             // 普罗米修斯配置
}

func (sc ServiceConf) MustSetup() {
	if err := sc.Setup(); err != nil {
		log.Fatal(err)
	}
}

// Setup 设置并初始化服务配置（初始化启动模式、普罗米修斯代理、统计输出器等）
func (sc ServiceConf) Setup() error {
	if len(sc.Log.ServiceName) == 0 {
		sc.Log.ServiceName = sc.Name
	}

	// 初始化日志
	if err := logx.Setup(sc.Log); err != nil {
		return err
	}

	// 非生产模式禁用负载均衡和日志汇报
	sc.initMode()

	// 方式一：启动普罗米修斯http服务端口
	prometheus.StartAgent(sc.Prometheus)

	// 方式二：设置统计报告书写器（写入普罗米修斯）
	if len(sc.MetricsUrl) > 0 {
		stat.SetReportWriter(stat.NewRemoteWriter(sc.MetricsUrl))
	}

	return nil
}

func (sc ServiceConf) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, PreMode:
		// 开发、测试、预发布模式，不启用负载均衡和统计上报
		load.Disable()
		stat.SetReporter(nil)
	}
}
