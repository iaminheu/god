package internal

import (
	"context"
	"errors"
	"fmt"
	"git.zc0901.com/go/god/rpc/internal/balancer/p2c"
	"git.zc0901.com/go/god/rpc/internal/clientinterceptors"
	"git.zc0901.com/go/god/rpc/internal/resolver"
	"google.golang.org/grpc"
	"strings"
	"time"
)

const (
	dialTimeout = time.Second * 3
	separator   = '/'
)

type (
	// Client 包装 Conn 方法的接口
	Client interface {
		Conn() *grpc.ClientConn
	}

	// ClientOptions 是RPC客户端选择项
	ClientOptions struct {
		Timeout     time.Duration
		DialOptions []grpc.DialOption
	}

	// ClientOption 是自定义客户端选项的方法
	ClientOption func(options *ClientOptions)

	client struct {
		conn *grpc.ClientConn
	}
)

func init() {
	// 注册服务直连和服务发现两种模式
	resolver.RegisterResolver()
}

// NewClient 返回RPC客户端
func NewClient(target string, opts ...ClientOption) (*client, error) {
	var cli client
	opts = append([]ClientOption{WithDialOption(grpc.WithBalancerName(p2c.Name))}, opts...)
	if err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return &cli, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
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
			clientinterceptors.TraceInterceptor,                    // 线路跟踪
			clientinterceptors.DurationInterceptor,                 // 慢查询日志
			clientinterceptors.BreakerInterceptor,                  // 自动熔断
			clientinterceptors.PrometheusInterceptor,               // 监控报警
			clientinterceptors.TimeoutInterceptor(cliOpts.Timeout), // 超时控制
		),
	}

	return append(options, cliOpts.DialOptions...)
}

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
		return fmt.Errorf("rpc dial: %s, 错误：%s, 请确保rpc服务 %s 已启动，如使用etcd也需确保已启动",
			server, err.Error(), service)
	}

	c.conn = conn
	return nil
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
