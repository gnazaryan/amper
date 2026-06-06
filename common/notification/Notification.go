package notification

import (
	"amper/common/email"
	"amper/common/structs"
	"amper/common/util"
	"amper/properties/application"
	"fmt"
	"log"
	"strings"
)

// UserRegistration is a constant representing the remplate name for registering a new user
var UserRegistration string = "userRegistration"

// Send is used to send a notification to Amper users
func Send(user *structs.User, templateName *string, data interface{}) (result bool, err error) {
	config, errAP := application.Get()
	if errAP != nil {
		log.Printf("unable to get application properties: %v", errAP)
	}
	uinodes := config.GetString("uinodes", "http://localhost:3000/")
	uinode := strings.Split(uinodes, ",")[0]
	var emailRequest email.Message = email.Message{
		TemplateName: templateName,
		Notification: &structs.Notification{
			Data:           data,
			AmperResources: util.PointerString(config.GetString("resources", "192.168.2.11:8080")),
			AmperURL:       &uinode,
		},
		To: util.PSA([]string{*user.Email}),
	}
	errP := emailRequest.Parse()
	if errP != nil {
		return false, fmt.Errorf("unable to parse notification for user: %s", *user.Username)
	}
	return emailRequest.Send()
	/*path, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	if err != nil {
		log.Printf("unable to get root directory: %v", err)
		return fmt.Errorf("unable to get the root directory to send notification for user: %s", *user.Username)
	}
	t, errT := template.ParseFiles(path)
	if errT != nil {
		log.Printf("unable to parse template file: %v", errT)
		return fmt.Errorf("unable to parse template to send notification for user: %s", *user.Username)
	}
	email, errE := business.Read(util.PointerString(path + "/templates/notification/email.tpl"))
	emailContent, errC := business.Read(util.PointerString(path + fmt.Sprintf("/templates/notification/user/%s.tpl", *template)))
	subject, errS := business.Read(util.PointerString(path + fmt.Sprintf("/templates/notification/user/%sSubject.tpl", *template)))
	if errE != nil || errC != nil || errS != nil {
		log.Printf("Unable to read email template: %v, %v, %v", errE, errC, errS)
	}
	replacer := strings.NewReplacer(util.FlatArray(values)...)
	emailContent = replacer.Replace(emailContent)
	subject = replacer.Replace(subject)

	config, errAP := application.Get()
	if errAP != nil {
		log.Printf("unable to get application properties: %v", errAP)
		return fmt.Errorf("unable to get application properties to send notification for user: %s", *user.Username)
	}
	uinodes := config.GetString("uinodes", "http://localhost:3000/")
	uinode := strings.Split(uinodes, ",")[0]
	(*values)["{amperResources}"] = config.GetString("resources", "192.168.2.11:8080")
	(*values)["{amperUrl}"] = uinode
	(*values)["{headline}"] = replacer.Replace(email)
	(*values)["{emailContent}"] = replacer.Replace(email)

	replacer = strings.NewReplacer(util.FlatArray(values)...)
	email = replacer.Replace(email)*/
}
