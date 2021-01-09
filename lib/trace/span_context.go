package trace

type spanContext struct {
	traceId string // 表示tracer的全局唯一ID
	spanId  string // 表示单个trace中某一个span的唯一ID，在trace中唯一
}

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
