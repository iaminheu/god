package api

import (
	"time"

	"git.zc0901.com/go/god/lib/service"
)

type (
	// PrivateKeyConf 私钥配置
	PrivateKeyConf struct {
		Fingerprint string // 指纹配置
		KeyFile     string // 秘钥配置
	}

	// SignatureConf 签名配置
	SignatureConf struct {
		Strict      bool             `json:",default=false"` // 签名是否为严格模式，若是则签名秘钥(PrivateKeys)必填
		Expire      time.Duration    `json:",default=1h"`    // 签名有效期，默认1小时
		PrivateKeys []PrivateKeyConf // 签名秘钥相关配置
	}

	// Conf 接口配置
	//
	// 为什么不命名为 Conf或Config，因为需要考虑如下用法：
	//
	// type Config struct {
	// 		rpc.RpcConf
	// 		api.Conf
	// }
	//
	// 所以如果命名为 Conf 的话，Config中就会有两个 Conf 了。
	Conf struct {
		service.Conf               // service 配置
		Host         string        `json:",default=0.0.0.0"` // http 监听ip，默认0.0.0.0
		Port         int           // http 监听端口，必填
		CertFile     string        `json:",optional"`                           // http 证书文件，可选
		KeyFile      string        `json:",optional"`                           // http 私钥文件，可选
		Verbose      bool          `json:",optional"`                           // 是否打印详细http请求日志
		MaxConns     int           `json:",default=10000"`                      // http同时可接受最大请求数（限流数），默认10000
		MaxBytes     int64         `json:",default=1048576,range=[0:33554432]"` // 最大文件上传大小（字节）默认1Mb，最大32Mb
		Timeout      int64         `json:",default=3000"`                       // 请求超时时间，默认3000毫秒
		CpuThreshold int64         `json:",default=900,range=[0:1000]"`         // cpu降载阈值，默认900，可选范围0-1000
		Signature    SignatureConf `json:",optional"`                           // 签名配置
	}
)
