package alerts
import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
	"observability-system/internal/collector"
)
type AlertConfig struct {
	CPUThreshold    float64
	MemoryThreshold float64
	CooldownPeriod  time.Duration
}
type Alerter struct {
	config      AlertConfig
	redisClient *redis.Client
	notifiers   []Notifier
}
type Notifier interface {
	Send(alert Alert) error
}
type Alert struct {
	ContainerID   string
	ContainerName string
	Type          string
	Value         float64
	Threshold     float64
	Timestamp     time.Time
	Message       string
}
func NewAlerter(redisClient *redis.Client, config AlertConfig) *Alerter {
	return &Alerter{
		config:      config,
		redisClient: redisClient,
		notifiers:   []Notifier{&ConsoleNotifier{}},
	}
}
func (a *Alerter) AddNotifier(notifier Notifier) {
	a.notifiers = append(a.notifiers, notifier)
}
func (a *Alerter) CheckMetrics(ctx context.Context, stats *collector.ContainerStats) error {
	if stats.CPUPercent > a.config.CPUThreshold {
		if err := a.triggerAlert(ctx, Alert{
			ContainerID:   stats.ContainerID,
			ContainerName: stats.ContainerName,
			Type:          "CPU",
			Value:         stats.CPUPercent,
			Threshold:     a.config.CPUThreshold,
			Timestamp:     stats.Timestamp,
			Message:       fmt.Sprintf("CPU usage (%.2f%%) exceeded threshold (%.2f%%)", stats.CPUPercent, a.config.CPUThreshold),
		}); err != nil {
			return err
		}
	}
	if stats.MemoryPercent > a.config.MemoryThreshold {
		if err := a.triggerAlert(ctx, Alert{
			ContainerID:   stats.ContainerID,
			ContainerName: stats.ContainerName,
			Type:          "MEMORY",
			Value:         stats.MemoryPercent,
			Threshold:     a.config.MemoryThreshold,
			Timestamp:     stats.Timestamp,
			Message:       fmt.Sprintf("Memory usage (%.2f%%) exceeded threshold (%.2f%%)", stats.MemoryPercent, a.config.MemoryThreshold),
		}); err != nil {
			return err
		}
	}
	return nil
}
func (a *Alerter) triggerAlert(ctx context.Context, alert Alert) error {
	key := fmt.Sprintf("alert:%s:%s", alert.ContainerID, alert.Type)
	exists, err := a.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	for _, notifier := range a.notifiers {
		if err := notifier.Send(alert); err != nil {
			log.Printf("Failed to send alert: %v", err)
		}
	}
	return a.redisClient.Set(ctx, key, "1", a.config.CooldownPeriod).Err()
}
type ConsoleNotifier struct{}
func (n *ConsoleNotifier) Send(alert Alert) error {
	log.Printf("🚨 ALERT [%s] %s - %s: %s",
		alert.Timestamp.Format(time.RFC3339),
		alert.ContainerName,
		alert.Type,
		alert.Message,
	)
	return nil
}