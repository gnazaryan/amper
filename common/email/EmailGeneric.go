package email

import (
	"amper/common/util"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"
)

// Send Gmail email with the specified request properties
func (m *Message) SendEmail(username string, password string, serverName string, port int, auth string) (bool, error) {

	var authentication smtp.Auth

	if auth == "Plain auth" {
		authentication = smtp.PlainAuth("", username, password, serverName)
	} else {
		authentication = LoginAuth(username, password)
	}
	err := smtp.SendMail(serverName+":"+strconv.Itoa(port),
		authentication,
		*m.From, *m.To, m.ToBytes())

	if err != nil {
		util.Loggify(err)
		return false, fmt.Errorf("not able to successfully send email for host %s", serverName)
	}
	return true, nil
}

//https://stackoverflow.com/questions/58804817/setting-up-standard-go-net-smtp-with-office-365-fails-with-error-tls-first-rec

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
