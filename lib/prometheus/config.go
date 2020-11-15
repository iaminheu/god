package prometheus

type PromConf struct {
	Host string `json:",optional"`
	Port string `json:",default=9101"`
	Path string `json:",default=/metrics"`
}
