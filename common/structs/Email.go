package structs

import (
	"amper/api/email/imap"
	"amper/api/email/message/mail"
	"amper/api/email/message/textproto"
	"amper/common/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type EmailStatusMetadata struct {
	Email     *string `json:"email"`
	Mailboxes *map[string]map[string]interface{}
}

func (esm *EmailStatusMetadata) Parse(data *[]byte) error {
	errUM := json.Unmarshal(*data, esm)
	if errUM != nil {
		util.Loggify(errUM)
		return fmt.Errorf("not able to parse the email status metadata")
	}
	return nil
}

func (esm *EmailStatusMetadata) Json() (*string, error) {
	if esm != nil {
		marshaled, errM := json.Marshal(*esm)
		if errM != nil {
			util.Loggify(errM)
			return nil, fmt.Errorf("email status metadata is invalid and can not be converted to json")
		}
		return util.PointerString(string(marshaled)), nil
	}
	return nil, fmt.Errorf("the email status metadata is nil and caan't be converted to json")
}

type Attachment struct {
	Name    *string
	Body    []byte
	Headers *map[string][]*textproto.HeaderField
}

type AttachmentMetadata struct {
	Name    *string
	Headers *map[string][]*textproto.HeaderField
}

type Email struct {
	ID           *string                              `json:"id"`
	Email        *string                              `json:"email"`
	SeqNum       uint32                               `json:"seqNum"`
	Flags        []imap.Flag                          `json:"flags"`
	Envelope     *imap.Envelope                       `json:"envelope"`
	InternalDate time.Time                            `json:"internalDate"`
	RFC822Size   int64                                `json:"RFC822Size"`
	UID          uint32                               `json:"uId"`
	Body         *string                              `json:"body"`
	BodyHTML     *string                              `json:"bodyHtml"`
	Headers      *map[string][]*textproto.HeaderField `json:"headers"`
	Attachments  *[]AttachmentMetadata                `json:"attachments"`
}

func (e *Email) Process(userId *int64, BodySection map[*imap.FetchItemBodySection][]byte) (attachments []Attachment, err error) {
	attachments = make([]Attachment, 0)
	emailAttachments := make([]AttachmentMetadata, 0)
	e.Attachments = &emailAttachments
	for _, bodySection := range BodySection {
		reader := bytes.NewReader(bodySection)
		mr, errR := mail.CreateReader(reader)
		if errR == nil {
			defer mr.Close()
			// Read each mail's part
			for {
				p, errNP := mr.NextPart()
				if errNP == io.EOF {
					break
				} else if errNP != nil {
					util.Loggify(errNP)
					err = fmt.Errorf("partial error, receivied a wrong email body, skipping")
					break
				}

				switch h := p.Header.(type) {
				case *mail.InlineHeader:
					e.Headers = p.Header.GetHeaders()
					contentType := p.Header.Get("Content-Type")
					b, errBR := ioutil.ReadAll(p.Body)
					if errBR == nil {
						bs := util.PointerString(string(b))
						if strings.Contains(contentType, "text/html") {
							e.BodyHTML = bs
						} else {
							e.Body = bs
						}
					} else {
						err = fmt.Errorf("partial error, not able to read message body, skipping")
					}
				case *mail.AttachmentHeader:
					fileName, _ := h.Filename()
					attachmentBytes, errAR := ioutil.ReadAll(p.Body)
					if errAR == nil {
						attachments = append(attachments, Attachment{
							Name:    util.PointerString(fileName),
							Body:    attachmentBytes,
							Headers: p.Header.GetHeaders(),
						})
						emailAttachments = append(emailAttachments, AttachmentMetadata{
							Name:    util.PointerString(fileName),
							Headers: p.Header.GetHeaders(),
						})
					} else {
						err = fmt.Errorf("partial error, not able to read message aattachment, skipping")
					}
				}
			}
		} else {
			util.Loggify(errR)
			err = fmt.Errorf("not able to process the received email message due to reading error")
		}
	}
	return attachments, err
}

func (e *Email) Json() (*string, error) {
	if e != nil {
		marshaled, errM := json.Marshal(*e)
		if errM != nil {
			util.Loggify(errM)
			return nil, fmt.Errorf("email is invalid and can not be converted to json")
		}
		return util.PointerString(string(marshaled)), nil
	}
	return nil, fmt.Errorf("the email is nil and caan't be converted to json")
}

type EmailsResult struct {
	Result
	Data  *[]Email `json:"data"`
	Count uint32   `json:"count"`
}

type Mailbox struct {
	Label       string `json:"label"`
	NumMessages int    `json:"numMessages"`
	SyncNumber  int    `json:"syncNumber"`
	Count       int    `json:"count"`
	All         bool   `json:"all"`
}

type ConfigureEmailResult struct {
	Result
	Data *[]Mailbox `json:"data"`
}

type SaveDraftResult struct {
	Result
	Id *string `json:"id"`
}

type PagePointer struct {
	Pages map[int]string `json:"pages"`
}
