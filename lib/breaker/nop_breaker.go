package breaker

const noOpBreakerName = "nopBreaker"

// 无操作断路器
type (
	nopBreaker struct{}
	nopPromise struct{}
)

func newNopBreaker() Breaker {
	return nopBreaker{}
}

func (b nopBreaker) Name() string {
	return noOpBreakerName
}

func (b nopBreaker) Allow() (Promise, error) {
	return nopPromise{}, nil
}

func (b nopBreaker) Do(req Request) error {
	return req()
}

func (b nopBreaker) DoWithFallback(req Request, _ Fallback) error {
	return req()
}

func (b nopBreaker) DoWithAcceptable(req Request, _ Acceptable) error {
	return req()
}

func (b nopBreaker) DoWithFallbackAcceptable(req Request, _ Fallback, _ Acceptable) error {
	return req()
}

func (p nopPromise) Accept() {}

func (p nopPromise) Reject(string) {}
