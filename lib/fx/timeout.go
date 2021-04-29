package fx

import (
	"context"
	"time"
)

var (
	// ErrCanceled 是上下文被取消时返回的错误
	ErrCanceled = context.Canceled

	// ErrTimeout 是上下文截止时间超时返回的错误
	ErrTimeout = context.DeadlineExceeded
)

type DoOption func() context.Context

// DoWithTimeout 带超时控制执行函数fn
func DoWithTimeout(fn func() error, timeout time.Duration, opts ...DoOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		done <- fn()
	}()

	// 处理结果和panic
	select {
	case p := <-panicChan:
		panic(p)
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithTimeout 设置带超时时间的上下文
func WithContext(ctx context.Context) DoOption {
	return func() context.Context {
		return ctx
	}
}
