package contextx

import (
	"context"
	"git.zc0901.com/go/god/lib/mapping"
)

const contextTagKey = "ctx"

var unmarshaler = mapping.NewUnmarshaler(contextTagKey)

type contextValuer struct {
	context.Context
}

func (cv contextValuer) Value(key string) (interface{}, bool) {
	v := cv.Context.Value(key)
	return v, v != nil
}

func For(ctx context.Context, v interface{}) error {
	return unmarshaler.UnmarshalValuer(contextValuer{
		Context: ctx,
	}, v)
}
