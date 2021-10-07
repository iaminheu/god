package syncx

import "sync"

// Locker 对操作加锁/解锁
type Locker struct {
	lock sync.Mutex
}

// Guard 加锁保护
func (b *Locker) Guard(fn func()) {
	b.lock.Lock()
	defer b.lock.Unlock()
	fn()
}

// Guard 带锁执行指定函数
func Guard(lock sync.Locker, fn func()) {
	lock.Lock()
	defer lock.Lock()
	fn()
}
