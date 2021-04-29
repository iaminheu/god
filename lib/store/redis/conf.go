package redis

import "errors"

var (
	ErrEmptyHost = errors.New("redis host 为空")
	ErrEmptyMode = errors.New("redis mode 为空")
	ErrEmptyKey  = errors.New("redis key 为空")
)

type (
	Conf struct {
		Host     string // redis地址
		Mode     string `json:",default=standalone,options=standalone|cluster"` // 默认单点模式，可选集群模式
		Password string `json:",optional"`                                      // redis密码
		TLS      bool   `json:",default=false,options=true|false"`
	}

	// KeyConf 带有指定key的redis配置，一般用于rpc调用鉴权
	KeyConf struct {
		Conf
		Key string `json:",optional"`
	}
)

// 新建 Redis 示例
func (c Conf) NewRedis() *Redis {
	var opts []Option
	if c.Mode == ClusterMode {
		opts = append(opts, Cluster())
	}
	if len(c.Password) > 0 {
		opts = append(opts, WithPassword(c.Password))
	}
	if c.TLS {
		opts = append(opts, WithTLS())
	}

	return New(c.Host, opts...)
}

func (c Conf) Validate() error {
	if len(c.Host) == 0 {
		return ErrEmptyHost
	}

	if len(c.Mode) == 0 {
		return ErrEmptyMode
	}

	return nil
}

func (kc KeyConf) Validate() error {
	if err := kc.Conf.Validate(); err != nil {
		return err
	}

	if len(kc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
