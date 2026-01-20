package ports
import (
	"context"
	"time"
	"observability-system/internal/domain/entities"
)
type MetricsRepository interface {
	Save(ctx context.Context, metrics *entities.ContainerMetrics) error
	FindByContainerID(ctx context.Context, containerID string, duration time.Duration) ([]*entities.ContainerMetrics, error)
	FindAll(ctx context.Context, duration time.Duration) ([]*entities.ContainerMetrics, error)
	Close() error
}
type AlertRepository interface {
	Save(ctx context.Context, alert *entities.Alert) error
	IsInCooldown(ctx context.Context, containerID string, alertType entities.AlertType) (bool, error)
	SetCooldown(ctx context.Context, containerID string, alertType entities.AlertType, duration time.Duration) error
	Close() error
}