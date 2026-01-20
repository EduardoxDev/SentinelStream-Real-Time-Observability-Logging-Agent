package adapters
import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"observability-system/internal/domain/entities"
)
type ProcessCollectorAdapter struct {
	processes []string
}
func NewProcessCollectorAdapter() (*ProcessCollectorAdapter, error) {
	processes := []string{
		"observability-agent",
		"observability-server",
		"redis-server",
		"influxdb",
		"chrome",
		"vscode",
		"powershell",
	}
	return &ProcessCollectorAdapter{
		processes: processes,
	}, nil
}
func (d *ProcessCollectorAdapter) ListContainers(ctx context.Context) ([]string, error) {
	return d.processes, nil
}
func (d *ProcessCollectorAdapter) CollectMetrics(ctx context.Context, processName string) (*entities.ContainerMetrics, error) {
	baseLoad := rand.Float64() * 30
	spike := 0.0
	if rand.Float64() > 0.8 {
		spike = rand.Float64() * 40
	}
	cpuPercent := baseLoad + spike
	memoryPercent := 20 + rand.Float64()*50
	if rand.Float64() > 0.95 {
		cpuPercent = 92 + rand.Float64()*5
	}
	memoryUsage := uint64(memoryPercent * 10 * 1024 * 1024)
	memoryLimit := uint64(1024 * 1024 * 1024)
	networkRx := uint64(rand.Int63n(1024 * 1024 * 10))
	networkTx := uint64(rand.Int63n(1024 * 1024 * 5))
	return &entities.ContainerMetrics{
		ContainerID:   fmt.Sprintf("process-%s", processName),
		ContainerName: processName,
		CPUPercent:    cpuPercent,
		MemoryUsage:   memoryUsage,
		MemoryLimit:   memoryLimit,
		MemoryPercent: memoryPercent,
		NetworkRx:     networkRx,
		NetworkTx:     networkTx,
		Timestamp:     time.Now(),
	}, nil
}
func (d *ProcessCollectorAdapter) Close() error {
	return nil
}