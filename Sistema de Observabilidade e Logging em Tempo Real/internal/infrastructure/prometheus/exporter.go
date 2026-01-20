package prometheus
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)
type MetricsExporter struct {
	cpuUsage    *prometheus.GaugeVec
	memoryUsage *prometheus.GaugeVec
	networkRx   *prometheus.CounterVec
	networkTx   *prometheus.CounterVec
}
func NewMetricsExporter() *MetricsExporter {
	return &MetricsExporter{
		cpuUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "container_cpu_usage_percent",
				Help: "Container CPU usage percentage",
			},
			[]string{"container_id", "container_name"},
		),
		memoryUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "container_memory_usage_percent",
				Help: "Container memory usage percentage",
			},
			[]string{"container_id", "container_name"},
		),
		networkRx: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "container_network_rx_bytes_total",
				Help: "Container network received bytes",
			},
			[]string{"container_id", "container_name"},
		),
		networkTx: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "container_network_tx_bytes_total",
				Help: "Container network transmitted bytes",
			},
			[]string{"container_id", "container_name"},
		),
	}
}
func (e *MetricsExporter) RecordMetrics(containerID, containerName string, cpuPercent, memoryPercent float64, networkRx, networkTx uint64) {
	labels := prometheus.Labels{
		"container_id":   containerID,
		"container_name": containerName,
	}
	e.cpuUsage.With(labels).Set(cpuPercent)
	e.memoryUsage.With(labels).Set(memoryPercent)
	e.networkRx.With(labels).Add(float64(networkRx))
	e.networkTx.With(labels).Add(float64(networkTx))
}