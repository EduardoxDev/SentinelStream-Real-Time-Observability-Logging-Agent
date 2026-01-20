package adapters
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"observability-system/internal/domain/entities"
)
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}
type slackMessage struct {
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments"`
}
type slackAttachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
	Ts     int64  `json:"ts"`
}
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
func (n *SlackNotifier) Notify(ctx context.Context, alert *entities.Alert) error {
	color := "warning"
	if alert.Value > alert.Threshold*1.2 {
		color = "danger"
	}
	msg := slackMessage{
		Text: "🚨 *Container Alert*",
		Attachments: []slackAttachment{
			{
				Color:  color,
				Title:  fmt.Sprintf("%s - %s Alert", alert.ContainerName, alert.Type),
				Text:   alert.Message,
				Footer: "Observability System",
				Ts:     alert.Timestamp.Unix(),
			},
		},
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status %d", resp.StatusCode)
	}
	return nil
}