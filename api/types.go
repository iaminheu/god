package api

import "net/http"

type (
	// 路由
	Route struct {
		Method  string           // 路由方法
		Path    string           // 路由路径
		Handler http.HandlerFunc // 路由处理器
	}

	// 路由可选项函数
	RouteOption func(r *featuredRoutes)

	// 中间件函数（接收一个处理函数，并返回另一个处理函数）
	Middleware func(next http.HandlerFunc) http.HandlerFunc

	// jsonWebToken 设置
	jwtSetting struct {
		enabled    bool   // 是否启用jwt验证
		secret     string // jwt秘钥
		prevSecret string // 上一个jwt秘钥
	}

	// 签名设置
	signatureSetting struct {
		SignatureConf
		enabled bool // 是否启用签名校验
	}

	// 特色路由，支持高优先级、jwt令牌校验、签名校验
	featuredRoutes struct {
		priority  bool             // 带有高优先级的路由
		jwt       jwtSetting       // JWT 鉴权
		signature signatureSetting // 签名校验
		routes    []Route
	}
)
