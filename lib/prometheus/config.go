package prometheus

type PromConf struct {
	Host string `json:",optional"` // Prometheus Server，默认为空即不开启监控
	Port int    `json:",default=9101"`
	Path string `json:",default=/metrics"` // Prometheus拉去指标的路径
}
