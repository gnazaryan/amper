package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

// CheckEmail is running a request to users email server to retrieve the latest emails
func CheckEmails(userID *int64, Email *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "email": Email})
	if err != nil {
		return false, err
	}
	//TODO perform authorization for get user action with userId
	success, error := business.CheckEmails(userID, Email)
	return success, error
}

// FetchEmail is retrieving the users latest emails
func FetchEmails(userID *int64, Email *string, Box *string, Search *string, Start *int, Limit *int, Pointer *structs.PagePointer) ([]structs.Email, uint32, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "email": Email, "start": Start, "limit": Limit})
	if err != nil {
		return nil, 0, err
	}
	//TODO perform authorization for get user action with userId
	data, count, error := business.FetchEmails(userID, Email, Box, Search, Start, Limit, Pointer)
	return data, count, error
}

func MoveEmails(userID *int64, Emails *[]map[string]interface{}, From *string, To *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "emails": Emails, "From": From, "to": To})
	if err != nil {
		return false, err
	}
	return business.MoveEmails(userID, Emails, From, To)
}

func FlagEmails(userID *int64, Emails *[]map[string]interface{}, Box *string, Flags *[]string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "emails": Emails, "box": Box, "flags": Flags})
	if err != nil {
		return false, err
	}
	return business.FlagEmails(userID, Emails, Box, Flags)
}

func ConfigureEmail(userID *int64, Email *string, Password *string) (*[]structs.Mailbox, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "email": Email, "password": Password})
	if err != nil {
		return nil, err
	}
	return business.ConfigureEmail(userID, Email, Password)
}

func Mailboxes(userID *int64, Email *string) (*[]structs.Mailbox, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "email": Email})
	if err != nil {
		return nil, err
	}
	return business.Mailboxes(userID, Email)
}

func SaveEmailDraft(userID *int64, Id *string, From *string, TO *string, CC *string, BCC *string, Subject *string, Content *string) (*string, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "from": From})
	if err != nil {
		return nil, err
	}
	return business.SaveEmailDraft(userID, Id, From, TO, CC, BCC, Subject, Content)
}

func Send(userID *int64, Id *string, From *string, TO *string, CC *string, BCC *string, Subject *string, Content *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "from": From, "to": TO})
	if err != nil {
		return false, err
	}
	return business.Send(userID, Id, From, TO, CC, BCC, Subject, Content)
}
