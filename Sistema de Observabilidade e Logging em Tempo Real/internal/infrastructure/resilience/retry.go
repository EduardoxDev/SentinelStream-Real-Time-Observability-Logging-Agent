package resilience
import (
	"context"
	"time"
)
type RetryPolicy struct {
	maxAttempts int
	delay       time.Duration
	backoff     float64
}
func NewRetryPolicy(maxAttempts int, delay time.Duration, backoff float64) *RetryPolicy {
	return &RetryPolicy{
		maxAttempts: maxAttempts,
		delay:       delay,
		backoff:     backoff,
	}
}
func (rp *RetryPolicy) Execute(ctx context.Context, fn func() error) error {
	var err error
	currentDelay := rp.delay
	for attempt := 0; attempt < rp.maxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt < rp.maxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(currentDelay):
				currentDelay = time.Duration(float64(currentDelay) * rp.backoff)
			}
		}
	}
	return err
}