package adapters

import (
	"fmt"
	"net/smtp"

	"github.com/suryansh74/chat_app/pkg/logger"
)

type MailpitAdapter struct {
	smtpHost string
	smtpPort string
}

func NewMailpitAdapter(host, port string) *MailpitAdapter {
	return &MailpitAdapter{
		smtpHost: host,
		smtpPort: port,
	}
}

func (a *MailpitAdapter) SendEmail(to, subject, body string) error {
	from := "test@example.com"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", a.smtpHost, a.smtpPort)

	logger.Log.Info("MailpitAdapter: sending email", "to", to, "subject", subject, "addr", addr)

	// Mailpit doesn't require auth - send directly
	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
	if err != nil {
		logger.Log.Error("MailpitAdapter: failed to send", "error", err.Error(), "to", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Log.Info("MailpitAdapter: email sent", "to", to)
	return nil
}
