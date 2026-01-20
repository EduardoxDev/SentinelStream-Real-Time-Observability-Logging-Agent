package entities
import "time"
type ContainerMetrics struct {
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
func (m *ContainerMetrics) IsHealthy(cpuThreshold, memoryThreshold float64) bool {
	return m.CPUPercent <= cpuThreshold && m.MemoryPercent <= memoryThreshold
}
func (m *ContainerMetrics) ExceedsThreshold(cpuThreshold, memoryThreshold float64) []string {
	var violations []string
	if m.CPUPercent > cpuThreshold {
		violations = append(violations, "CPU")
	}
	if m.MemoryPercent > memoryThreshold {
		violations = append(violations, "MEMORY")
	}
	return violations
}