package trace

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNopSpan_Fork(t *testing.T) {
	ctx, span := emptyNopSpan.Fork(context.Background(), "", "")
	assert.Equal(t, emptyNopSpan, span)
	assert.Equal(t, context.Background(), ctx)
}

func TestNopSpan_Follow(t *testing.T) {
	ctx, span := emptyNopSpan.Follow(context.Background(), "", "")
	assert.Equal(t, emptyNopSpan, span)
	assert.Equal(t, context.Background(), ctx)
}

func TestNopSpan(t *testing.T) {
	emptyNopSpan.Visit(func(key string, value string) bool {
		assert.Fail(t, "不该到此")
		return true
	})

	ctx, span := emptyNopSpan.Follow(context.Background(), "", "")
	assert.Equal(t, context.Background(), ctx)
	assert.Equal(t, "", span.TraceId())
	assert.Equal(t, "", span.SpanId())
}
