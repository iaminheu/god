package trace

import "context"

const (
	traceIdKey = "X-Trace-ID"
	spanIdKey  = "X-Span-ID"
)

type (
	SpanContext interface {
		TraceId() string
		SpanId() string
		Visit(fn func(key, value string) bool)
	}

	Trace interface {
		SpanContext

		Finish()
		Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
		Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
	}

	spanContext struct {
		traceId string
		spanId  string
	}

	contextKey string
)

func (sc spanContext) TraceId() string {
	return sc.traceId
}

func (sc spanContext) SpanId() string {
	return sc.spanId
}

func (sc spanContext) Visit(fn func(key, val string) bool) {
	fn(traceIdKey, sc.traceId)
	fn(spanIdKey, sc.spanId)
}

var TracingKey = contextKey("X-Trace")

func (c contextKey) String() string {
	return "trace/spec.go 上下文的键 " + string(c)
}
