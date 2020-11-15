package syncx

import (
	"runtime"
	"sync/atomic"
)

// SwitchLock 开关锁
type SwitchLock struct {
	lock uint32
}

func (sl *SwitchLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}

func (sl *SwitchLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

func (sl *SwitchLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}
