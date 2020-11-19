package resolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
)

const (
	DirectSchema    = "direct"    // 服务直连模式
	DiscoverySchema = "discovery" // 服务发现模式

	EndpointSepChar = ','
	subsetSize      = 32
)

var (
	EndpointsSep = fmt.Sprintf("%c", EndpointSepChar)
	dirBuilder   directBuilder
	disBuilder   discoveryBuilder
)

type nopResolver struct {
	conn resolver.ClientConn
}

func RegisterResolver() {
	resolver.Register(&dirBuilder)
	resolver.Register(&disBuilder)
}

func (r *nopResolver) Close() {

}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {

}
