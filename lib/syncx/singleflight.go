package syncx

import (
	"sync"
)

type (
	// SharedCalls 是 SingleFlight 的别名。
	// 已废弃: 使用 SingleFlight。
	SharedCalls = SingleFlight

	// SingleFlight 让相同 key 的并发调用，共享同一调用结果。
	// 如，A 先调用 X，B 后调用了 X，则先执行 A，然后 B 共享该结果。
	SingleFlight interface {
		Do(key string, fn func() (val interface{}, err error)) (val interface{}, hit bool, err error)
	}

	// call 调用，封装其等待组、调用结果和错误信息
	call struct {
		wg  sync.WaitGroup
		val interface{}
		err error
	}

	// flightGroup 共享调用组，封装了所有调用及互斥锁
	flightGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

// NewSingleFlight 返回一个 SingleFlight。
func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

// NewSharedCalls 返回一个 SingleFlight。
// Deprecated: 使用 NewSingleFlight。
func NewSharedCalls() SingleFlight {
	return NewSingleFlight()
}

// Do 获取指定 key 的调用、是否命中及错误信息，如未命中则重新调用并共享。
func (g *flightGroup) Do(key string, fn func() (interface{}, error)) (val interface{}, hit bool, err error) {
	c, hit := g.get(key)

	// 命中共享调用，则直接返回结果及错误
	if hit {
		return c.val, hit, c.err
	}

	// 未命中，则调取，并共享结果
	g.share(key, c, fn)
	return c.val, hit, c.err
}

// get 获取 key 对应的组内共享调用，如不存在则创建并返回
func (g *flightGroup) get(key string) (c *call, hit bool) {
	g.lock.Lock()

	// 如果 key 对应的调用已在组内存在，则直接返回共享结果
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	// 如不存在，则新增该 key 的调用
	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()
	return c, false
}

func (g *flightGroup) share(key string, c *call, fn func() (interface{}, error)) {
	defer func() {
		// 先删除，后完成。顺序不可反，否则其他调用可能一直等待不到计数器归零。
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
