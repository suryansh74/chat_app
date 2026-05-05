package adapters

type EmailSenderPort interface {
	SendEmail(to, subject, body string) error
}
