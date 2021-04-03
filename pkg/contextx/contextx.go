package contextx

import (
	"context"
	"time"
)

func ShrinkDeadline(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if deadline, ok := ctx.Deadline(); ok {
		if leftTime := time.Until(deadline); leftTime < timeout {
			timeout = leftTime
		}
	}
	return context.WithDeadline(ctx, time.Now().Add(timeout))
}
