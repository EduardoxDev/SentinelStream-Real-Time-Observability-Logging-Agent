package usecases
import (
	"context"
	"log"
	"observability-system/internal/domain/entities"
	"observability-system/internal/domain/ports"
)
type CollectMetricsUseCase struct {
	collector  ports.ContainerCollector
	repository ports.MetricsRepository
}
func NewCollectMetricsUseCase(collector ports.ContainerCollector, repository ports.MetricsRepository) *CollectMetricsUseCase {
	return &CollectMetricsUseCase{
		collector:  collector,
		repository: repository,
	}
}
func (uc *CollectMetricsUseCase) Execute(ctx context.Context) ([]*entities.ContainerMetrics, error) {
	containers, err := uc.collector.ListContainers(ctx)
	if err != nil {
		return nil, err
	}
	var allMetrics []*entities.ContainerMetrics
	for _, containerID := range containers {
		metrics, err := uc.collector.CollectMetrics(ctx, containerID)
		if err != nil {
			log.Printf("Failed to collect metrics for %s: %v", containerID, err)
			continue
		}
		if metrics == nil {
			continue
		}
		if err := uc.repository.Save(ctx, metrics); err != nil {
			log.Printf("Failed to save metrics for %s: %v", containerID, err)
		}
		allMetrics = append(allMetrics, metrics)
	}
	return allMetrics, nil
}