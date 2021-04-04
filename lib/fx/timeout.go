package fx

import (
	"context"
	"time"
)

var (
	ErrCanceled = context.Canceled
	ErrTimeout  = context.DeadlineExceeded
)

type TimeoutOption func() context.Context

func DoWithTimeout(fn func() error, timeout time.Duration, opts ...TimeoutOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	timeoutCtx, cancel := context.WithTimeout(parentCtx, timeout)
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
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

// WithTimeout 设置带超时时间的上下文
func WithContext(ctx context.Context) TimeoutOption {
	return func() context.Context {
		return ctx
	}
}
