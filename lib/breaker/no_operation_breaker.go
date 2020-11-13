package breaker

const noOpBreakerName = "noOpBreaker"

type noOpBreaker struct{}
type noOpPromise struct{}

func newNoOpBreaker() Breaker {
	return noOpBreaker{}
}

func (b noOpBreaker) Name() string {
	return noOpBreakerName
}

func (b noOpBreaker) Allow() (Promise, error) {
	return noOpPromise{}, nil
}

func (b noOpBreaker) Do(req Request) error {
	return req()
}

func (b noOpBreaker) DoWithFailback(req Request, fallback Fallback) error {
	panic("implement me")
}

func (b noOpBreaker) DoWithAcceptable(req Request, acceptable Acceptable) error {
	panic("implement me")
}

func (b noOpBreaker) DoWithFailbackAcceptable(req Request, fallback Fallback, acceptable Acceptable) error {
	panic("implement me")
}

func (p noOpPromise) Accept() {
	panic("implement me")
}

func (p noOpPromise) Reject(reason string) {
	panic("implement me")
}
