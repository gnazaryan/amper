package email

import (
	"amper/properties/application"
	"fmt"
	"log"
	"net/smtp"
)

// Send SMTP email with the specified request properties
func (r *Message) SendSmtp() (bool, error) {
	config, errAP := application.Get()
	if errAP != nil {
		log.Printf("unable to get application properties: %v", errAP)
		return false, fmt.Errorf("unable to get application properties: %v", errAP)
	}
	smtpHost := config.GetString("mail.smtp.host", "192.168.2.11")
	smtpPort := config.GetString("mail.smtp.port", "25")
	smtpFrom := config.GetString("mail.smtp.from", "admin@mail.amp-er.com")
	smtpUsername := config.GetString("mail.smtp.username", "admin@mail.amp-er.com")
	smtpPassword := config.GetString("mail.smtp.password", "amper123")
	if r.From == nil {
		r.From = &smtpFrom
	}
	var auth smtp.Auth = smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + *r.Subject + "!\n"
	msg := []byte(subject + mime + "\n" + *r.Body)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	if err := smtp.SendMail(addr, auth, *r.From, *r.To, msg); err != nil {
		return false, err
	}
	return true, nil
}
