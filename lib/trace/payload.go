package trace

import (
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidPayload = errors.New("无效 payload 参数")

type (
	Payload interface {
		Get(key string) string
		Set(key, value string)
	}

	httpPayload http.Header
	grpcPayload map[string][]string // grpc metadata keys 不区分大小写
)

func (h httpPayload) Get(key string) string {
	return http.Header(h).Get(key)
}

func (h httpPayload) Set(key, val string) {
	http.Header(h).Set(key, val)
}

func (g grpcPayload) Get(key string) string {
	if payload, ok := g[strings.ToLower(key)]; ok && len(payload) > 0 {
		return payload[0]
	} else {
		return ""
	}
}

func (g grpcPayload) Set(key, val string) {
	key = strings.ToLower(key)
	g[key] = append(g[key], val)
}
