package email

import (
	"amper/common/structs"
	"amper/common/util"
	"amper/properties/application"
	"amper/templates"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"text/template"
)

// Request is used to hold the email send properties
type Message struct {
	TemplateName *string
	Notification *structs.Notification
	From         *string
	To           *[]string
	CC           *[]string
	BCC          *[]string
	Subject      *string
	Body         *string
	Attachments  map[string][]byte
}

// Send email with the specified request properties
func (r *Message) Send() (bool, error) {
	config, errAP := application.Get()
	if errAP != nil {
		log.Printf("unable to get application properties: %v", errAP)
		return false, fmt.Errorf("unable to get application properties: %v", errAP)
	}
	mailInterface := config.GetString("mail.interface", "gmail")
	r.AddAttachment("notification/amper_logo.png", "001")
	r.AddAttachment("notification/facebook_icon.png", "002")
	r.AddAttachment("notification/twitter_icon.png", "003")
	switch mailInterface {
	case "smtp":
	default:
		r.From = util.PointerString(config.GetString("gmail.from", "admin"))
		return r.SendGmail()
	}
	return false, nil
}

func (r *Message) AddAttachment(fileName string, cid string) error {
	fs := templates.GetFS()
	data, err := fs.ReadFile(fileName)
	if err != nil {
		log.Printf("not able to add a attachment to the email with error %s", err.Error())
	} else {
		if r.Attachments == nil {
			r.Attachments = make(map[string][]byte)
		}
		r.Attachments[cid] = data
	}
	return err
}

// Parse parse the email template and apply the values
func (r *Message) Parse() error {
	fs := templates.GetFS()
	tE, errE := template.ParseFS(fs, "notification/email.tpl")
	if errE != nil {
		return errE
	}
	tC, errC := template.ParseFS(fs, fmt.Sprintf("notification/user/%s.tpl", *r.TemplateName))
	if errC != nil {
		return errC
	}
	tS, errS := template.ParseFS(fs, fmt.Sprintf("notification/user/%sSubject.tpl", *r.TemplateName))
	if errS != nil {
		return errS
	}
	bufC := new(bytes.Buffer)
	if err := tC.Execute(bufC, r.Notification); err != nil {
		return err
	}
	r.Notification.EmailContent = util.PointerString(bufC.String())
	bufS := new(bytes.Buffer)
	if err := tS.Execute(bufS, r.Notification); err != nil {
		return err
	}
	r.Notification.Headline = util.PointerString(bufS.String())
	r.Subject = r.Notification.Headline
	bufE := new(bytes.Buffer)
	if err := tE.Execute(bufE, r.Notification); err != nil {
		return err
	}
	r.Body = util.PointerString(bufE.String())
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := m.Attachments != nil && len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", *m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(*m.To, ",")))
	if m.CC != nil && len(*m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(*m.CC, ",")))
	}

	if m.BCC != nil && len(*m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(*m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}
	buf.WriteString("Content-Type:text/html; charset=\"UTF-8\";\n\n\r\n")
	buf.WriteString(*m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))
			buf.WriteString(fmt.Sprintf("Content-ID: %s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
