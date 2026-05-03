package email

import (
	"fmt"
	"net/smtp"
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

	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
