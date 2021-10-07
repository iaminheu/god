package load

// nopShedder 无操作的负载泄流器
type (
	nopShedder struct{}
	nopPromise struct{}
)

func newNopShedder() nopShedder {
	return nopShedder{}
}

func (s nopShedder) Allow() (Promise, error) {
	return nopPromise{}, nil
}

func (p nopPromise) Pass() {}

func (p nopPromise) Fail() {}
