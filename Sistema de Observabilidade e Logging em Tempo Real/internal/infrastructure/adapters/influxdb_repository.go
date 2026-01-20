package adapters
import (
	"context"
	"fmt"
	"log"
	"time"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"observability-system/internal/domain/entities"
	"observability-system/internal/infrastructure/resilience"
)
type InfluxDBRepository struct {
	client         influxdb2.Client
	writeAPI       api.WriteAPI
	queryAPI       api.QueryAPI
	org            string
	bucket         string
	circuitBreaker *resilience.CircuitBreaker
	retryPolicy    *resilience.RetryPolicy
}
func NewInfluxDBRepository(url, token, org, bucket string) *InfluxDBRepository {
	client := influxdb2.NewClient(url, token)
	return &InfluxDBRepository{
		client:         client,
		writeAPI:       client.WriteAPI(org, bucket),
		queryAPI:       client.QueryAPI(org),
		org:            org,
		bucket:         bucket,
		circuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		retryPolicy:    resilience.NewRetryPolicy(3, 1*time.Second, 2.0),
	}
}
func (r *InfluxDBRepository) Save(ctx context.Context, metrics *entities.ContainerMetrics) error {
	return r.circuitBreaker.Execute(ctx, func() error {
		return r.retryPolicy.Execute(ctx, func() error {
			p := influxdb2.NewPoint(
				"container_metrics",
				map[string]string{
					"container_id":   metrics.ContainerID,
					"container_name": metrics.ContainerName,
				},
				map[string]interface{}{
					"cpu_percent":    metrics.CPUPercent,
					"memory_usage":   metrics.MemoryUsage,
					"memory_limit":   metrics.MemoryLimit,
					"memory_percent": metrics.MemoryPercent,
					"network_rx":     metrics.NetworkRx,
					"network_tx":     metrics.NetworkTx,
				},
				metrics.Timestamp,
			)
			r.writeAPI.WritePoint(p)
			r.writeAPI.Flush()
			select {
			case err := <-r.writeAPI.Errors():
				return err
			default:
				return nil
			}
		})
	})
}
func (r *InfluxDBRepository) FindByContainerID(ctx context.Context, containerID string, duration time.Duration) ([]*entities.ContainerMetrics, error) {
	var result []*entities.ContainerMetrics
	err := r.circuitBreaker.Execute(ctx, func() error {
		query := fmt.Sprintf(`
			from(bucket: "%s")
			|> range(start: -%s)
			|> filter(fn: (r) => r["_measurement"] == "container_metrics")
			|> filter(fn: (r) => r["container_id"] == "%s")
		`, r.bucket, duration.String(), containerID)
		queryResult, err := r.queryAPI.Query(ctx, query)
		if err != nil {
			return err
		}
		for queryResult.Next() {
			log.Printf("Query result: %v", queryResult.Record().Values())
		}
		return queryResult.Err()
	})
	return result, err
}
func (r *InfluxDBRepository) FindAll(ctx context.Context, duration time.Duration) ([]*entities.ContainerMetrics, error) {
	return nil, nil
}
func (r *InfluxDBRepository) Close() error {
	r.writeAPI.Flush()
	r.client.Close()
	return nil
}