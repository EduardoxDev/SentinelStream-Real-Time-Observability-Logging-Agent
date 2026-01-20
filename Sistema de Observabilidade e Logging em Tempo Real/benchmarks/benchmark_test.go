package benchmarks
import (
	"context"
	"testing"
	"time"
	"observability-system/internal/domain/entities"
	"observability-system/internal/infrastructure/adapters"
)
func BenchmarkMetricsCollection(b *testing.B) {
	collector, err := adapters.NewDockerCollectorAdapter()
	if err != nil {
		b.Fatalf("Failed to create collector: %v", err)
	}
	defer collector.Close()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containers, _ := collector.ListContainers(ctx)
		if len(containers) > 0 {
			collector.CollectMetrics(ctx, containers[0])
		}
	}
}
func BenchmarkInfluxDBWrite(b *testing.B) {
	repo := adapters.NewInfluxDBRepository(
		"http:
		"my-super-secret-token",
		"observability",
		"metrics",
	)
	defer repo.Close()
	ctx := context.Background()
	metrics := &entities.ContainerMetrics{
		ContainerID:   "test-container",
		ContainerName: "test",
		CPUPercent:    50.0,
		MemoryUsage:   1024 * 1024 * 100,
		MemoryLimit:   1024 * 1024 * 500,
		MemoryPercent: 20.0,
		NetworkRx:     1000,
		NetworkTx:     2000,
		Timestamp:     time.Now(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Save(ctx, metrics)
	}
}
func BenchmarkRedisAlertCheck(b *testing.B) {
	repo := adapters.NewRedisAlertRepository("localhost:6379")
	defer repo.Close()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.IsInCooldown(ctx, "test-container", entities.AlertTypeCPU)
	}
}
func BenchmarkCircuitBreaker(b *testing.B) {
	repo := adapters.NewInfluxDBRepository(
		"http:
		"my-super-secret-token",
		"observability",
		"metrics",
	)
	defer repo.Close()
	ctx := context.Background()
	metrics := &entities.ContainerMetrics{
		ContainerID:   "test-container",
		ContainerName: "test",
		CPUPercent:    50.0,
		MemoryUsage:   1024 * 1024 * 100,
		MemoryLimit:   1024 * 1024 * 500,
		MemoryPercent: 20.0,
		Timestamp:     time.Now(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Save(ctx, metrics)
	}
}
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = &entities.ContainerMetrics{
			ContainerID:   "test-container-id-12345",
			ContainerName: "/test-container-name",
			CPUPercent:    75.5,
			MemoryUsage:   1024 * 1024 * 256,
			MemoryLimit:   1024 * 1024 * 512,
			MemoryPercent: 50.0,
			NetworkRx:     1024 * 1024 * 10,
			NetworkTx:     1024 * 1024 * 5,
			Timestamp:     time.Now(),
		}
	}
}