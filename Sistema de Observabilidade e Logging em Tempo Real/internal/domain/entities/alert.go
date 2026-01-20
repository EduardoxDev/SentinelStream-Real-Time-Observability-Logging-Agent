package entities
import "time"
type AlertType string
const (
	AlertTypeCPU    AlertType = "CPU"
	AlertTypeMemory AlertType = "MEMORY"
)
type Alert struct {
	ID            string
	ContainerID   string
	ContainerName string
	Type          AlertType
	Value         float64
	Threshold     float64
	Timestamp     time.Time
	Message       string
}
func NewAlert(containerID, containerName string, alertType AlertType, value, threshold float64) *Alert {
	return &Alert{
		ContainerID:   containerID,
		ContainerName: containerName,
		Type:          alertType,
		Value:         value,
		Threshold:     threshold,
		Timestamp:     time.Now(),
	}
}