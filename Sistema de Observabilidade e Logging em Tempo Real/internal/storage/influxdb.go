package storage
import (
	"context"
	"fmt"
	"time"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"observability-system/internal/collector"
)
type InfluxDBStorage struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	org      string
	bucket   string
}
func NewInfluxDBStorage(url, token, org, bucket string) *InfluxDBStorage {
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)
	queryAPI := client.QueryAPI(org)
	return &InfluxDBStorage{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		org:      org,
		bucket:   bucket,
	}
}
func (s *InfluxDBStorage) WriteMetrics(stats *collector.ContainerStats) error {
	p := influxdb2.NewPoint(
		"container_metrics",
		map[string]string{
			"container_id":   stats.ContainerID,
			"container_name": stats.ContainerName,
		},
		map[string]interface{}{
			"cpu_percent":    stats.CPUPercent,
			"memory_usage":   stats.MemoryUsage,
			"memory_limit":   stats.MemoryLimit,
			"memory_percent": stats.MemoryPercent,
			"network_rx":     stats.NetworkRx,
			"network_tx":     stats.NetworkTx,
		},
		stats.Timestamp,
	)
	s.writeAPI.WritePoint(p)
	return nil
}
func (s *InfluxDBStorage) QueryMetrics(ctx context.Context, containerID string, duration time.Duration) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -%s)
		|> filter(fn: (r) => r["_measurement"] == "container_metrics")
		|> filter(fn: (r) => r["container_id"] == "%s")
	`, s.bucket, duration.String(), containerID)
	result, err := s.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var metrics []map[string]interface{}
	for result.Next() {
		metrics = append(metrics, result.Record().Values())
	}
	return metrics, result.Err()
}
func (s *InfluxDBStorage) Close() {
	s.writeAPI.Flush()
	s.client.Close()
}