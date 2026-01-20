package adapters
import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
	"observability-system/internal/domain/entities"
	"observability-system/internal/infrastructure/resilience"
)
type RedisAlertRepository struct {
	client         *redis.Client
	circuitBreaker *resilience.CircuitBreaker
	retryPolicy    *resilience.RetryPolicy
}
func NewRedisAlertRepository(addr string) *RedisAlertRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisAlertRepository{
		client:         client,
		circuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		retryPolicy:    resilience.NewRetryPolicy(3, 500*time.Millisecond, 2.0),
	}
}
func (r *RedisAlertRepository) Save(ctx context.Context, alert *entities.Alert) error {
	return r.circuitBreaker.Execute(ctx, func() error {
		return r.retryPolicy.Execute(ctx, func() error {
			key := fmt.Sprintf("alert:%s:%s:%d", alert.ContainerID, alert.Type, alert.Timestamp.Unix())
			return r.client.Set(ctx, key, alert.Message, 24*time.Hour).Err()
		})
	})
}
func (r *RedisAlertRepository) IsInCooldown(ctx context.Context, containerID string, alertType entities.AlertType) (bool, error) {
	var result bool
	err := r.circuitBreaker.Execute(ctx, func() error {
		key := fmt.Sprintf("cooldown:%s:%s", containerID, alertType)
		exists, err := r.client.Exists(ctx, key).Result()
		result = exists > 0
		return err
	})
	return result, err
}
func (r *RedisAlertRepository) SetCooldown(ctx context.Context, containerID string, alertType entities.AlertType, duration time.Duration) error {
	return r.circuitBreaker.Execute(ctx, func() error {
		return r.retryPolicy.Execute(ctx, func() error {
			key := fmt.Sprintf("cooldown:%s:%s", containerID, alertType)
			return r.client.Set(ctx, key, "1", duration).Err()
		})
	})
}
func (r *RedisAlertRepository) Close() error {
	return r.client.Close()
}