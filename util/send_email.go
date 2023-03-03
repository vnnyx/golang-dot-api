package util

import (
	"bytes"

	"gopkg.in/gomail.v2"
)

func SendEmailTo(email string, buff *bytes.Buffer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply-sixmath@email.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html", buff.String())
	// port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "d026e1e22a90dd", "a5adcbfe188b75")
	err := d.DialAndSend(m)
	return err
}
