package adapters

import (
	"fmt"
	"net/smtp"

	"github.com/suryansh74/chat_app/pkg/logger"
)

type GmailAdapter struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

func NewGmailAdapter(host, port, username, password string) *GmailAdapter {
	return &GmailAdapter{
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
	}
}

func (a *GmailAdapter) SendEmail(to, subject, body string) error {
	from := a.username
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", a.smtpHost, a.smtpPort)

	logger.Log.Info("GmailAdapter: sending email", "to", to, "subject", subject)

	auth := smtp.PlainAuth("", a.username, a.password, a.smtpHost)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		logger.Log.Error("GmailAdapter: failed to send", "error", err.Error(), "to", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Log.Info("GmailAdapter: email sent", "to", to)
	return nil
}
