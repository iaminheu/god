package syncx

import (
	"errors"
	"git.zc0901.com/go/god/lib/lang"
)

var ErrLimitReturn = errors.New("请求限制")

// Limit 控制并发请求。
type Limit struct {
	pool chan lang.PlaceholderType
}

// NewLimit 新建并发控制
func NewLimit(n int) Limit {
	return Limit{
		pool: make(chan lang.PlaceholderType, n),
	}
}

// Borrow 在阻塞模式下从 Limit 借用一个元素。
func (l Limit) Borrow() {
	l.pool <- lang.Placeholder
}

// Return 返回借用资源，当返回的比借用的多则返回错误。
func (l Limit) Return() error {
	select {
	case <-l.pool:
		return nil
	default:
		return ErrLimitReturn
	}
}

// TryBorrow 尝试从 Limit 借用一个元素（非阻塞模式下）。
// 如果成功则返回 true，否则返回 false。
func (l Limit) TryBorrow() bool {
	select {
	case l.pool <- lang.Placeholder:
		return true
	default:
		return false
	}
}
