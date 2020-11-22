package internal

import (
	"context"
	"errors"
	"fmt"
	"git.zc0901.com/go/god/rpc/internal/balancer/p2c"
	"git.zc0901.com/go/god/rpc/internal/client_interceptors"
	"git.zc0901.com/go/god/rpc/internal/resolver"
	"google.golang.org/grpc"
	"strings"
	"time"
)

const (
	dialTimeout = time.Second * 2
	separator   = '/'
)

type (
	ClientOptions struct {
		Timeout     time.Duration
		DialOptions []grpc.DialOption
	}

	ClientOption func(options *ClientOptions)

	client struct {
		conn *grpc.ClientConn
	}
)

func (c *client) dial(server string, opts ...ClientOption) error {
	options := c.buildDialOptions(opts...)
	timeCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(timeCtx, server, options...)
	if err != nil {
		service := server
		if errors.Is(err, context.DeadlineExceeded) {
			pos := strings.LastIndexByte(server, separator)
			if 0 < pos && pos < len(server)-1 {
				service = server[pos+1:]
			}
		}
		return fmt.Errorf("rpc dial: %s, 错误：%s, 请确保rpc服务 %s 是否已启动，如使用etcd也需确保已启动",
			server, err.Error(), service)
	}

	c.conn = conn
	return nil
}

func (c *client) buildDialOptions(opts ...ClientOption) []grpc.DialOption {
	var cliOpts ClientOptions
	for _, opt := range opts {
		opt(&cliOpts)
	}

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		WithUnaryClientInterceptors(
			client_interceptors.TraceInterceptor,                    // 线路跟踪
			client_interceptors.DurationInterceptor,                 // 慢查询日志
			client_interceptors.BreakerInterceptor,                  // 自动熔断
			client_interceptors.PrometheusInterceptor,               // 监控报警
			client_interceptors.TimeoutInterceptor(cliOpts.Timeout), // 超时控制
		),
	}

	return append(options, cliOpts.DialOptions...)
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func init() {
	// 注册服务直连和服务发现两种模式
	resolver.RegisterResolver()
}

func NewClient(target string, opts ...ClientOption) (*client, error) {
	var cli client
	opts = append([]ClientOption{WithDialOption(grpc.WithBalancerName(p2c.Name))}, opts...)
	if err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return &cli, nil
}

func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, opt)
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

func WithUnaryClientInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, WithUnaryClientInterceptors(interceptor))
	}
}
