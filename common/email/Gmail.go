package email

import (
	"log"
	"net/smtp"
)

// Send Gmail email with the specified request properties
func (m *Message) SendGmail() (bool, error) {
	from := "smtp.amper.cloud"
	pass := "exaxlndmfbnxgiky" //cloud.amper.smtp

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		*m.From, *m.To, m.ToBytes())

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false, err
	}
	return false, nil
}
