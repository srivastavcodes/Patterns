package stability_patterns

import (
	"context"
	"sync"
	"time"
)

func DebounceFirstContext(circuit Circuit, dur time.Duration) Circuit {
	var threshold time.Time
	var mu sync.Mutex

	var lastCtx context.Context
	var lastCancel context.CancelFunc

	return func(ctx context.Context) (string, error) {
		mu.Lock()

		if time.Now().Before(threshold) {
			lastCancel()
		}
		lastCtx, lastCancel = context.WithCancel(ctx)
		threshold = time.Now().Add(dur)

		mu.Unlock()
		result, err := circuit(lastCtx)
		return result, err
	}
}
