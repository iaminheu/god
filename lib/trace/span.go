package trace

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/stringx"
	"git.zc0901.com/go/god/lib/timex"
	"strconv"
	"strings"
	"time"
)

const (
	initSpanId  = "0"
	clientFlag  = "client"
	serverFlag  = "server"
	spanSepRune = '.'
)

var spanSep = string([]byte{spanSepRune})

// 链路中的一个操作（api调用、数据库操作等），存储时间、服务名称等信息
// 是分布式追踪的最小单元。
// 一个 Trace 由多个 Span 组成。
type Span struct {
	ctx           spanContext // 传递的上下文
	serviceName   string      // 服务名
	operationName string      // 操作
	startTime     time.Time   // 开始时间戳
	flag          string      // 标记 trace 是 server 还是 client
	children      int         // 本 span fork 的子节点数
}

func (s *Span) TraceId() string {
	return s.ctx.TraceId()
}

func (s *Span) SpanId() string {
	return s.ctx.SpanId()
}

func (s *Span) Visit(fn func(key string, value string) bool) {
	s.ctx.Visit(fn)
}

func (s *Span) Finish() {

}

func (s *Span) Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.forkSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          clientFlag,
	}
	return context.WithValue(ctx, TracingKey, span), span
}

func (s *Span) Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.followSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          s.flag,
	}
	return context.WithValue(ctx, TracingKey, span), span
}

func (s *Span) forkSpanId() string {
	s.children++
	return fmt.Sprintf("%s.%d", s.ctx.spanId, s.children)
}

func (s *Span) followSpanId() string {
	fields := strings.FieldsFunc(s.ctx.spanId, func(r rune) bool {
		return r == spanSepRune
	})
	if len(fields) == 0 {
		return s.ctx.spanId
	}

	last := fields[len(fields)-1]
	val, err := strconv.Atoi(last)
	if err != nil {
		return s.ctx.spanId
	}

	last = strconv.Itoa(val + 1)
	fields[len(fields)-1] = last

	return strings.Join(fields, spanSep)
}

func newServerSpan(payload Payload, serviceName, operationName string) Trace {
	// 从上述的 payload「也就是header」获取traceId，spanId。
	traceId := stringx.TakeWithPriority(func() string {
		if payload != nil {
			return payload.Get(traceIdKey)
		}

		return ""
	}, stringx.RandId)
	spanId := stringx.TakeWithPriority(func() string {
		if payload != nil {
			return payload.Get(spanIdKey)
		}
		return ""
	}, func() string {
		return initSpanId
	})

	return &Span{
		ctx:           spanContext{traceId: traceId, spanId: spanId},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          serverFlag, // 标记为server
	}
}

// 开启客户端链路操作
func StartClientSpan(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	// **1** 获取上游（api 或其他 rpc 调用端）带下来的 span 上下文信息
	if span, ok := ctx.Value(TracingKey).(*Span); ok {
		// **2** 在此客户端中分裂出子Span，（从获取的 span 中创建 子span，「继承父span的traceId」）
		return span.Fork(ctx, serviceName, operationName)
	}

	return ctx, emptyNopSpan
}

// 开启服务端链路操作
func StartServerSpan(ctx context.Context, payload Payload, serviceName, operationName string) (context.Context, Trace) {
	span := newServerSpan(payload, serviceName, operationName)
	// **4** - 看header中是否设置
	return context.WithValue(ctx, TracingKey, span), span
}
