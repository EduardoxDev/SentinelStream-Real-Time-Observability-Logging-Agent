package usecases
import (
	"context"
	"fmt"
	"observability-system/internal/domain/entities"
	"observability-system/internal/domain/ports"
)
type CheckAlertsUseCase struct {
	alertRepo     ports.AlertRepository
	notifier      ports.Notifier
	cpuThreshold  float64
	memThreshold  float64
}
func NewCheckAlertsUseCase(alertRepo ports.AlertRepository, notifier ports.Notifier, cpuThreshold, memThreshold float64) *CheckAlertsUseCase {
	return &CheckAlertsUseCase{
		alertRepo:    alertRepo,
		notifier:     notifier,
		cpuThreshold: cpuThreshold,
		memThreshold: memThreshold,
	}
}
func (uc *CheckAlertsUseCase) Execute(ctx context.Context, metrics *entities.ContainerMetrics) error {
	violations := metrics.ExceedsThreshold(uc.cpuThreshold, uc.memThreshold)
	for _, violation := range violations {
		alertType := entities.AlertType(violation)
		inCooldown, err := uc.alertRepo.IsInCooldown(ctx, metrics.ContainerID, alertType)
		if err != nil {
			return fmt.Errorf("failed to check cooldown: %w", err)
		}
		if inCooldown {
			continue
		}
		var value, threshold float64
		if alertType == entities.AlertTypeCPU {
			value = metrics.CPUPercent
			threshold = uc.cpuThreshold
		} else {
			value = metrics.MemoryPercent
			threshold = uc.memThreshold
		}
		alert := entities.NewAlert(metrics.ContainerID, metrics.ContainerName, alertType, value, threshold)
		alert.Message = fmt.Sprintf("%s usage (%.2f%%) exceeded threshold (%.2f%%)", alertType, value, threshold)
		if err := uc.alertRepo.Save(ctx, alert); err != nil {
			return fmt.Errorf("failed to save alert: %w", err)
		}
		if err := uc.notifier.Notify(ctx, alert); err != nil {
			return fmt.Errorf("failed to send notification: %w", err)
		}
	}
	return nil
}