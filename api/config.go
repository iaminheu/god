package api

import (
	"git.zc0901.com/go/god/lib/service"
	"time"
)

type (
	// 私钥配置
	PrivateKeyConf struct {
		Fingerprint string // 指纹
		KeyFile     string // 私钥文件路径
	}

	// 签名配置
	SignatureConf struct {
		Strict      bool          `json:",default=false"`
		Expire      time.Duration `json:",default=1h"`
		PrivateKeys []PrivateKeyConf
	}

	// 接口配置
	//
	// 为什么不命名为 Conf或Config，因为需要考虑如下用法：
	//
	// type Config struct {
	// 		rpc.RpcConf
	// 		api.ApiConf
	// }
	//
	// 所以如果命名为 Conf 的话，Config中就会有两个 Conf 了。
	ApiConf struct {
		service.ServiceConf
		Host         string `json:",default=0.0.0.0"`
		Port         int
		CertFile     string        `json:",optional"`
		KeyFile      string        `json:",optional"`
		Verbose      bool          `json:",optional"`
		MaxConns     int           `json:",default=10000"`
		MaxBytes     int64         `json:",default=1048576,range=[0:33554432]"` // 最大文件上传大小（字节）默认1Mb，最大32Mb
		Timeout      int64         `json:",default=3000"`                       // 超时时间（毫秒）
		CpuThreshold int64         `json:",default=900,range=[0:1000]"`
		Signature    SignatureConf `json:",optional"`
	}
)
