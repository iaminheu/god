package collection

import "sync"

// Ring 是一个定长的环形切片。
type Ring struct {
	elements []interface{}
	index    int
	lock     sync.Mutex
}

func NewRing(size int) *Ring {
	if size < 1 {
		panic("size 应大于 0")
	}

	return &Ring{
		elements: make([]interface{}, size),
	}
}

// Add 加值到环。
func (r *Ring) Add(v interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.elements[r.index%len(r.elements)] = v
	r.index++
}

// Take 从环取值。
func (r *Ring) Take() []interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()

	var size int
	var start int
	if r.index > len(r.elements) {
		size = len(r.elements)
		start = r.index % len(r.elements)
	} else {
		size = r.index
	}

	elems := make([]interface{}, size)
	for i := 0; i < size; i++ {
		elems[i] = r.elements[(start+i)%len(r.elements)]
	}

	return elems
}
