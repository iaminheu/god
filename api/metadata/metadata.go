package metadata

import (
	"context"
)

var (
	RemoteIP   = "remote-ip"
	RemotePort = "remote-port"
)

// MD is a mapping from metadata keys to values. Users should use the following
// two convenience functions New and Pairs to generate MD.
type MD map[string]string

// Len returns the number of items in md.
func (md MD) Len() int {
	return len(md)
}

// NewIncomingContext creates a new context with incoming md attached.
func NewContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, struct{}{}, md)
}
