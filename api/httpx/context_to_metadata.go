package httpx

import (
	"git.zc0901.com/go/god/api/metadata"
	"git.zc0901.com/go/god/lib/netx"
	"net/http"
	"strings"
)

// 结合 ClientPublicIP 和 ClientIP 获取用户的真实 ip
func remoteIP(r *http.Request) string {
	ip := netx.ClientPublicIP(r)
	if ip == "" {
		ip = netx.ClientIP(r)
	}
	return ip
}

func remotePort(r *http.Request) string {
	part := strings.Split(r.Host, ":")
	if len(part) == 2 {
		return part[1]
	}
	return ""
}

func ContextToMetadata(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		md := metadata.MD{
			metadata.RemoteIP:   remoteIP(req),
			metadata.RemotePort: remotePort(req),
		}
		for rawKey := range req.Header {
			md[rawKey] = req.Header.Get(rawKey)
		}
		for _, cookie := range req.Cookies() {
			md[cookie.Name] = cookie.Value
		}

		ctx := metadata.NewContext(req.Context(), md)
		next(w, req.WithContext(ctx))
	}
}
