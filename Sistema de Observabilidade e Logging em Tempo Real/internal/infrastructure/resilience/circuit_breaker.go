package resilience
import (
	"context"
	"errors"
	"sync"
	"time"
)
type State int
const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)
var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests")
)
type CircuitBreaker struct {
	maxFailures  uint32
	timeout      time.Duration
	state        State
	failures     uint32
	lastFailTime time.Time
	mu           sync.RWMutex
}
func NewCircuitBreaker(maxFailures uint32, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       StateClosed,
	}
}
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}
	err := fn()
	cb.afterRequest(err)
	return err
}
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = StateHalfOpen
			return nil
		}
		return ErrCircuitOpen
	case StateHalfOpen:
		return ErrTooManyRequests
	default:
		return nil
	}
}
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	} else {
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
		}
		cb.failures = 0
	}
}
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}