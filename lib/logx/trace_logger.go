package logx

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/timex"
	"git.zc0901.com/go/god/lib/trace"
	"io"
	"time"
)

type traceLogger struct {
	logEntry
	Trace string `json:"trace,omitempty"`
	Span  string `json:"span,omitempty"`
	ctx   context.Context
}

func WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
	}
}

func (l *traceLogger) Info(v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLogger, infoLevel, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLogger, infoLevel, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Error(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLogger, errorLevel, formatWithCaller(fmt.Sprint(v...), durationCallerDepth))
	}
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLogger, errorLevel, formatWithCaller(fmt.Sprintf(format, v...), durationCallerDepth))
	}
}

func (l *traceLogger) Slow(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLogger, slowLevel, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLogger, slowLevel, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *traceLogger) write(writer io.WriteCloser, level string, content string) {
	l.Timestamp = getTimestamp()
	l.Level = level
	l.Content = content
	l.Trace = traceIdFromContext(l.ctx)
	l.Span = spanIdFromContext(l.ctx)
	outputJson(writer, l)
}

func traceIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(trace.TracingKey).(trace.Trace)
	if !ok {
		return ""
	}

	return t.TraceId()
}

func spanIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(trace.TracingKey).(trace.Trace)
	if !ok {
		return ""
	}

	return t.SpanId()
}
