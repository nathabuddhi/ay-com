package email

import (
	"os"

	"gopkg.in/gomail.v2"
)

var (
	smtpHost     string
	smtpPort     int
	smtpEmail    string
	smtpPassword string
)

func InitEmailSender() {
	smtpHost = "smtp.gmail.com"
	smtpPort = 587
	smtpEmail = os.Getenv("SMTP_EMAIL")
	smtpPassword = os.Getenv("SMTP_KEY")
}

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	return d.DialAndSend(m)
}
