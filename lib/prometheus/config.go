package prometheus

type PromConf struct {
	Host string `json:",optional"`         // prometheus 监听ip，默认为空即不开启监控
	Port int    `json:",default=9101"`     // prometheus 监听端口
	Path string `json:",default=/metrics"` // Prometheus 上报地址
}
