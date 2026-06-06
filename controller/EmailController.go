package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// WidgetController is responsible for dispatching requests related to
// dashboard and underlying widgets managment functionalities
func EmailController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "check":
			resultStruct = checkEmail(userID, w, r)
		case "fetch":
			resultStruct = fetchEmails(userID, w, r)
		case "move":
			resultStruct = moveEmails(userID, w, r)
		case "flag":
			resultStruct = flagEmails(userID, w, r)
		case "configure":
			resultStruct = configureEmail(userID, w, r)
		case "mailboxes":
			resultStruct = mailboxes(userID, w, r)
		case "saveDraft":
			resultStruct = saveDraft(userID, w, r)
		case "send":
			resultStruct = send(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func send(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Id      *string
		From    *string
		TO      *string
		CC      *string
		BCC     *string
		Subject *string
		Content *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.Send(userID, parameters.Id, parameters.From, parameters.TO, parameters.CC, parameters.BCC, parameters.Subject, parameters.Content)
	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func saveDraft(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.SaveDraftResult) {
	var parameters struct {
		Id      *string
		From    *string
		TO      *string
		CC      *string
		BCC     *string
		Subject *string
		Content *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	id, err := authorization.SaveEmailDraft(userID, parameters.Id, parameters.From, parameters.TO, parameters.CC, parameters.BCC, parameters.Subject, parameters.Content)
	if err == nil {
		result.Success = true
		result.Id = id
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func mailboxes(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ConfigureEmailResult) {
	var parameters struct {
		Email *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	Mailboxes, err := authorization.Mailboxes(userID, parameters.Email)
	if err == nil {
		result.Success = true
		result.Data = Mailboxes
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func configureEmail(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ConfigureEmailResult) {
	var parameters struct {
		Email    *string
		Password *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	Mailboxes, err := authorization.ConfigureEmail(userID, parameters.Email, parameters.Password)
	if err == nil {
		result.Success = true
		result.Data = Mailboxes
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func flagEmails(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Emails *[]map[string]interface{}
		Box    *string
		Flags  *[]string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	var boxUnescaped string
	if parameters.Box != nil {
		boxUnescaped, _ = url.QueryUnescape(*parameters.Box)
	} else {
		boxUnescaped = ""
	}

	parameters.Box = &boxUnescaped
	success, err := authorization.FlagEmails(userID, parameters.Emails, parameters.Box, parameters.Flags)
	if err == nil && success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func moveEmails(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Emails *[]map[string]interface{}
		From   *string
		To     *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.MoveEmails(userID, parameters.Emails, parameters.From, parameters.To)
	if err == nil && success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func fetchEmails(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.EmailsResult) {
	var parameters struct {
		Email   *string
		Box     *string
		Search  *string
		Start   *int
		Limit   *int
		Pointer *structs.PagePointer
	}
	json.NewDecoder(r.Body).Decode(&parameters)

	var boxUnescaped string
	if parameters.Box != nil {
		boxUnescaped, _ = url.QueryUnescape(*parameters.Box)
	} else {
		boxUnescaped = ""
	}

	parameters.Box = &boxUnescaped
	data, count, err := authorization.FetchEmails(userID, parameters.Email, parameters.Box, parameters.Search, parameters.Start, parameters.Limit, parameters.Pointer)
	if err == nil {
		result.Success = true
		result.Data = &data
		result.Count = count
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func checkEmail(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Email *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.CheckEmails(userID, parameters.Email)
	if err == nil || success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}
