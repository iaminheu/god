package trace

import "context"

var emptyNopSpan = nopSpan{}

type nopSpan struct{}

func (s nopSpan) TraceId() string {
	return ""
}

func (s nopSpan) SpanId() string {
	return ""
}

func (s nopSpan) Visit(fn func(key string, value string) bool) {}

func (s nopSpan) Finish() {}

func (s nopSpan) Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	return ctx, emptyNopSpan
}

func (s nopSpan) Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	return ctx, emptyNopSpan
}
