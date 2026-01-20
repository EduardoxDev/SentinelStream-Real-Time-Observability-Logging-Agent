package adapters
import (
	"context"
	"log"
	"sync"
	"observability-system/internal/domain/entities"
	"observability-system/internal/domain/ports"
)
type MultiNotifier struct {
	notifiers []ports.Notifier
}
func NewMultiNotifier(notifiers ...ports.Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}
func (m *MultiNotifier) AddNotifier(notifier ports.Notifier) {
	m.notifiers = append(m.notifiers, notifier)
}
func (m *MultiNotifier) Notify(ctx context.Context, alert *entities.Alert) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(m.notifiers))
	for _, notifier := range m.notifiers {
		wg.Add(1)
		go func(n ports.Notifier) {
			defer wg.Done()
			if err := n.Notify(ctx, alert); err != nil {
				log.Printf("Notifier failed: %v", err)
				errChan <- err
			}
		}(notifier)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}