package internal

import (
	"fmt"
	"git.zc0901.com/go/god/rpc/internal/resolver"
	"strings"
)

func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectSchema,
		strings.Join(endpoints, resolver.EndpointsSep))
}

func BuildDiscoveryTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscoverySchema,
		strings.Join(endpoints, resolver.EndpointsSep), key)
}
