package cache

import "time"

const (
	defaultExpires         = time.Hour * 24 * 7
	defaultNotFoundExpires = time.Minute // 防缓存穿透，设置未找到记录一分钟过期
)

type (
	// Options 存储缓存选项的结构体。
	Options struct {
		Expires         time.Duration
		NotFoundExpires time.Duration
	}

	// Option 自定义缓存选项的函数。
	Option func(o *Options)
)

func newOptions(opts ...Option) Options {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}

	if o.Expires <= 0 {
		o.Expires = defaultExpires
	}
	if o.NotFoundExpires <= 0 {
		o.NotFoundExpires = defaultNotFoundExpires
	}

	return o
}

// WithExpires 返回一个自定义缓存过期时间的函数。
func WithExpires(expires time.Duration) Option {
	return func(o *Options) {
		o.Expires = expires
	}
}

// WithNotFoundExpires 返回未找到缓存记录时自定义的函数。
func WithNotFoundExpires(expires time.Duration) Option {
	return func(o *Options) {
		o.NotFoundExpires = expires
	}
}
