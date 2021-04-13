package syncx

import (
	"math"
	"sync/atomic"
)

// AtomicFloat64 float64类型 原子类
type AtomicFloat64 uint64

// NewAtomicFloat64 返回一个 AtomicFloat64
func NewAtomicFloat64() *AtomicFloat64 {
	return new(AtomicFloat64)
}

func ForAtomicFloat64(v float64) *AtomicFloat64 {
	f := NewAtomicFloat64()
	f.Set(v)
	return f
}

// Set 设置当前值到val。
func (f *AtomicFloat64) Set(val float64) {
	atomic.StoreUint64((*uint64)(f), math.Float64bits(val))
}

// Load 加载当前值。
func (f *AtomicFloat64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(f)))
}

// CompareAndSwap 比较当前值和老值，如果相等则将当前值赋给val。
func (f *AtomicFloat64) CompareAndSwap(old, val float64) bool {
	return atomic.CompareAndSwapUint64((*uint64)(f), math.Float64bits(old), math.Float64bits(val))
}

// Add 添加val到当前值
func (f *AtomicFloat64) Add(val float64) float64 {
	for {
		old := f.Load()
		nv := old + val
		if f.CompareAndSwap(old, nv) {
			return nv
		}
	}
}
