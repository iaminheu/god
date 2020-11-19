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

type Span struct {
	ctx           spanContext
	serviceName   string
	operationName string
	startTime     time.Time
	flag          string
	children      int
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
		flag:          serverFlag,
	}
}

func StartClientSpan(ctx context.Context, serviceName, operationName string) (context.Context, Trace) {
	if span, ok := ctx.Value(TracingKey).(*Span); ok {
		return span.Fork(ctx, serviceName, operationName)
	}

	return ctx, emptyNopSpan
}

func StartServerSpan(ctx context.Context, payload Payload, serviceName, operationName string) (context.Context, Trace) {
	span := newServerSpan(payload, serviceName, operationName)
	return context.WithValue(ctx, TracingKey, span), span
}
