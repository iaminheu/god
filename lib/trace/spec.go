package trace

import "context"

const (
	traceIdKey = "X-Trace-ID"
	spanIdKey  = "X-Span-ID"
)

type (
	// 链路的操作上下文：TraceID、SpanID或其他想要传递的内容
	SpanContext interface {
		TraceId() string
		SpanId() string
		Visit(fn func(key, value string) bool) // 自定义操作TraceId，SpanId
	}

	// 链路接口
	Trace interface {
		SpanContext

		Finish()
		Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
		Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
	}

	contextKey string
)

var TracingKey = contextKey("X-Trace")

func (c contextKey) String() string {
	return "trace/spec.go 上下文的键 " + string(c)
}
