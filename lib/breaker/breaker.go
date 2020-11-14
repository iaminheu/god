package breaker

import (
	"errors"
	"god/lib/stringx"
)

const (
	StateClosed State = iota // 断路器关闭 0
	StateOpen                // 断路器打开 1
)

var ErrServiceUnavailable = errors.New("断路器已打开，服务不可用")

type (
	State      = int32                     // 断路器状态
	Request    func() error                // 待执行的请求
	Acceptable func(reqError error) bool   // 判断错误是否可接受的函数
	Fallback   func(acceptErr error) error // 备用函数
	Option     func(b *breaker)            // 断路器可选项应用函数

	Breaker interface {
		// Name 是 netflixBreaker 断路器的名称
		Name() string

		// Allow 检查请求是否被允许。
		// 若允许，则返回 Promise，调用者需在成功时调用 promise.Resolve()，失败时调用 promise.Reject()。
		// 若不允许，则返回 ErrServiceUnavailable。
		Allow() (Promise, error)

		// Do 若断路器允许则执行请求，否则返回错误。
		Do(req Request) error

		// DoWithFailback 若断路器允许则执行请求，反之则调用备用函数，再之返回错误。
		DoWithFailback(req Request, fallback Fallback) error

		// DoWithAcceptable 若断路器允许则执行请求，反之则返回错误，并判断错误可否标记为成功请求。
		DoWithAcceptable(req Request, acceptable Acceptable) error

		// DoWithFailbackAcceptable 若断路器允许则执行请求，反之则调用备用函数，再之返回错误并判断错误可否标记为成功请求。
		DoWithFailbackAcceptable(req Request, fallback Fallback, acceptable Acceptable) error
	}

	Promise interface {
		Accept()
		Reject(reason string)
	}

	throttle interface {
		allow() (Promise, error)
		doReq(req Request, fallback Fallback, acceptable Acceptable) error
	}

	breaker struct {
		name string
		throttle
	}
)

func NewBreaker(opts ...Option) Breaker {
	var b breaker
	for _, opt := range opts {
		opt(&b)
	}
	if len(b.name) == 0 {
		b.name = stringx.Rand()
	}
	b.throttle = newLoggedThrottle(b.name, newGoogleBreaker())
	return &b
}

func WithName(name string) Option {
	return func(b *breaker) {
		b.name = name
	}
}

func (b breaker) Name() string {
	return b.name
}

func (b breaker) Allow() (Promise, error) {
	return b.throttle.allow()
}

func (b breaker) Do(req Request) error {
	return b.throttle.doReq(req, nil, defaultAcceptable)
}

func (b breaker) DoWithFailback(req Request, fallback Fallback) error {
	return b.throttle.doReq(req, fallback, defaultAcceptable)
}
func (b breaker) DoWithAcceptable(req Request, acceptable Acceptable) error {
	return b.throttle.doReq(req, nil, acceptable)
}

func (b breaker) DoWithFailbackAcceptable(req Request, fallback Fallback, acceptable Acceptable) error {
	return b.throttle.doReq(req, fallback, acceptable)
}

func defaultAcceptable(err error) bool {
	return err == nil
}
