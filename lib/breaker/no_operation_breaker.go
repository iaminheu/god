package breaker

const noOpBreakerName = "nopBreaker"

// 无操作熔断器
type nopBreaker struct{}
type nopPromise struct{}

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

func (b nopBreaker) DoWithFallback(req Request, fallback Fallback) error {
	return req()
}

func (b nopBreaker) DoWithAcceptable(req Request, acceptable Acceptable) error {
	return req()
}

func (b nopBreaker) DoWithFallbackAcceptable(req Request, fallback Fallback, acceptable Acceptable) error {
	return req()
}

func (p nopPromise) Accept() {}

func (p nopPromise) Reject(reason string) {}
