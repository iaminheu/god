package contextx

import (
	"context"
	"time"
)

// ShrinkDeadline 根据上下文截止时间按需缩短（通过上下文控制级联联超时？）
func ShrinkDeadline(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if deadline, ok := ctx.Deadline(); ok {
		leftTime := time.Until(deadline)
		if leftTime < timeout {
			timeout = leftTime
		}
	}

	return context.WithDeadline(ctx, time.Now().Add(timeout))
}
