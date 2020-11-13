package breaker

import (
	"sync"
)

var (
	lock     sync.RWMutex
	breakers = make(map[string]Breaker)
)

// Do 使用指定名称的断路器，执行请求函数。
func Do(name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.Do(req)
	})
}

// DoWithFallback 使用指定名称的断路器，若断路器允许则执行请求，反之则根据error调用备用应急函数。
func DoWithFallback(name string, req func() error, fallback Fallback) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFailback(req, fallback)
	})
}

// DoWithAcceptable 使用指定名称的断路器，若断路器允许则执行请求，若执行错误则判断可否标记为请求成功。
func DoWithAcceptable(name string, req func() error, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptable(req, acceptable)
	})
}

// DoWithFallbackAcceptable 使用指定名称的断路器，若断路器允许则执行请求，反之则根据error调用备用应急函数。
// 若请求的执行，出现错误，则判断可否标记为请求成功。。
func DoWithFallbackAcceptable(name string, req func() error, fallback Fallback, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFailbackAcceptable(req, fallback, acceptable)
	})
}

func GetBreaker(name string) Breaker {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return b
	}

	lock.Lock()
	b, ok = breakers[name]
	if !ok {
		b = NewBreaker(WithName(name))
		breakers[name] = b
	}
	lock.Unlock()

	return b
}

func NoBreakFor(name string) {
	lock.Lock()
	breakers[name] = newNoOpBreaker()
	lock.Unlock()
}

func do(name string, execute func(b Breaker) error) error {
	return execute(GetBreaker(name))
}
