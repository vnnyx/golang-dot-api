package util

import (
	"bytes"

	"github.com/vnnyx/golang-dot-api/infrastructure"
	"gopkg.in/gomail.v2"
)

func SendEmailTo(email string, buff *bytes.Buffer) error {
	configuration := infrastructure.NewConfig(".env")
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply-dot@email.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html", buff.String())
	d := gomail.NewDialer(configuration.MailHost, configuration.MailPort, configuration.MailUsername, configuration.MailPassword)
	err := d.DialAndSend(m)
	return err
}
