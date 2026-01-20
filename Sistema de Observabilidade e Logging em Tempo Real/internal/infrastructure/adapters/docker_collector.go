package adapters
import (
	"context"
	"encoding/json"
	"io"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"observability-system/internal/domain/entities"
	"time"
)
type DockerCollectorAdapter struct {
	client *client.Client
}
func NewDockerCollectorAdapter() (*DockerCollectorAdapter, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerCollectorAdapter{client: cli}, nil
}
func (d *DockerCollectorAdapter) ListContainers(ctx context.Context) ([]string, error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(containers))
	for i, container := range containers {
		ids[i] = container.ID
	}
	return ids, nil
}
func (d *DockerCollectorAdapter) CollectMetrics(ctx context.Context, containerID string) (*entities.ContainerMetrics, error) {
	stats, err := d.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()
	var v types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	containerInfo, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	cpuPercent := calculateCPUPercent(&v)
	memPercent := float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
	var networkRx, networkTx uint64
	for _, net := range v.Networks {
		networkRx += net.RxBytes
		networkTx += net.TxBytes
	}
	return &entities.ContainerMetrics{
		ContainerID:   containerID,
		ContainerName: containerInfo.Name,
		CPUPercent:    cpuPercent,
		MemoryUsage:   v.MemoryStats.Usage,
		MemoryLimit:   v.MemoryStats.Limit,
		MemoryPercent: memPercent,
		NetworkRx:     networkRx,
		NetworkTx:     networkTx,
		Timestamp:     time.Now(),
	}, nil
}
func calculateCPUPercent(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}
func (d *DockerCollectorAdapter) Close() error {
	return d.client.Close()
}