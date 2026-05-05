package email

import (
	"fmt"
	"net/smtp"

	"github.com/suryansh74/chat_app/pkg/logger"
)

type Sender struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

func NewSender(host, port, username, password string) *Sender {
	return &Sender{
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
	}
}

func (s *Sender) SendEmail(to, subject, body string) error {
	from := s.username
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	logger.Log.Info("EmailSender: preparing to send", "to", to, "subject", subject, "smtp_addr", addr)

	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	logger.Log.Info("EmailSender: attempting to send email", "to", to)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		logger.Log.Error("EmailSender: failed to send email", "error", err.Error(), "to", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Log.Info("EmailSender: email sent successfully", "to", to)
	return nil
}
