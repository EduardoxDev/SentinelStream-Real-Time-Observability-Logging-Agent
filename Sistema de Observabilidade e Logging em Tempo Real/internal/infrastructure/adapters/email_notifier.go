package adapters
import (
	"context"
	"fmt"
	"net/smtp"
	"time"
	"observability-system/internal/domain/entities"
)
type EmailNotifier struct {
	smtpHost string
	smtpPort string
	from     string
	password string
	to       []string
}
func NewEmailNotifier(smtpHost, smtpPort, from, password string, to []string) *EmailNotifier {
	return &EmailNotifier{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
		password: password,
		to:       to,
	}
}
func (n *EmailNotifier) Notify(ctx context.Context, alert *entities.Alert) error {
	subject := fmt.Sprintf("🚨 Alert: %s - %s", alert.ContainerName, alert.Type)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .alert { background: #fee; border-left: 4px solid #f44; padding: 20px; }
        .info { color: #666; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="alert">
        <h2>🚨 Container Alert</h2>
        <p><strong>Container:</strong> %s</p>
        <p><strong>Type:</strong> %s</p>
        <p><strong>Message:</strong> %s</p>
        <p><strong>Value:</strong> %.2f%% (Threshold: %.2f%%)</p>
        <div class="info">
            <p>Time: %s</p>
            <p>Container ID: %s</p>
        </div>
    </div>
</body>
</html>
	`, alert.ContainerName, alert.Type, alert.Message, alert.Value, alert.Threshold,
		alert.Timestamp.Format(time.RFC3339), alert.ContainerID)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", n.from, n.to[0], subject, body))
	auth := smtp.PlainAuth("", n.from, n.password, n.smtpHost)
	addr := fmt.Sprintf("%s:%s", n.smtpHost, n.smtpPort)
	return smtp.SendMail(addr, auth, n.from, n.to, msg)
}