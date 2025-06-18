package stability_patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Circuit is the type you will receive
type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, threshold int) Circuit {
	last := time.Now()
	var failures int
	var rwm sync.RWMutex

	return func(ctx context.Context) (string, error) {
		rwm.RLock()
		tracker := failures - threshold

		if tracker >= 0 {
			shouldRetryAt := last.Add((2 << tracker) * time.Second)

			if !time.Now().After(shouldRetryAt) {
				rwm.RUnlock()
				return "", errors.New("service unavailable")
			}
		}
		rwm.RUnlock()
		response, err := circuit(ctx)

		rwm.Lock()
		defer rwm.Unlock()

		last = time.Now()
		if err != nil {
			failures++
			return response, err
		}
		failures = 0
		return response, nil
	}
}
