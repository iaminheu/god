package conf

type (
	Option func(opt *options)

	options struct {
		env bool
	}
)

// UseEnv 是否加载环境配置信息
func UseEnv() Option {
	return func(opt *options) {
		opt.env = true
	}
}
