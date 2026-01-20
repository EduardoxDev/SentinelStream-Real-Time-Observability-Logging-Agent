package main
import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"observability-system/internal/application/usecases"
	"observability-system/internal/domain/entities"
	"observability-system/internal/infrastructure/adapters"
)
func main() {
	log.Println("🚀 Starting Observability Agent (Clean Architecture)...")
	processCollector, err := adapters.NewProcessCollectorAdapter()
	if err != nil {
		log.Fatalf("Failed to create Process collector: %v", err)
	}
	defer processCollector.Close()

	metricsRepo := adapters.NewInfluxDBRepository(
		getEnv("INFLUXDB_URL", "http://localhost:8086"),
		getEnv("INFLUXDB_TOKEN", "my-super-secret-token"),
		getEnv("INFLUXDB_ORG", "observability"),
		getEnv("INFLUXDB_BUCKET", "metrics"),
	)
	defer metricsRepo.Close()

	alertRepo := adapters.NewRedisAlertRepository(getEnv("REDIS_ADDR", "localhost:6379"))
	defer alertRepo.Close()

	notifier := adapters.NewConsoleNotifier()

	collectMetricsUC := usecases.NewCollectMetricsUseCase(processCollector, metricsRepo)
	checkAlertsUC := usecases.NewCheckAlertsUseCase(alertRepo, notifier, 90.0, 85.0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go collectMetrics(ctx, collectMetricsUC, checkAlertsUC, alertRepo)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("🛑 Shutting down agent...")
}
func collectMetrics(ctx context.Context, collectUC *usecases.CollectMetricsUseCase, alertUC *usecases.CheckAlertsUseCase, alertRepo *adapters.RedisAlertRepository) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			allMetrics, err := collectUC.Execute(ctx)
			if err != nil {
				log.Printf("Error collecting metrics: %v", err)
				continue
			}
			for _, metrics := range allMetrics {
				if err := alertUC.Execute(ctx, metrics); err != nil {
					log.Printf("Error checking alerts: %v", err)
				}
				if !metrics.IsHealthy(90.0, 85.0) {
					violations := metrics.ExceedsThreshold(90.0, 85.0)
					for _, v := range violations {
						var alertType entities.AlertType
						if v == "CPU" {
							alertType = entities.AlertTypeCPU
						} else {
							alertType = entities.AlertTypeMemory
						}
						alertRepo.SetCooldown(ctx, metrics.ContainerID, alertType, 5*time.Minute)
					}
				}
				log.Printf("📊 %s - CPU: %.2f%% | Memory: %.2f%% | Net RX: %d TX: %d",
					metrics.ContainerName,
					metrics.CPUPercent,
					metrics.MemoryPercent,
					metrics.NetworkRx,
					metrics.NetworkTx,
				)
			}
		}
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}