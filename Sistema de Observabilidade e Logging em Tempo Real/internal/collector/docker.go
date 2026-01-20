package collector
import (
	"context"
	"encoding/json"
	"io"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)
type DockerCollector struct {
	client *client.Client
}
type ContainerStats struct {
	ContainerID   string
	ContainerName string
	CPUPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	MemoryPercent float64
	NetworkRx     uint64
	NetworkTx     uint64
	Timestamp     time.Time
}
func NewDockerCollector() (*DockerCollector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerCollector{client: cli}, nil
}
func (dc *DockerCollector) GetRunningContainers(ctx context.Context) ([]types.Container, error) {
	return dc.client.ContainerList(ctx, types.ContainerListOptions{})
}
func (dc *DockerCollector) CollectStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	stats, err := dc.client.ContainerStats(ctx, containerID, false)
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
	containerInfo, err := dc.client.ContainerInspect(ctx, containerID)
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
	return &ContainerStats{
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
func (dc *DockerCollector) Close() error {
	return dc.client.Close()
}