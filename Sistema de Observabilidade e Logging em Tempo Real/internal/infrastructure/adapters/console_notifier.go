package adapters
import (
	"context"
	"log"
	"time"
	"observability-system/internal/domain/entities"
)
type ConsoleNotifier struct{}
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}
func (n *ConsoleNotifier) Notify(ctx context.Context, alert *entities.Alert) error {
	log.Printf("🚨 ALERT [%s] %s - %s: %s",
		alert.Timestamp.Format(time.RFC3339),
		alert.ContainerName,
		alert.Type,
		alert.Message,
	)
	return nil
}