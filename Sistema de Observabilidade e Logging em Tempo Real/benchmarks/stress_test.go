package benchmarks
import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	"observability-system/internal/domain/entities"
	"observability-system/internal/infrastructure/adapters"
)
func TestStressInfluxDB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}
	repo := adapters.NewInfluxDBRepository(
		"http:
		"my-super-secret-token",
		"observability",
		"metrics",
	)
	defer repo.Close()
	ctx := context.Background()
	concurrency := 100
	iterations := 1000
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				metrics := &entities.ContainerMetrics{
					ContainerID:   fmt.Sprintf("container-%d", workerID),
					ContainerName: fmt.Sprintf("test-%d", workerID),
					CPUPercent:    float64(j % 100),
					MemoryUsage:   uint64(j * 1024),
					MemoryLimit:   1024 * 1024 * 500,
					MemoryPercent: float64(j % 100),
					NetworkRx:     uint64(j * 100),
					NetworkTx:     uint64(j * 50),
					Timestamp:     time.Now(),
				}
				repo.Save(ctx, metrics)
			}
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	totalOps := concurrency * iterations
	opsPerSec := float64(totalOps) / duration.Seconds()
	t.Logf("Stress Test Results:")
	t.Logf("  Total Operations: %d", totalOps)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Operations/sec: %.2f", opsPerSec)
}
func TestStressRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}
	repo := adapters.NewRedisAlertRepository("localhost:6379")
	defer repo.Close()
	ctx := context.Background()
	concurrency := 50
	iterations := 500
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				containerID := fmt.Sprintf("container-%d", workerID)
				repo.IsInCooldown(ctx, containerID, entities.AlertTypeCPU)
				if j%10 == 0 {
					repo.SetCooldown(ctx, containerID, entities.AlertTypeCPU, 1*time.Minute)
				}
			}
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	totalOps := concurrency * iterations
	opsPerSec := float64(totalOps) / duration.Seconds()
	t.Logf("Redis Stress Test Results:")
	t.Logf("  Total Operations: %d", totalOps)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Operations/sec: %.2f", opsPerSec)
}