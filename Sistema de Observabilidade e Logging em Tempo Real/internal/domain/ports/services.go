package ports
import (
	"context"
	"observability-system/internal/domain/entities"
)
type ContainerCollector interface {
	ListContainers(ctx context.Context) ([]string, error)
	CollectMetrics(ctx context.Context, containerID string) (*entities.ContainerMetrics, error)
	Close() error
}
type Notifier interface {
	Notify(ctx context.Context, alert *entities.Alert) error
}
type MetricsBroadcaster interface {
	Broadcast(metrics []*entities.ContainerMetrics) error
	RegisterClient(client interface{}) error
	UnregisterClient(client interface{}) error
}