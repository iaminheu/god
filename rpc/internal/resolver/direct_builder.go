package resolver

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

// directBuilder 服务直连构建器
type directBuilder struct{}

// Build 构建服务直连解析器
func (b *directBuilder) Build(target resolver.Target, conn resolver.ClientConn, options resolver.BuildOptions) (
	resolver.Resolver, error) {
	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.Endpoint, func(r rune) bool {
		return r == EndpointSepChar
	})

	for _, addr := range subset(endpoints, subsetSize) {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	conn.UpdateState(resolver.State{Addresses: addrs})

	return &nopResolver{conn: conn}, nil
}

// Schema 服务直连模式
func (b *directBuilder) Scheme() string {
	return DirectSchema
}
