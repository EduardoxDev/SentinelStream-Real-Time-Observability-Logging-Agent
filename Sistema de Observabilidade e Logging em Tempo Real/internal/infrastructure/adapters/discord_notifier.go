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
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}
type discordMessage struct {
	Content string         `json:"content"`
	Embeds  []discordEmbed `json:"embeds"`
}
type discordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []discordEmbedField `json:"fields"`
	Footer      discordEmbedFooter  `json:"footer"`
	Timestamp   string              `json:"timestamp"`
}
type discordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type discordEmbedFooter struct {
	Text string `json:"text"`
}
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
func (n *DiscordNotifier) Notify(ctx context.Context, alert *entities.Alert) error {
	color := 16776960
	if alert.Value > alert.Threshold*1.2 {
		color = 16711680
	}
	msg := discordMessage{
		Content: "🚨 **Container Alert**",
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("%s - %s Alert", alert.ContainerName, alert.Type),
				Description: alert.Message,
				Color:       color,
				Fields: []discordEmbedField{
					{
						Name:   "Container",
						Value:  alert.ContainerName,
						Inline: true,
					},
					{
						Name:   "Type",
						Value:  string(alert.Type),
						Inline: true,
					},
					{
						Name:   "Value",
						Value:  fmt.Sprintf("%.2f%%", alert.Value),
						Inline: true,
					},
					{
						Name:   "Threshold",
						Value:  fmt.Sprintf("%.2f%%", alert.Threshold),
						Inline: true,
					},
					{
						Name:   "Container ID",
						Value:  alert.ContainerID[:12],
						Inline: false,
					},
				},
				Footer: discordEmbedFooter{
					Text: "Observability System",
				},
				Timestamp: alert.Timestamp.Format(time.RFC3339),
			},
		},
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord notification: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord returned status %d", resp.StatusCode)
	}
	return nil
}