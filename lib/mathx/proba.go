package mathx

import (
	"math/rand"
	"sync"
	"time"
)

// 维护可能性的概率
type Proba struct {
	r    *rand.Rand
	lock sync.Mutex
}

// 新建可能性概率
func NewProba() *Proba {
	return &Proba{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 返回 v 值是否大于随机值
func (p *Proba) TrueOnProba(v float64) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return v > p.r.Float64()
}
